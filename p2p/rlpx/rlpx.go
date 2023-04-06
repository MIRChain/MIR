// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package rlpx implements the RLPx transport protocol.
package rlpx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"io"
	mrand "math/rand"
	"net"
	"reflect"
	"time"

	"github.com/golang/snappy"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/ecies"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"golang.org/x/crypto/sha3"
)

// Conn[T,P] is an RLPx network Connection. It wraps a low-level network Connection. The
// underlying Connection should not be used for other activity when it is wrapped by Conn[T,P].
//
// Before sending messages, a handshake must be performed by calling the Handshake method.
// This type is not generally safe for concurrent use, but reading and writing of messages
// may happen concurrently after the handshake.
type Conn[T crypto.PrivateKey, P crypto.PublicKey] struct {
	dialDest  P
	conn      net.Conn
	handshake *handshakeState
	snappy    bool
}

type handshakeState struct {
	enc cipher.Stream
	dec cipher.Stream

	macCipher  cipher.Block
	egressMAC  hash.Hash
	ingressMAC hash.Hash
}

// NewConn[T,P] wraps the given network Connection. If dialDest is non-nil, the Connection
// behaves as the initiator during the handshake.
func NewConn[T crypto.PrivateKey, P crypto.PublicKey](conn net.Conn, dialDest P) *Conn[T,P] {
	return &Conn[T,P]{
		dialDest: dialDest,
		conn:     conn,
	}
}

// SetSnappy enables or disables snappy compression of messages. This is usually called
// after the devp2p Hello message exchange when the negotiated version indicates that
// compression is available on both ends of the Connection.
func (c *Conn[T,P]) SetSnappy(snappy bool) {
	c.snappy = snappy
}

// SetReadDeadline sets the deadline for all future read operations.
func (c *Conn[T,P]) SetReadDeadline(time time.Time) error {
	return c.conn.SetReadDeadline(time)
}

// SetWriteDeadline sets the deadline for all future write operations.
func (c *Conn[T,P]) SetWriteDeadline(time time.Time) error {
	return c.conn.SetWriteDeadline(time)
}

// SetDeadline sets the deadline for all future read and write operations.
func (c *Conn[T,P]) SetDeadline(time time.Time) error {
	return c.conn.SetDeadline(time)
}

// Read reads a message from the Connection.
func (c *Conn[T,P]) Read() (code uint64, data []byte, wireSize int, err error) {
	if c.handshake == nil {
		panic("can't ReadMsg before handshake")
	}

	frame, err := c.handshake.readFrame(c.conn)
	if err != nil {
		return 0, nil, 0, err
	}
	code, data, err = rlp.SplitUint64(frame)
	if err != nil {
		return 0, nil, 0, fmt.Errorf("invalid message code: %v", err)
	}
	wireSize = len(data)

	// If snappy is enabled, verify and decompress message.
	if c.snappy {
		var actualSize int
		actualSize, err = snappy.DecodedLen(data)
		if err != nil {
			return code, nil, 0, err
		}
		if actualSize > maxUint24 {
			return code, nil, 0, errPlainMessageTooLarge
		}
		data, err = snappy.Decode(nil, data)
	}
	return code, data, wireSize, err
}

func (h *handshakeState) readFrame(conn io.Reader) ([]byte, error) {
	// read the header
	headbuf := make([]byte, 32)
	if _, err := io.ReadFull(conn, headbuf); err != nil {
		return nil, err
	}

	// verify header mac
	shouldMAC := updateMAC(h.ingressMAC, h.macCipher, headbuf[:16])
	if !hmac.Equal(shouldMAC, headbuf[16:]) {
		return nil, errors.New("bad header MAC")
	}
	h.dec.XORKeyStream(headbuf[:16], headbuf[:16]) // first half is now decrypted
	fsize := readInt24(headbuf)
	// ignore protocol type for now

	// read the frame content
	var rsize = fsize // frame size rounded up to 16 byte boundary
	if padding := fsize % 16; padding > 0 {
		rsize += 16 - padding
	}
	framebuf := make([]byte, rsize)
	if _, err := io.ReadFull(conn, framebuf); err != nil {
		return nil, err
	}

	// read and validate frame MAC. we can re-use headbuf for that.
	h.ingressMAC.Write(framebuf)
	fmacseed := h.ingressMAC.Sum(nil)
	if _, err := io.ReadFull(conn, headbuf[:16]); err != nil {
		return nil, err
	}
	shouldMAC = updateMAC(h.ingressMAC, h.macCipher, fmacseed)
	if !hmac.Equal(shouldMAC, headbuf[:16]) {
		return nil, errors.New("bad frame MAC")
	}

	// decrypt frame content
	h.dec.XORKeyStream(framebuf, framebuf)
	return framebuf[:fsize], nil
}

