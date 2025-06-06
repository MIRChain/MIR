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

package common

import (
	"bytes"
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"reflect"
	"strings"

	"github.com/pavelkrolevets/MIR-pro/common/hexutil"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"golang.org/x/crypto/sha3"
)

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 20
	// length of the hash returned by Private Transaction Manager
	EncryptedPayloadHashLength = 64
)

var (
	ErrNotPrivateContract = errors.New("the provided address is not a private contract")
	ErrNoAccountExtraData = errors.New("no account extra data found")

	hashT    = reflect.TypeOf(Hash{})
	addressT = reflect.TypeOf(Address{})
)

// Hash, returned by Private Transaction Manager, represents the 64-byte hash of encrypted payload
type EncryptedPayloadHash [EncryptedPayloadHashLength]byte

// Using map to enable fast lookup
type EncryptedPayloadHashes map[EncryptedPayloadHash]struct{}

func (h *EncryptedPayloadHash) MarshalJSON() (j []byte, err error) {
	return json.Marshal(h.ToBase64())
}

func (h *EncryptedPayloadHash) UnmarshalJSON(j []byte) (err error) {
	var ephStr string
	err = json.Unmarshal(j, &ephStr)
	if err != nil {
		return err
	}
	eph, err := Base64ToEncryptedPayloadHash(ephStr)
	if err != nil {
		return err
	}
	h.SetBytes(eph.Bytes())
	return nil
}

func (h *EncryptedPayloadHashes) MarshalJSON() (j []byte, err error) {
	return json.Marshal(h.ToBase64s())
}

func (h *EncryptedPayloadHashes) UnmarshalJSON(j []byte) (err error) {
	var ephStrArray []string
	err = json.Unmarshal(j, &ephStrArray)
	if err != nil {
		return err
	}
	for _, str := range ephStrArray {
		eph, err := Base64ToEncryptedPayloadHash(str)
		if err != nil {
			return err
		}
		h.Add(eph)
	}
	return nil
}

// BytesToEncryptedPayloadHash sets b to EncryptedPayloadHash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToEncryptedPayloadHash(b []byte) EncryptedPayloadHash {
	var h EncryptedPayloadHash
	h.SetBytes(b)
	return h
}

func Base64ToEncryptedPayloadHash(b64 string) (EncryptedPayloadHash, error) {
	bytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return EncryptedPayloadHash{}, fmt.Errorf("unable to convert base64 string %s to EncryptedPayloadHash. Cause: %v", b64, err)
	}
	return BytesToEncryptedPayloadHash(bytes), nil
}

func (eph *EncryptedPayloadHash) SetBytes(b []byte) {
	if len(b) > len(eph) {
		b = b[len(b)-EncryptedPayloadHashLength:]
	}

	copy(eph[EncryptedPayloadHashLength-len(b):], b)
}

func (eph EncryptedPayloadHash) Hex() string {
	return hexutil.Encode(eph[:])
}

func (eph EncryptedPayloadHash) Bytes() []byte {
	return eph[:]
}

func (eph EncryptedPayloadHash) String() string {
	return eph.Hex()
}

func (eph EncryptedPayloadHash) ToBase64() string {
	return base64.StdEncoding.EncodeToString(eph[:])
}

func (eph EncryptedPayloadHash) TerminalString() string {
	return fmt.Sprintf("%x…%x", eph[:3], eph[EncryptedPayloadHashLength-3:])
}

func (eph EncryptedPayloadHash) BytesTypeRef() *hexutil.Bytes {
	b := hexutil.Bytes(eph.Bytes())
	return &b
}

func EmptyEncryptedPayloadHash(eph EncryptedPayloadHash) bool {
	return eph == EncryptedPayloadHash{}
}

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

func StringToHash(s string) Hash { return BytesToHash([]byte(s)) } // dep: Istanbul

// BigToHash sets byte representation of b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BigToHash(b *big.Int) Hash { return BytesToHash(b.Bytes()) }

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte { return h[:] }

// Big converts a hash to a big integer.
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex converts a hash to a hex string.
func (h Hash) Hex() string { return hexutil.Encode(h[:]) }

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (h Hash) TerminalString() string {
	return fmt.Sprintf("%x..%x", h[:3], h[29:])
}

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h Hash) String() string {
	return h.Hex()
}

// Format implements fmt.Formatter.
// Hash supports the %v, %s, %v, %x, %X and %d format verbs.
func (h Hash) Format(s fmt.State, c rune) {
	hexb := make([]byte, 2+len(h)*2)
	copy(hexb, "0x")
	hex.Encode(hexb[2:], h[:])

	switch c {
	case 'x', 'X':
		if !s.Flag('#') {
			hexb = hexb[2:]
		}
		if c == 'X' {
			hexb = bytes.ToUpper(hexb)
		}
		fallthrough
	case 'v', 's':
		s.Write(hexb)
	case 'q':
		q := []byte{'"'}
		s.Write(q)
		s.Write(hexb)
		s.Write(q)
	case 'd':
		fmt.Fprint(s, ([len(h)]byte)(h))
	default:
		fmt.Fprintf(s, "%%!%c(hash=%x)", c, h)
	}
}

