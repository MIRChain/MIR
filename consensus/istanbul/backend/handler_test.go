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

package backend

import (
	"bytes"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus/istanbul"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

func TestIstanbulMessage(t *testing.T) {
	_, backend := NewBlockChain(1, nil)
	defer backend.Stop()

	// generate one msg
	data := []byte("data1")
	hash := istanbul.RLPHash(data)
	msg := makeMsg(istanbulMsg, data)
	addr := common.StringToAddress("address")

	// 1. this message should not be in cache
	// for peers
	if _, ok := backend.recentMessages.Get(addr); ok {
		t.Fatalf("the cache of messages for this peer should be nil")
	}

	// for self
	if _, ok := backend.knownMessages.Get(hash); ok {
		t.Fatalf("the cache of messages should be nil")
	}

	// 2. this message should be in cache after we handle it
	_, err := backend.HandleMsg(addr, msg)
	if err != nil {
		t.Fatalf("handle message failed: %v", err)
	}
	// for peers
	if ms, ok := backend.recentMessages.Get(addr); ms == nil || !ok {
		t.Fatalf("the cache of messages for this peer cannot be nil")
	} else if m, ok := ms.(*lru.ARCCache); !ok {
		t.Fatalf("the cache of messages for this peer cannot be casted")
	} else if _, ok := m.Get(hash); !ok {
		t.Fatalf("the cache of messages for this peer cannot be found")
	}

	// for self
	if _, ok := backend.knownMessages.Get(hash); !ok {
		t.Fatalf("the cache of messages cannot be found")
	}
}

func makeMsg(msgcode uint64, data interface{}) p2p.Msg {
	size, r, _ := rlp.EncodeToReader(data)
	return p2p.Msg{Code: msgcode, Size: uint32(size), Payload: r}
}

func tryUntilMessageIsHandled(backend *Backend[nist.PrivateKey,nist.PublicKey], arbitraryAddress common.Address, arbitraryP2PMessage p2p.Msg) (handled bool, err error) {
	for i := 0; i < 5; i++ { // make 5 tries if a little wait
		handled, err = backend.HandleMsg(arbitraryAddress, arbitraryP2PMessage)
		if handled && err == nil {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	return
}

func TestHandleNewBlockMessage_whenTypical(t *testing.T) {
	_, backend := NewBlockChain(1, nil)
	defer backend.Stop()
	arbitraryAddress := common.StringToAddress("arbitrary")
	arbitraryBlock, arbitraryP2PMessage := buildArbitraryP2PNewBlockMessage[nist.PublicKey](t, false)
	postAndWait(backend, arbitraryBlock, t)

	handled, err := tryUntilMessageIsHandled(backend, arbitraryAddress, arbitraryP2PMessage)
	if err != nil {
		t.Errorf("expected message being handled successfully but got %s", err)
	}
	if !handled {
		t.Errorf("expected message being handled but not")
	}
	if _, err := ioutil.ReadAll(arbitraryP2PMessage.Payload); err != nil {
		t.Errorf("expected p2p message payload is restored")
	}
}

func TestHandleNewBlockMessage_whenNotAProposedBlock(t *testing.T) {
	_, backend := NewBlockChain(1, nil)
	defer backend.Stop()
	arbitraryAddress := common.StringToAddress("arbitrary")
	_, arbitraryP2PMessage := buildArbitraryP2PNewBlockMessage[nist.PublicKey](t, false)
	postAndWait(backend, types.NewBlock[nist.PublicKey](&types.Header[nist.PublicKey]{
		Number:    big.NewInt(1),
		Root:      common.StringToHash("someroot"),
		GasLimit:  1,
		MixDigest: types.IstanbulDigest,
	}, nil, nil, nil, new(trie.Trie[nist.PublicKey])), t)

	handled, err := tryUntilMessageIsHandled(backend, arbitraryAddress, arbitraryP2PMessage)
	if err != nil {
		t.Errorf("expected message being handled successfully but got %s", err)
	}
	if handled {
		t.Errorf("expected message not being handled")
	}
	if _, err := ioutil.ReadAll(arbitraryP2PMessage.Payload); err != nil {
		t.Errorf("expected p2p message payload is restored")
	}
}

func TestHandleNewBlockMessage_whenFailToDecode(t *testing.T) {
	_, backend := NewBlockChain(1, nil)
	defer backend.Stop()
	arbitraryAddress := common.StringToAddress("arbitrary")
	_, arbitraryP2PMessage := buildArbitraryP2PNewBlockMessage[nist.PublicKey](t, true)
	postAndWait(backend, types.NewBlock[nist.PublicKey](&types.Header[nist.PublicKey]{
		Number:    big.NewInt(1),
		GasLimit:  1,
		MixDigest: types.IstanbulDigest,
	}, nil, nil, nil, new(trie.Trie[nist.PublicKey])), t)

	handled, err := tryUntilMessageIsHandled(backend, arbitraryAddress, arbitraryP2PMessage)
	if err != nil {
		t.Errorf("expected message being handled successfully but got %s", err)
	}
	if handled {
		t.Errorf("expected message not being handled")
	}
	if _, err := ioutil.ReadAll(arbitraryP2PMessage.Payload); err != nil {
		t.Errorf("expected p2p message payload is restored")
	}
}

func postAndWait[T crypto.PrivateKey,P crypto.PublicKey](backend *Backend[T,P], block *types.Block[P], t *testing.T) {
	eventSub := backend.EventMux().Subscribe(istanbul.RequestEvent{})
	defer eventSub.Unsubscribe()
	stop := make(chan struct{}, 1)
	eventLoop := func() {
		<-eventSub.Chan()
		stop <- struct{}{}
	}
	go eventLoop()
	if err := backend.EventMux().Post(istanbul.RequestEvent{
		Proposal: block,
	}); err != nil {
		t.Fatalf("%s", err)
	}
	<-stop
}

func buildArbitraryP2PNewBlockMessage[P crypto.PublicKey](t *testing.T, invalidMsg bool) (*types.Block[P], p2p.Msg) {
	arbitraryBlock := types.NewBlock[P](&types.Header[P]{
		Number:    big.NewInt(1),
		GasLimit:  0,
		MixDigest: types.IstanbulDigest,
	}, nil, nil, nil, new(trie.Trie[P]))
	request := []interface{}{&arbitraryBlock, big.NewInt(1)}
	if invalidMsg {
		request = []interface{}{"invalid msg"}
	}
	size, r, err := rlp.EncodeToReader(request)
	if err != nil {
		t.Fatalf("can't encode due to %s", err)
	}
	payload, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("can't read payload due to %s", err)
	}
	arbitraryP2PMessage := p2p.Msg{Code: 0x07, Size: uint32(size), Payload: bytes.NewReader(payload)}
	return arbitraryBlock, arbitraryP2PMessage
}