// Write writes a message to the connection.
//
// Write returns the written size of the message data. This may be less than or equal to
// len(data) depending on whether snappy compression is enabled.
func (c *Conn[T,P]) Write(code uint64, data []byte) (uint32, error) {
	if c.handshake == nil {
		panic("can't WriteMsg before handshake")
	}
	if len(data) > maxUint24 {
		return 0, errPlainMessageTooLarge
	}
	if c.snappy {
		data = snappy.Encode(nil, data)
	}

	wireSize := uint32(len(data))
	err := c.handshake.writeFrame(c.conn, code, data)
	return wireSize, err
}

func (h *handshakeState) writeFrame(conn io.Writer, code uint64, data []byte) error {
	ptype, _ := rlp.EncodeToBytes(code)

	// write header
	headbuf := make([]byte, 32)
	fsize := len(ptype) + len(data)
	if fsize > maxUint24 {
		return errPlainMessageTooLarge
	}
	putInt24(uint32(fsize), headbuf)
	copy(headbuf[3:], zeroHeader)
	h.enc.XORKeyStream(headbuf[:16], headbuf[:16]) // first half is now encrypted

	// write header MAC
	copy(headbuf[16:], updateMAC(h.egressMAC, h.macCipher, headbuf[:16]))
	if _, err := conn.Write(headbuf); err != nil {
		return err
	}

	// write encrypted frame, updating the egress MAC hash with
	// the data written to conn.
	tee := cipher.StreamWriter{S: h.enc, W: io.MultiWriter(conn, h.egressMAC)}
	if _, err := tee.Write(ptype); err != nil {
		return err
	}
	if _, err := tee.Write(data); err != nil {
		return err
	}
	if padding := fsize % 16; padding > 0 {
		if _, err := tee.Write(zero16[:16-padding]); err != nil {
			return err
		}
	}

	// write frame MAC. egress MAC hash is up to date because
	// frame content was written to it as well.
	fmacseed := h.egressMAC.Sum(nil)
	mac := updateMAC(h.egressMAC, h.macCipher, fmacseed)
	_, err := conn.Write(mac)
	return err
}

func readInt24(b []byte) uint32 {
	return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
}

func putInt24(v uint32, b []byte) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

// updateMAC reseeds the given hash with encrypted seed.
// it returns the first 16 bytes of the hash sum after seeding.
func updateMAC(mac hash.Hash, block cipher.Block, seed []byte) []byte {
	aesbuf := make([]byte, aes.BlockSize)
	block.Encrypt(aesbuf, mac.Sum(nil))
	for i := range aesbuf {
		aesbuf[i] ^= seed[i]
	}
	mac.Write(aesbuf)
	return mac.Sum(nil)[:16]
}

// Handshake performs the handshake. This must be called before any data is written
// or read from the Connection.
func (c *Conn[T,P]) Handshake(prv T) (P, error) {
	var (
		sec Secrets[T,P]
		err error
	)
	if !reflect.ValueOf(&c.dialDest).Elem().IsZero() {
		sec, err = initiatorEncHandshake(c.conn, prv, c.dialDest)
	} else {
		sec, err = receiverEncHandshake[T,P](c.conn, prv)
	}
	if err != nil {
		return crypto.ZeroPublicKey[P](), err
	}
	c.InitWithSecrets(sec)
	return sec.remote, err
}