// UnmarshalText parses a hash in hex syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Hash", input, h[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(hashT, input, h[:])
}

// MarshalText returns the hex representation of h.
func (h Hash) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

func EmptyHash(h Hash) bool {
	return h == Hash{}
}

// Generate implements testing/quick.Generator.
func (h Hash) Generate(rand *rand.Rand, size int) reflect.Value {
	m := rand.Intn(len(h))
	for i := len(h) - 1; i > m; i-- {
		h[i] = byte(rand.Uint32())
	}
	return reflect.ValueOf(h)
}

func (h Hash) ToBase64() string {
	return base64.StdEncoding.EncodeToString(h.Bytes())
}

// Decode base64 string to Hash
// if String is empty then return empty hash
func Base64ToHash(b64 string) (Hash, error) {
	if b64 == "" {
		return Hash{}, nil
	}
	bytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return Hash{}, fmt.Errorf("unable to convert base64 string %s to Hash. Cause: %v", b64, err)
	}
	return BytesToHash(bytes), nil
}

// Scan implements Scanner for database/sql.
func (h *Hash) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Hash", src)
	}
	if len(srcB) != HashLength {
		return fmt.Errorf("can't scan []byte of len %d into Hash, want %d", len(srcB), HashLength)
	}
	copy(h[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (h Hash) Value() (driver.Value, error) {
	return h[:], nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (Hash) ImplementsGraphQLType(name string) bool { return name == "Bytes32" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (h *Hash) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		err = h.UnmarshalText([]byte(input))
	default:
		err = fmt.Errorf("unexpected type %T for Hash", input)
	}
	return err
}

// UnprefixedHash allows marshaling a Hash without 0x prefix.
type UnprefixedHash Hash

// UnmarshalText decodes the hash from hex. The 0x prefix is optional.
func (h *UnprefixedHash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedHash", input, h[:])
}

// MarshalText encodes the hash as hex.
func (h UnprefixedHash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

func (ephs EncryptedPayloadHashes) ToBase64s() []string {
	a := make([]string, 0, len(ephs))
	for eph := range ephs {
		a = append(a, eph.ToBase64())
	}
	return a
}

func (ephs EncryptedPayloadHashes) NotExist(eph EncryptedPayloadHash) bool {
	_, ok := ephs[eph]
	return !ok
}

func (ephs EncryptedPayloadHashes) Add(eph EncryptedPayloadHash) {
	ephs[eph] = struct{}{}
}

func (ephs EncryptedPayloadHashes) EncodeRLP(writer io.Writer) error {
	encryptedPayloadHashesArray := make([]EncryptedPayloadHash, len(ephs))
	idx := 0
	for key := range ephs {
		encryptedPayloadHashesArray[idx] = key
		idx++
	}
	return rlp.Encode(writer, encryptedPayloadHashesArray)
}

func (ephs EncryptedPayloadHashes) DecodeRLP(stream *rlp.Stream) error {
	var encryptedPayloadHashesRLP []EncryptedPayloadHash
	if err := stream.Decode(&encryptedPayloadHashesRLP); err != nil {
		return err
	}
	for _, val := range encryptedPayloadHashesRLP {
		ephs.Add(val)
	}
	return nil
}

func Base64sToEncryptedPayloadHashes(b64s []string) (EncryptedPayloadHashes, error) {
	ephs := make(EncryptedPayloadHashes)
	for _, b64 := range b64s {
		data, err := Base64ToEncryptedPayloadHash(b64)
		if err != nil {
			return nil, err
		}
		ephs.Add(data)
	}
	return ephs, nil
}

// Print hex but only first 3 and last 3 bytes
func FormatTerminalString(data []byte) string {
	l := len(data)
	if l > 0 {
		if l > 6 {
			return fmt.Sprintf("%x…%x", data[:3], data[l-3:])
		} else {
			return fmt.Sprintf("%x", data[:])
		}
	}
	return ""
}

/////////// Address

// Address represents the 20 byte address of an Ethereum account.
type Address [AddressLength]byte

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

func StringToAddress(s string) Address { return BytesToAddress([]byte(s)) } // dep: Istanbul

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

// IsHexAddress verifies whether a string can represent a valid hex-encoded
// Ethereum address or not.
func IsHexAddress(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

// Bytes gets the string representation of the underlying address.
func (a Address) Bytes() []byte { return a[:] }

// Hash converts an address to a hash by left-padding it with zeros.
func (a Address) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns an EIP55-compliant hex string representation of the address.
func (a Address) Hex() string {
	return string(a.checksumHex())
}

// String implements fmt.Stringer.
func (a Address) String() string {
	return a.Hex()
}

func (a *Address) checksumHex() []byte {
	buf := a.hex()

	// compute checksum
	sha := sha3.NewLegacyKeccak256()
	sha.Write(buf[2:])
	hash := sha.Sum(nil)
	for i := 2; i < len(buf); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	return buf[:]
}

func (a Address) hex() []byte {
	var buf [len(a)*2 + 2]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], a[:])
	return buf[:]
}

// Format implements fmt.Formatter.
// Address supports the %v, %s, %v, %x, %X and %d format verbs.
func (a Address) Format(s fmt.State, c rune) {
	switch c {
	case 'v', 's':
		s.Write(a.checksumHex())
	case 'q':
		q := []byte{'"'}
		s.Write(q)
		s.Write(a.checksumHex())
		s.Write(q)
	case 'x', 'X':
		// %x disables the checksum.
		hex := a.hex()
		if !s.Flag('#') {
			hex = hex[2:]
		}
		if c == 'X' {
			hex = bytes.ToUpper(hex)
		}
		s.Write(hex)
	case 'd':
		fmt.Fprint(s, ([len(a)]byte)(a))
	default:
		fmt.Fprintf(s, "%%!%c(address=%x)", c, a)
	}
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a), b will be cropped from the left.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a Address) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *Address) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Address", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Address) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}

