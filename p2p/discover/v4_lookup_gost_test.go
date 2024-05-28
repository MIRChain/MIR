// Copyright 2019 The go-ethereum Authors
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

package discover

import (
	"testing"

	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/p2p/enode"
)

func TestUDPv4_Lookup_Gost(t *testing.T) {
	t.Parallel()
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)

	// Lookup on empty table returns no nodes.
	targetKey, _ := decodePubkey[gost3410.PublicKey](lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]().target[:])
	if results := test.udp.LookupPubkey(targetKey); len(results) > 0 {
		t.Fatalf("lookup on empty table returned %d results: %#v", len(results), results)
	}

	// Seed table with initial node.
	fillTable(test.table, []*node[gost3410.PublicKey]{wrapNode(lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]().node(256, 0))})

	// Start the lookup.
	resultC := make(chan []*enode.Node[gost3410.PublicKey], 1)
	go func() {
		resultC <- test.udp.LookupPubkey(targetKey)
		test.close()
	}()

	// Answer lookup packets.
	serveTestnet(test, lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]())

	// Verify result nodes.
	results := <-resultC
	t.Logf("results:")
	for _, e := range results {
		t.Logf("  ld=%d, %x", enode.LogDist(lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]().target.id(), e.ID()), e.ID().Bytes())
	}
	if len(results) != bucketSize {
		t.Errorf("wrong number of results: got %d, want %d", len(results), bucketSize)
	}
	checkLookupResults(t, lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey](), results)
}

func TestUDPv4_LookupIterator_Gost(t *testing.T) {
	t.Parallel()
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)
	defer test.close()

	// Seed table with initial nodes.
	bootnodes := make([]*node[gost3410.PublicKey], len(lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]().dists[256]))
	for i := range lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]().dists[256] {
		bootnodes[i] = wrapNode(lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]().node(256, i))
	}
	fillTable(test.table, bootnodes)
	go serveTestnet(test, lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]())

	// Create the iterator and collect the nodes it yields.
	iter := test.udp.RandomNodes()
	seen := make(map[enode.ID]*enode.Node[gost3410.PublicKey])
	for limit := lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]().len(); iter.Next() && len(seen) < limit; {
		seen[iter.Node().ID()] = iter.Node()
	}
	iter.Close()

	// Check that all nodes in lookupTestnet were seen by the iterator.
	results := make([]*enode.Node[gost3410.PublicKey], 0, len(seen))
	for _, n := range seen {
		results = append(results, n)
	}
	sortByID(results)
	want := lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]().nodes()
	if err := checkNodesEqual(results, want); err != nil {
		t.Fatal(err)
	}
}

// TestUDPv4_LookupIteratorClose checks that lookupIterator ends when its Close
// method is called.
func TestUDPv4_LookupIteratorClose_Gost(t *testing.T) {
	t.Parallel()
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)
	defer test.close()

	// Seed table with initial nodes.
	bootnodes := make([]*node[gost3410.PublicKey], len(lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]().dists[256]))
	for i := range lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]().dists[256] {
		bootnodes[i] = wrapNode(lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]().node(256, i))
	}
	fillTable(test.table, bootnodes)
	go serveTestnet(test, lookupTestnetGost[gost3410.PrivateKey, gost3410.PublicKey]())

	it := test.udp.RandomNodes()
	if ok := it.Next(); !ok || it.Node() == nil {
		t.Fatalf("iterator didn't return any node")
	}

	it.Close()

	ncalls := 0
	for ; ncalls < 100 && it.Next(); ncalls++ {
		if it.Node() == nil {
			t.Error("iterator returned Node() == nil node after Next() == true")
		}
	}
	t.Logf("iterator returned %d nodes after close", ncalls)
	if it.Next() {
		t.Errorf("Next() == true after close and %d more calls", ncalls)
	}
	if n := it.Node(); n != nil {
		t.Errorf("iterator returned non-nil node after close and %d more calls", ncalls)
	}
}