// InitWithSecrets injects connection secrets as if a handshake had
// been performed. This cannot be called after the handshake.
func (c *Conn[T,P]) InitWithSecrets(sec Secrets[T,P]) {
	if c.handshake != nil {
		panic("can't handshake twice")
	}
	macc, err := aes.NewCipher(sec.MAC)
	if err != nil {
		panic("invalid MAC secret: " + err.Error())
	}
	encc, err := aes.NewCipher(sec.AES)
	if err != nil {
		panic("invalid AES secret: " + err.Error())
	}
	// we use an all-zeroes IV for AES because the key used
	// for encryption is ephemeral.
	iv := make([]byte, encc.BlockSize())
	c.handshake = &handshakeState{
		enc:        cipher.NewCTR(encc, iv),
		dec:        cipher.NewCTR(encc, iv),
		macCipher:  macc,
		egressMAC:  sec.EgressMAC,
		ingressMAC: sec.IngressMAC,
	}
}

// Close closes the underlying network connection.
func (c *Conn[T,P]) Close() error {
	return c.conn.Close()
}

// Constants for the handshake.
const (
	maxUint24 = int(^uint32(0) >> 8)

	sskLen = 16                     // ecies.MaxSharedKeyLength(pubKey) / 2
	sigLen = crypto.SignatureLength // elliptic S256
	pubLen = 64                     // 512 bit pubkey in uncompressed representation without format byte
	shaLen = 32                     // hash length (for nonce etc)

	authMsgLen  = sigLen + shaLen + pubLen + shaLen + 1
	authRespLen = pubLen + shaLen + 1

	eciesOverhead = 65 /* pubkey */ + 16 /* IV */ + 32 /* MAC */

	encAuthMsgLen  = authMsgLen + eciesOverhead  // size of encrypted pre-EIP-8 initiator handshake
	encAuthRespLen = authRespLen + eciesOverhead // size of encrypted pre-EIP-8 handshake reply
)

var (
	// this is used in place of actual frame header data.
	// TODO: replace this when Msg contains the protocol type code.
	zeroHeader = []byte{0xC2, 0x80, 0x80}
	// sixteen zero bytes
	zero16 = make([]byte, 16)

	// errPlainMessageTooLarge is returned if a decompressed message length exceeds
	// the allowed 24 bits (i.e. length >= 16MB).
	errPlainMessageTooLarge = errors.New("message length >= 16MB")
)

// Secrets represents the connection secrets which are negotiated during the handshake.
type Secrets [T crypto.PrivateKey, P crypto.PublicKey] struct {
	AES, MAC              []byte
	EgressMAC, IngressMAC hash.Hash
	remote                P
}

// encHandshake contains the state of the encryption handshake.
type encHandshake [T crypto.PrivateKey, P crypto.PublicKey] struct {
	initiator            bool
	remote               *ecies.PublicKey[P]  // remote-pubk
	initNonce, respNonce []byte            // nonce
	randomPrivKey        *ecies.PrivateKey[T,P] // ecdhe-random
	remoteRandomPub      *ecies.PublicKey[P]  // ecdhe-random-pubk
}

// RLPx v4 handshake auth (defined in EIP-8).
type authMsgV4 [T crypto.PrivateKey, P crypto.PublicKey] struct {
	gotPlain bool // whether read packet had plain format.

	Signature       [sigLen]byte
	InitiatorPubkey [pubLen]byte
	Nonce           [shaLen]byte
	Version         uint

	// Ignore additional fields (forward-compatibility)
	Rest []rlp.RawValue `rlp:"tail"`
}

// RLPx v4 handshake response (defined in EIP-8).
type authRespV4 [T crypto.PrivateKey, P crypto.PublicKey]struct {
	RandomPubkey [pubLen]byte
	Nonce        [shaLen]byte
	Version      uint

	// Ignore additional fields (forward-compatibility)
	Rest []rlp.RawValue `rlp:"tail"`
}