// Scan implements Scanner for database/sql.
func (a *Address) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Address", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into Address, want %d", len(srcB), AddressLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a Address) Value() (driver.Value, error) {
	return a[:], nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (a Address) ImplementsGraphQLType(name string) bool { return name == "Address" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (a *Address) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		err = a.UnmarshalText([]byte(input))
	default:
		err = fmt.Errorf("unexpected type %T for Address", input)
	}
	return err
}

// UnprefixedAddress allows marshaling an Address without 0x prefix.
type UnprefixedAddress Address

// UnmarshalText decodes the address from hex. The 0x prefix is optional.
func (a *UnprefixedAddress) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedAddress", input, a[:])
}

// MarshalText encodes the address as hex.
func (a UnprefixedAddress) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(a[:])), nil
}

// MixedcaseAddress retains the original string, which may or may not be
// correctly checksummed
type MixedcaseAddress struct {
	addr     Address
	original string
}

// NewMixedcaseAddress constructor (mainly for testing)
func NewMixedcaseAddress(addr Address) MixedcaseAddress {
	return MixedcaseAddress{addr: addr, original: addr.Hex()}
}

// NewMixedcaseAddressFromString is mainly meant for unit-testing
func NewMixedcaseAddressFromString(hexaddr string) (*MixedcaseAddress, error) {
	if !IsHexAddress(hexaddr) {
		return nil, errors.New("invalid address")
	}
	a := FromHex(hexaddr)
	return &MixedcaseAddress{addr: BytesToAddress(a), original: hexaddr}, nil
}

// UnmarshalJSON parses MixedcaseAddress
func (ma *MixedcaseAddress) UnmarshalJSON(input []byte) error {
	if err := hexutil.UnmarshalFixedJSON(addressT, input, ma.addr[:]); err != nil {
		return err
	}
	return json.Unmarshal(input, &ma.original)
}

// MarshalJSON marshals the original value
func (ma *MixedcaseAddress) MarshalJSON() ([]byte, error) {
	if strings.HasPrefix(ma.original, "0x") || strings.HasPrefix(ma.original, "0X") {
		return json.Marshal(fmt.Sprintf("0x%s", ma.original[2:]))
	}
	return json.Marshal(fmt.Sprintf("0x%s", ma.original))
}

// Address returns the address
func (ma *MixedcaseAddress) Address() Address {
	return ma.addr
}

// String implements fmt.Stringer
func (ma *MixedcaseAddress) String() string {
	if ma.ValidChecksum() {
		return fmt.Sprintf("%s [chksum ok]", ma.original)
	}
	return fmt.Sprintf("%s [chksum INVALID]", ma.original)
}

// ValidChecksum returns true if the address has valid checksum
func (ma *MixedcaseAddress) ValidChecksum() bool {
	return ma.original == ma.addr.Hex()
}

// Original returns the mixed-case input string
func (ma *MixedcaseAddress) Original() string {
	return ma.original
}

type DecryptRequest struct {
	SenderKey       []byte   `json:"senderKey"`
	CipherText      []byte   `json:"cipherText"`
	CipherTextNonce []byte   `json:"cipherTextNonce"`
	RecipientBoxes  []string `json:"recipientBoxes"`
	RecipientNonce  []byte   `json:"recipientNonce"`
	RecipientKeys   []string `json:"recipientKeys"`
}