// This is the test network for the Lookup test.
// The nodes were obtained by running lookupTestnet.mine with a random NodeID as target.
func lookupTestnetGost[T crypto.PrivateKey, P crypto.PublicKey]() *preminedTestnet[T, P] {
	return &preminedTestnet[T, P]{
		target: hexEncPubkey[P]("ba6303fbb832ef2a24a5a2b8603bc18981073dd7634b9c6454aa8183ccf07c3b60aabbca2fd8601a29958dd7d7799eb725051d25244635417b9f5a34f1c29a28"),
		dists: [257][]T{
			251: {
				hexEncPrivkey[T]("29738ba0c1a4397d6a65f292eee07f02df8e58d41594ba2be3cf84ce0fc58169"),
				hexEncPrivkey[T]("511b1686e4e58a917f7f848e9bf5539d206a68f5ad6b54b552c2399fe7d174ae"),
				hexEncPrivkey[T]("d09e5eaeec0fd596236faed210e55ef45112409a5aa7f3276d26646080dcfaeb"),
				hexEncPrivkey[T]("c1e20dbbf0d530e50573bd0a260b32ec15eb9190032b4633d44834afc8afe578"),
				hexEncPrivkey[T]("ed5f38f5702d92d306143e5d9154fb21819777da39af325ea359f453d179e80b"),
			},
			252: {
				hexEncPrivkey[T]("1c9b1cafbec00848d2c174b858219914b42a7d5c9359b1ca03fd650e8239ae94"),
				hexEncPrivkey[T]("e0e1e8db4a6f13c1ffdd3e96b72fa7012293ced187c9dcdcb9ba2af37a46fa10"),
				hexEncPrivkey[T]("3d53823e0a0295cb09f3e11d16c1b44d07dd37cec6f739b8df3a590189fe9fb9"),
			},
			253: {
				hexEncPrivkey[T]("2d0511ae9bf590166597eeab86b6f27b1ab761761eaea8965487b162f8703847"),
				hexEncPrivkey[T]("6cfbd7b8503073fc3dbdb746a7c672571648d3bd15197ccf7f7fef3d904f53a2"),
				hexEncPrivkey[T]("a30599b12827b69120633f15b98a7f6bc9fc2e9a0fd6ae2ebb767c0e64d743ab"),
				hexEncPrivkey[T]("14a98db9b46a831d67eff29f3b85b1b485bb12ae9796aea98d91be3dc78d8a91"),
				hexEncPrivkey[T]("2369ff1fc1ff8ca7d20b17e2673adc3365c3674377f21c5d9dafaff21fe12e24"),
				hexEncPrivkey[T]("9ae91101d6b5048607f41ec0f690ef5d09507928aded2410aabd9237aa2727d7"),
				hexEncPrivkey[T]("05e3c59090a3fd1ae697c09c574a36fcf9bedd0afa8fe3946f21117319ca4973"),
				hexEncPrivkey[T]("06f31c5ea632658f718a91a1b1b9ae4b7549d7b3bc61cbc2be5f4a439039f3ad"),
			},
			254: {
				hexEncPrivkey[T]("dec742079ec00ff4ec1284d7905bc3de2366f67a0769431fd16f80fd68c58a7c"),
				hexEncPrivkey[T]("ff02c8861fa12fbd129d2a95ea663492ef9c1e51de19dcfbbfe1c59894a28d2b"),
				hexEncPrivkey[T]("4dded9e4eefcbce4262be4fd9e8a773670ab0b5f448f286ec97dfc8cf681444a"),
				hexEncPrivkey[T]("750d931e2a8baa2c9268cb46b7cd851f4198018bed22f4dceb09dd334a2395f6"),
				hexEncPrivkey[T]("ce1435a956a98ffec484cd11489c4f165cf1606819ab6b521cee440f0c677e9e"),
				hexEncPrivkey[T]("996e7f8d1638be92d7328b4770f47e5420fc4bafecb4324fd33b1f5d9f403a75"),
				hexEncPrivkey[T]("ebdc44e77a6cc0eb622e58cf3bb903c3da4c91ca75b447b0168505d8fc308b9c"),
				hexEncPrivkey[T]("46bd1eddcf6431bea66fc19ebc45df191c1c7d6ed552dcdc7392885009c322f0"),
			},
			255: {
				hexEncPrivkey[T]("da8645f90826e57228d9ea72aff84500060ad111a5d62e4af831ed8e4b5acfb8"),
				hexEncPrivkey[T]("3c944c5d9af51d4c1d43f5d0f3a1a7ef65d5e82744d669b58b5fed242941a566"),
				hexEncPrivkey[T]("5ebcde76f1d579eebf6e43b0ffe9157e65ffaa391175d5b9aa988f47df3e33da"),
				hexEncPrivkey[T]("97f78253a7d1d796e4eaabce721febcc4550dd68fb11cc818378ba807a2cb7de"),
				hexEncPrivkey[T]("a38cd7dc9b4079d1c0406afd0fdb1165c285f2c44f946eca96fc67772c988c7d"),
				hexEncPrivkey[T]("d64cbb3ffdf712c372b7a22a176308ef8f91861398d5dbaf326fd89c6eaeef1c"),
				hexEncPrivkey[T]("d269609743ef29d6446e3355ec647e38d919c82a4eb5837e442efd7f4218944f"),
				hexEncPrivkey[T]("d8f7bcc4a530efde1d143717007179e0d9ace405ddaaf151c4d863753b7fd64c"),
			},
			256: {
				hexEncPrivkey[T]("8c5b422155d33ea8e9d46f71d1ad3e7b24cb40051413ffa1a81cff613d243ba9"),
				hexEncPrivkey[T]("937b1af801def4e8f5a3a8bd225a8bcff1db764e41d3e177f2e9376e8dd87233"),
				hexEncPrivkey[T]("120260dce739b6f71f171da6f65bc361b5fad51db74cf02d3e973347819a6518"),
				hexEncPrivkey[T]("1fa56cf25d4b46c2bf94e82355aa631717b63190785ac6bae545a88aadc304a9"),
				hexEncPrivkey[T]("3c38c503c0376f9b4adcbe935d5f4b890391741c764f61b03cd4d0d42deae002"),
				hexEncPrivkey[T]("3a54af3e9fa162bc8623cdf3e5d9b70bf30ade1d54cc3abea8659aba6cff471f"),
				hexEncPrivkey[T]("6799a02ea1999aefdcbcc4d3ff9544478be7365a328d0d0f37c26bd95ade0cda"),
				hexEncPrivkey[T]("e24a7bc9051058f918646b0f6e3d16884b2a55a15553b89bab910d55ebc36116"),
			},
		},
	}
}