// receiverEncHandshake negotiates a session token on conn.
// it should be called on the listening side of the connection.
//
// prv is the local client's private key.
func receiverEncHandshake[T crypto.PrivateKey, P crypto.PublicKey](conn io.ReadWriter, prv T) (s Secrets[T,P], err error) {
	authMsg := new(authMsgV4[T,P])
	authPacket, err := readHandshakeMsg[T,P](authMsg, encAuthMsgLen, prv, conn)
	if err != nil {
		return s, err
	}
	h := new(encHandshake[T,P])
	if err := h.handleAuthMsg(authMsg, prv); err != nil {
		return s, err
	}

	authRespMsg, err := h.makeAuthResp()
	if err != nil {
		return s, err
	}
	var authRespPacket []byte
	if authMsg.gotPlain {
		authRespPacket, err = authRespMsg.sealPlain(h)
	} else {
		authRespPacket, err = sealEIP8(authRespMsg, h)
	}
	if err != nil {
		return s, err
	}
	if _, err = conn.Write(authRespPacket); err != nil {
		return s, err
	}
	return h.secrets(authPacket, authRespPacket)
}

func (h *encHandshake[T,P]) handleAuthMsg(msg *authMsgV4[T,P], prv T) error {
	// Import the remote identity.
	rpub, err := importPublicKey[T,P](msg.InitiatorPubkey[:])
	if err != nil {
		return err
	}
	h.initNonce = msg.Nonce[:]
	h.remote = rpub

	// Generate random keypair for ECDH.
	// If a private key is already set, use it instead of generating one (for testing).
	if h.randomPrivKey == nil {
		h.randomPrivKey, err = ecies.GenerateKey[T,P](rand.Reader, crypto.S256(), nil)
		if err != nil {
			return err
		}
	}

	// Check the signature.
	token, err := h.staticSharedSecret(prv)
	if err != nil {
		return err
	}
	signedMsg := xor(token, h.initNonce)
	remoteRandomPub, err := crypto.Ecrecover[P](signedMsg, msg.Signature[:])
	if err != nil {
		return err
	}
	h.remoteRandomPub, _ = importPublicKey[T,P](remoteRandomPub)
	return nil
}

// secrets is called after the handshake is completed.
// It extracts the Connection secrets from the handshake values.
func (h *encHandshake[T,P]) secrets(auth, authResp []byte) (Secrets[T,P], error) {
	ecdheSecret, err := h.randomPrivKey.GenerateShared(h.remoteRandomPub, sskLen, sskLen)
	if err != nil {
		return Secrets[T,P]{}, err
	}

	// derive base secrets from ephemeral key agreement
	sharedSecret := crypto.Keccak256(ecdheSecret, crypto.Keccak256(h.respNonce, h.initNonce))
	aesSecret := crypto.Keccak256(ecdheSecret, sharedSecret)
	var remotePub P
	switch p:=any(&remotePub).(type){
	case *nist.PublicKey:
		*p = nist.PublicKey{&ecdsa.PublicKey{Curve: h.remote.Curve, X: h.remote.X, Y: h.remote.Y}}
	case *gost3410.PublicKey:
		*p = gost3410.PublicKey{C:h.remote.Curve, X: h.remote.X, Y: h.remote.Y}
	}
	s := Secrets[T,P]{
		remote: remotePub,
		AES:    aesSecret,
		MAC:    crypto.Keccak256(ecdheSecret, aesSecret),
	}

	// setup sha3 instances for the MACs
	mac1 := sha3.NewLegacyKeccak256()
	mac1.Write(xor(s.MAC, h.respNonce))
	mac1.Write(auth)
	mac2 := sha3.NewLegacyKeccak256()
	mac2.Write(xor(s.MAC, h.initNonce))
	mac2.Write(authResp)
	if h.initiator {
		s.EgressMAC, s.IngressMAC = mac1, mac2
	} else {
		s.EgressMAC, s.IngressMAC = mac2, mac1
	}

	return s, nil
}

// staticSharedSecret returns the static shared secret, the result
// of key agreement between the local and remote static node key.
func (h *encHandshake[T,P]) staticSharedSecret(prv T) ([]byte, error) {
	return ecies.ImportECDSA[T,P](prv).GenerateShared(h.remote, sskLen, sskLen)
}

// initiatorEncHandshake negotiates a session token on conn.
// it should be called on the dialing side of the connection.
//
// prv is the local client's private key.
func initiatorEncHandshake[T crypto.PrivateKey, P crypto.PublicKey](conn io.ReadWriter, prv T, remote P) (s Secrets[T,P], err error) {
	h := &encHandshake[T,P]{initiator: true, remote: ecies.ImportECDSAPublic(remote)}
	authMsg, err := h.makeAuthMsg(prv)
	if err != nil {
		return s, err
	}
	authPacket, err := sealEIP8(authMsg, h)
	if err != nil {
		return s, err
	}

	if _, err = conn.Write(authPacket); err != nil {
		return s, err
	}

	authRespMsg := new(authRespV4[T,P])
	authRespPacket, err := readHandshakeMsg[T,P](authRespMsg, encAuthRespLen, prv, conn)
	if err != nil {
		return s, err
	}
	if err := h.handleAuthResp(authRespMsg); err != nil {
		return s, err
	}
	return h.secrets(authPacket, authRespPacket)
}

// makeAuthMsg creates the initiator handshake message.
func (h *encHandshake[T,P]) makeAuthMsg(prv T) (*authMsgV4[T,P], error) {
	// Generate random initiator nonce.
	h.initNonce = make([]byte, shaLen)
	_, err := rand.Read(h.initNonce)
	if err != nil {
		return nil, err
	}
	// Generate random keypair to for ECDH.
	h.randomPrivKey, err = ecies.GenerateKey[T,P](rand.Reader, crypto.S256(), nil)
	if err != nil {
		return nil, err
	}

	// Sign known message: static-shared-secret ^ nonce
	token, err := h.staticSharedSecret(prv)
	if err != nil {
		return nil, err
	}
	signed := xor(token, h.initNonce)
	signature, err := crypto.Sign(signed, h.randomPrivKey.ExportECDSA())
	if err != nil {
		return nil, err
	}

	msg := new(authMsgV4[T,P])
	copy(msg.Signature[:], signature)
	var pub P
	switch t:=any(&prv).(type) {
	case *nist.PrivateKey:
		p := any(&pub).(*nist.PublicKey)
		*p = *t.Public()
	case *gost3410.PrivateKey:
		p := any(&pub).(*gost3410.PublicKey)
		*p = *t.Public()	
	}
	copy(msg.InitiatorPubkey[:], crypto.FromECDSAPub[P](pub)[1:])
	copy(msg.Nonce[:], h.initNonce)
	msg.Version = 4
	return msg, nil
}

func (h *encHandshake[T,P]) handleAuthResp(msg *authRespV4[T,P]) (err error) {
	h.respNonce = msg.Nonce[:]
	h.remoteRandomPub, err = importPublicKey[T,P](msg.RandomPubkey[:])
	return err
}

func (h *encHandshake[T,P]) makeAuthResp() (msg *authRespV4[T,P], err error) {
	// Generate random nonce.
	h.respNonce = make([]byte, shaLen)
	if _, err = rand.Read(h.respNonce); err != nil {
		return nil, err
	}

	msg = new(authRespV4[T,P])
	copy(msg.Nonce[:], h.respNonce)
	copy(msg.RandomPubkey[:], exportPubkey(&h.randomPrivKey.PublicKey))
	msg.Version = 4
	return msg, nil
}

func (msg *authMsgV4[T,P]) decodePlain(input []byte) {
	n := copy(msg.Signature[:], input)
	n += shaLen // skip sha3(initiator-ephemeral-pubk)
	n += copy(msg.InitiatorPubkey[:], input[n:])
	copy(msg.Nonce[:], input[n:])
	msg.Version = 4
	msg.gotPlain = true
}

func (msg *authRespV4[T,P]) sealPlain(hs *encHandshake[T,P]) ([]byte, error) {
	buf := make([]byte, authRespLen)
	n := copy(buf, msg.RandomPubkey[:])
	copy(buf[n:], msg.Nonce[:])
	return ecies.Encrypt[T](rand.Reader, hs.remote, buf, nil, nil)
}

func (msg *authRespV4[T,P]) decodePlain(input []byte) {
	n := copy(msg.RandomPubkey[:], input)
	copy(msg.Nonce[:], input[n:])
	msg.Version = 4
}

var padSpace = make([]byte, 300)

func sealEIP8[T crypto.PrivateKey, P crypto.PublicKey](msg interface{}, h *encHandshake[T,P]) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, msg); err != nil {
		return nil, err
	}
	// pad with random amount of data. the amount needs to be at least 100 bytes to make
	// the message distinguishable from pre-EIP-8 handshakes.
	pad := padSpace[:mrand.Intn(len(padSpace)-100)+100]
	buf.Write(pad)
	prefix := make([]byte, 2)
	binary.BigEndian.PutUint16(prefix, uint16(buf.Len()+eciesOverhead))

	enc, err := ecies.Encrypt[T](rand.Reader, h.remote, buf.Bytes(), nil, prefix)
	return append(prefix, enc...), err
}

type plainDecoder interface {
	decodePlain([]byte)
}

func readHandshakeMsg[T crypto.PrivateKey, P crypto.PublicKey](msg plainDecoder, plainSize int, prv T, r io.Reader) ([]byte, error) {
	buf := make([]byte, plainSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return buf, err
	}
	// Attempt decoding pre-EIP-8 "plain" format.
	key := ecies.ImportECDSA[T,P](prv)
	if dec, err := key.Decrypt(buf, nil, nil); err == nil {
		msg.decodePlain(dec)
		return buf, nil
	}
	// Could be EIP-8 format, try that.
	prefix := buf[:2]
	size := binary.BigEndian.Uint16(prefix)
	if size < uint16(plainSize) {
		return buf, fmt.Errorf("size underflow, need at least %d bytes", plainSize)
	}
	buf = append(buf, make([]byte, size-uint16(plainSize)+2)...)
	if _, err := io.ReadFull(r, buf[plainSize:]); err != nil {
		return buf, err
	}
	dec, err := key.Decrypt(buf[2:], nil, prefix)
	if err != nil {
		return buf, err
	}
	// Can't use rlp.DecodeBytes here because it rejects
	// trailing data (forward-compatibility).
	s := rlp.NewStream(bytes.NewReader(dec), 0)
	return buf, s.Decode(msg)
}

// importPublicKey unmarshals 512 bit public keys.
func importPublicKey[T crypto.PrivateKey, P crypto.PublicKey](pubKey []byte) (*ecies.PublicKey[P], error) {
	var pubKey65 []byte
	switch len(pubKey) {
	case 64:
		// add 'uncompressed key' flag
		pubKey65 = append([]byte{0x04}, pubKey...)
	case 65:
		pubKey65 = pubKey
	default:
		return nil, fmt.Errorf("invalid public key length %v (expect 64/65)", len(pubKey))
	}
	// TODO: fewer pointless conversions
	pub, err := crypto.UnmarshalPubkey[P](pubKey65)
	if err != nil {
		return nil, err
	}
	return ecies.ImportECDSAPublic(pub), nil
}

func exportPubkey[P crypto.PublicKey](pub *ecies.PublicKey[P]) []byte {
	if pub == nil {
		panic("nil pubkey")
	}
	return elliptic.Marshal(pub.Curve, pub.X, pub.Y)[1:]
}

func xor(one, other []byte) (xor []byte) {
	xor = make([]byte, len(one))
	for i := 0; i < len(one); i++ {
		xor[i] = one[i] ^ other[i]
	}
	return xor
}
