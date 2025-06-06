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

package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jpmorganchase/quorum-security-plugin-sdk-go/proto"
	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/bloombits"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/eth"
	"github.com/pavelkrolevets/MIR-pro/eth/downloader"
	"github.com/pavelkrolevets/MIR-pro/eth/ethconfig"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/event"
	"github.com/pavelkrolevets/MIR-pro/internal/ethapi"
	"github.com/pavelkrolevets/MIR-pro/multitenancy"
	"github.com/pavelkrolevets/MIR-pro/node"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/private"
	"github.com/pavelkrolevets/MIR-pro/private/engine"
	"github.com/pavelkrolevets/MIR-pro/private/engine/notinuse"
	"github.com/pavelkrolevets/MIR-pro/rpc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
)

func TestBuildSchema(t *testing.T) {
	ddir, err := ioutil.TempDir("", "graphql-buildschema")
	if err != nil {
		t.Fatalf("failed to create temporary datadir: %v", err)
	}
	// Copy config
	conf := node.DefaultConfig
	conf.DataDir = ddir
	stack, err := node.New(&conf)
	if err != nil {
		t.Fatalf("could not create new node: %v", err)
	}
	// Make sure the schema can be parsed and matched up to the object model.
	if err := newHandler(stack, nil, []string{}, []string{}); err != nil {
		t.Errorf("Could not construct GraphQL handler: %v", err)
	}
}

// Tests that a graphQL request is successfully handled when graphql is enabled on the specified endpoint
func TestGraphQLBlockSerialization(t *testing.T) {
	stack := createNode[nist.PrivateKey,nist.PublicKey](t, true, false)
	defer stack.Close()
	// start node
	if err := stack.Start(); err != nil {
		t.Fatalf("could not start node: %v", err)
	}

	for i, tt := range []struct {
		body string
		want string
		code int
	}{
		{ // Should return latest block
			body: `{"query": "{block{number}}","variables": null}`,
			want: `{"data":{"block":{"number":10}}}`,
			code: 200,
		},
		{ // Should return info about latest block
			body: `{"query": "{block{number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":{"number":10,"gasUsed":0,"gasLimit":11500000}}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:0){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":{"number":0,"gasUsed":0,"gasLimit":11500000}}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:-1){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":null}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:-500){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":null}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:\"0\"){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":{"number":0,"gasUsed":0,"gasLimit":11500000}}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:\"-33\"){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":null}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:\"1337\"){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"data":{"block":null}}`,
			code: 200,
		},
		{
			body: `{"query": "{block(number:\"0xbad\"){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"errors":[{"message":"strconv.ParseInt: parsing \"0xbad\": invalid syntax"}],"data":{}}`,
			code: 400,
		},
		{ // hex strings are currently not supported. If that's added to the spec, this test will need to change
			body: `{"query": "{block(number:\"0x0\"){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"errors":[{"message":"strconv.ParseInt: parsing \"0x0\": invalid syntax"}],"data":{}}`,
			code: 400,
		},
		{
			body: `{"query": "{block(number:\"a\"){number,gasUsed,gasLimit}}","variables": null}`,
			want: `{"errors":[{"message":"strconv.ParseInt: parsing \"a\": invalid syntax"}],"data":{}}`,
			code: 400,
		},
		{
			body: `{"query": "{bleh{number}}","variables": null}"`,
			want: `{"errors":[{"message":"Cannot query field \"bleh\" on type \"Query\".","locations":[{"line":1,"column":2}]}]}`,
			code: 400,
		},
		// should return `estimateGas` as decimal
		{
			body: `{"query": "{block{ estimateGas(data:{}) }}"}`,
			want: `{"data":{"block":{"estimateGas":53000}}}`,
			code: 200,
		},
		// should return `status` as decimal
		{
			body: `{"query": "{block {number call (data : {from : \"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b\", to: \"0x6295ee1b4f6dd65047762f924ecd367c17eabf8f\", data :\"0x12a7b914\"}){data status}}}"}`,
			want: `{"data":{"block":{"number":10,"call":{"data":"0x","status":1}}}}`,
			code: 200,
		},
	} {
		resp, err := http.Post(fmt.Sprintf("%s/graphql", stack.HTTPEndpoint()), "application/json", strings.NewReader(tt.body))
		if err != nil {
			t.Fatalf("could not post: %v", err)
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read from response body: %v", err)
		}
		if have := string(bodyBytes); have != tt.want {
			t.Errorf("testcase %d %s,\nhave:\n%v\nwant:\n%v", i, tt.body, have, tt.want)
		}
		if tt.code != resp.StatusCode {
			t.Errorf("testcase %d %s,\nwrong statuscode, have: %v, want: %v", i, tt.body, resp.StatusCode, tt.code)
		}
	}
}

func TestGraphQLBlockSerializationEIP2718(t *testing.T) {
	stack := createNode[nist.PrivateKey,nist.PublicKey](t, true, true)
	defer stack.Close()
	// start node
	if err := stack.Start(); err != nil {
		t.Fatalf("could not start node: %v", err)
	}

	for i, tt := range []struct {
		body string
		want string
		code int
	}{
		{
			body: `{"query": "{block {number transactions { from { address } to { address } value hash type accessList { address storageKeys } index}}}"}`,
			want: `{"data":{"block":{"number":1,"transactions":[{"from":{"address":"0x71562b71999873db5b286df957af199ec94617f7"},"to":{"address":"0x0000000000000000000000000000000000000dad"},"value":"0x64","hash":"0x4f7b8d718145233dcf7f29e34a969c63dd4de8715c054ea2af022b66c4f4633e","type":0,"accessList":[],"index":0},{"from":{"address":"0x71562b71999873db5b286df957af199ec94617f7"},"to":{"address":"0x0000000000000000000000000000000000000dad"},"value":"0x32","hash":"0x9c6c2c045b618fe87add0e49ba3ca00659076ecae00fd51de3ba5d4ccf9dbf40","type":1,"accessList":[{"address":"0x0000000000000000000000000000000000000dad","storageKeys":["0x0000000000000000000000000000000000000000000000000000000000000000"]}],"index":1}]}}}`,
			code: 200,
		},
	} {
		resp, err := http.Post(fmt.Sprintf("%s/graphql", stack.HTTPEndpoint()), "application/json", strings.NewReader(tt.body))
		if err != nil {
			t.Fatalf("could not post: %v", err)
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("could not read from response body: %v", err)
		}
		if have := string(bodyBytes); have != tt.want {
			t.Errorf("testcase %d %s,\nhave:\n%v\nwant:\n%v", i, tt.body, have, tt.want)
		}
		if tt.code != resp.StatusCode {
			t.Errorf("testcase %d %s,\nwrong statuscode, have: %v, want: %v", i, tt.body, resp.StatusCode, tt.code)
		}
	}
}

// Tests that a graphQL request is not handled successfully when graphql is not enabled on the specified endpoint
func TestGraphQLHTTPOnSamePort_GQLRequest_Unsuccessful(t *testing.T) {
	stack := createNode[nist.PrivateKey,nist.PublicKey](t, false, false)
	defer stack.Close()
	if err := stack.Start(); err != nil {
		t.Fatalf("could not start node: %v", err)
	}
	body := strings.NewReader(`{"query": "{block{number}}","variables": null}`)
	resp, err := http.Post(fmt.Sprintf("%s/graphql", stack.HTTPEndpoint()), "application/json", body)
	if err != nil {
		t.Fatalf("could not post: %v", err)
	}
	// make sure the request is not handled successfully
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func createNode[T crypto.PrivateKey, P crypto.PublicKey](t *testing.T, gqlEnabled bool, txEnabled bool) *node.Node[T,P] {
	stack, err := node.New(&node.Config[T,P]{
		HTTPHost: "127.0.0.1",
		HTTPPort: 0,
		WSHost:   "127.0.0.1",
		WSPort:   0,
	})
	if err != nil {
		t.Fatalf("could not create node: %v", err)
	}
	if !gqlEnabled {
		return stack
	}
	if !txEnabled {
		createGQLService(t, stack)
	} else {
		createGQLServiceWithTransactions(t, stack)
	}
	return stack
}

func createGQLService[T crypto.PrivateKey, P crypto.PublicKey](t *testing.T, stack *node.Node[T,P]) {
	// create backend
	ethConf := &ethconfig.Config[P]{
		Genesis: &core.Genesis[P]{
			Config:     params.AllEthashProtocolChanges,
			GasLimit:   11500000,
			Difficulty: big.NewInt(1048576),
		},
		Ethash: ethash.Config{
			PowMode: ethash.ModeFake,
		},
		NetworkId:               1337,
		TrieCleanCache:          5,
		TrieCleanCacheJournal:   "triecache",
		TrieCleanCacheRejournal: 60 * time.Minute,
		TrieDirtyCache:          5,
		TrieTimeout:             60 * time.Minute,
		SnapshotCache:           5,
	}
	ethBackend, err := eth.New[T,P](stack, ethConf)
	if err != nil {
		t.Fatalf("could not create eth backend: %v", err)
	}
	// Create some blocks and import them
	chain, _ := core.GenerateChain[P](params.AllEthashProtocolChanges, ethBackend.BlockChain().Genesis(),
		 ethash.NewFaker[P](), ethBackend.ChainDb(), 10, func(i int, gen *core.BlockGen[P]) {})
	_, err = ethBackend.BlockChain().InsertChain(chain)
	if err != nil {
		t.Fatalf("could not create import blocks: %v", err)
	}
	// create gql service
	err = New[T,P](stack, ethBackend.APIBackend, []string{}, []string{})
	if err != nil {
		t.Fatalf("could not create graphql service: %v", err)
	}
}

func createGQLServiceWithTransactions[T crypto.PrivateKey, P crypto.PublicKey](t *testing.T, stack *node.Node[T,P]) {
	// create backend
	key, _ := crypto.HexToECDSA[T]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	var pub P
	switch t:=any(&key).(type){
	case *nist.PrivateKey:
		p:=any(&pub).(*nist.PublicKey)
		*p= *t.Public()
	case *gost3410.PrivateKey:
		p:=any(&pub).(*gost3410.PublicKey)
		*p= *t.Public()
	}
	address := crypto.PubkeyToAddress[P](pub)
	funds := big.NewInt(1000000000)
	dad := common.HexToAddress("0x0000000000000000000000000000000000000dad")

	ethConf := &ethconfig.Config[P]{
		Genesis: &core.Genesis[P]{
			Config:     params.AllEthashProtocolChanges,
			GasLimit:   11500000,
			Difficulty: big.NewInt(1048576),
			Alloc: core.GenesisAlloc{
				address: {Balance: funds},
				// The address 0xdad sloads 0x00 and 0x01
				dad: {
					Code: []byte{
						byte(vm.PC),
						byte(vm.PC),
						byte(vm.SLOAD),
						byte(vm.SLOAD),
					},
					Nonce:   0,
					Balance: big.NewInt(0),
				},
			},
		},
		Ethash: ethash.Config{
			PowMode: ethash.ModeFake,
		},
		NetworkId:               1337,
		TrieCleanCache:          5,
		TrieCleanCacheJournal:   "triecache",
		TrieCleanCacheRejournal: 60 * time.Minute,
		TrieDirtyCache:          5,
		TrieTimeout:             60 * time.Minute,
		SnapshotCache:           5,
	}

	ethBackend, err := eth.New[T,P](stack, ethConf)
	if err != nil {
		t.Fatalf("could not create eth backend: %v", err)
	}
	signer := types.LatestSigner[P](ethConf.Genesis.Config)

	legacyTx, _ := types.SignNewTx(key, signer, &types.LegacyTx{
		Nonce:    uint64(0),
		To:       &dad,
		Value:    big.NewInt(100),
		Gas:      50000,
		GasPrice: big.NewInt(1),
	})
	envelopTx, _ := types.SignNewTx(key, signer, &types.AccessListTx{
		ChainID:  ethConf.Genesis.Config.ChainID,
		Nonce:    uint64(1),
		To:       &dad,
		Gas:      30000,
		GasPrice: big.NewInt(1),
		Value:    big.NewInt(50),
		AccessList: types.AccessList{{
			Address:     dad,
			StorageKeys: []common.Hash{{0}},
		}},
	})

	// Create some blocks and import them
	chain, _ := core.GenerateChain[P](params.AllEthashProtocolChanges, ethBackend.BlockChain().Genesis(),
		 ethash.NewFaker[P](), ethBackend.ChainDb(), 1, func(i int, b *core.BlockGen[P]) {
			b.SetCoinbase(common.Address{1})
			b.AddTx(legacyTx)
			b.AddTx(envelopTx)
		})

	_, err = ethBackend.BlockChain().InsertChain(chain)
	if err != nil {
		t.Fatalf("could not create import blocks: %v", err)
	}
	// create gql service
	err = New[T,P](stack, ethBackend.APIBackend, []string{}, []string{})
	if err != nil {
		t.Fatalf("could not create graphql service: %v", err)
	}
}

// Quorum

// Tests that 400 is returned when an invalid RPC request is made.
func TestGraphQL_BadRequest(t *testing.T) {
	stack := createNode[nist.PrivateKey,nist.PublicKey](t, false, true)
	defer stack.Close()
	// start node
	if err := stack.Start(); err != nil {
		t.Fatalf("could not start node: %v", err)
	}
	// create http request
	body := strings.NewReader("{\"query\": \"{bleh{number}}\",\"variables\": null}")
	gqlReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/graphql", stack.HTTPEndpoint()), body)
	if err != nil {
		t.Fatalf("could not post: %v", err)
	}
	// read from response
	resp := doHTTPRequest(t, gqlReq)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read from response body: %v", err)
	}
	assert.Equal(t, "", string(bodyBytes)) // TODO: geth1.10.2: check changes
	assert.Equal(t, 404, resp.StatusCode)
}

func doHTTPRequest(t *testing.T, req *http.Request) *http.Response {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("could not issue a GET request to the given endpoint", err)
	}
	return resp
}

func TestQuorumSchema_PublicTransaction(t *testing.T) {
	saved := private.Ptm
	defer func() {
		private.Ptm = saved
	}()
	private.Ptm = &stubPrivateTransactionManager{}

	publicTx := types.NewTransaction[nist.PublicKey](0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), []byte("some random public payload"))
	publicTxQuery := &Transaction[nist.PrivateKey,nist.PublicKey]{tx: publicTx, backend: &StubBackend[nist.PrivateKey, nist.PublicKey]{}}
	isPrivate, err := publicTxQuery.IsPrivate(context.Background())
	if err != nil {
		t.Fatalf("Expect no error: %v", err)
	}
	if *isPrivate {
		t.Fatalf("Expect isPrivate to be false for public TX")
	}
	privateInputData, err := publicTxQuery.PrivateInputData(context.Background())
	if err != nil {
		t.Fatalf("Expect no error: %v", err)
	}
	if privateInputData.String() != "0x" {
		t.Fatalf("Expect privateInputData to be: \"0x\" for public TX, actual: %v", privateInputData.String())
	}
	internalPrivateTxQuery, err := publicTxQuery.PrivateTransaction(context.Background())
	if err != nil {
		t.Fatalf("Expect no error: %v", err)
	}
	if internalPrivateTxQuery != nil {
		t.Fatalf("Expect PrivateTransaction to be nil for non privacy precompile public tx, actual: %v", *internalPrivateTxQuery)
	}
}

func TestQuorumSchema_PrivateTransaction(t *testing.T) {
	saved := private.Ptm
	defer func() {
		private.Ptm = saved
	}()

	payloadHashByt := sha3.Sum512([]byte("arbitrary key"))
	arbitraryPayloadHash := common.BytesToEncryptedPayloadHash(payloadHashByt[:])
	private.Ptm = &stubPrivateTransactionManager{
		responses: map[common.EncryptedPayloadHash]ptmResponse{
			arbitraryPayloadHash: {
				body: []byte("private payload"), // equals to 0x70726976617465207061796c6f6164 after converting to bytes
				err:  nil,
			},
		},
	}

	privateTx := types.NewTransaction[nist.PublicKey](0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), arbitraryPayloadHash.Bytes())
	privateTx.SetPrivate()
	privateTxQuery := &Transaction[nist.PrivateKey,nist.PublicKey]{tx: privateTx, backend: &StubBackend[nist.PrivateKey, nist.PublicKey]{}}
	isPrivate, err := privateTxQuery.IsPrivate(context.Background())
	if err != nil {
		t.Fatalf("Expect no error: %v", err)
	}
	if !*isPrivate {
		t.Fatalf("Expect isPrivate to be true for private TX")
	}
	privateInputData, err := privateTxQuery.PrivateInputData(context.Background())
	if err != nil {
		t.Fatalf("Expect no error: %v", err)
	}
	if privateInputData.String() != "0x70726976617465207061796c6f6164" {
		t.Fatalf("Expect privateInputData to be: \"0x70726976617465207061796c6f6164\" for private TX, actual: %v", privateInputData.String())
	}
	internalPrivateTxQuery, err := privateTxQuery.PrivateTransaction(context.Background())
	if err != nil {
		t.Fatalf("Expect no error: %v", err)
	}
	if internalPrivateTxQuery != nil {
		t.Fatalf("Expect PrivateTransaction to be nil for non privacy precompile private tx, actual: %v", *internalPrivateTxQuery)
	}
}

func TestQuorumSchema_PrivacyMarkerTransaction(t *testing.T) {
	saved := private.Ptm
	defer func() {
		private.Ptm = saved
	}()

	encryptedPayloadHashByt := sha3.Sum512([]byte("encrypted payload hash"))
	encryptedPayloadHash := common.BytesToEncryptedPayloadHash(encryptedPayloadHashByt[:])

	privateTx := types.NewTransaction[nist.PublicKey](1, common.Address{}, big.NewInt(0), 0, big.NewInt(0), encryptedPayloadHash.Bytes())
	privateTx.SetPrivate()
	// json decoding later in the test requires the private tx to have signature values, so set to some arbitrary values here
	_, r, s := privateTx.RawSignatureValues()
	r.SetUint64(10)
	s.SetUint64(10)

	privateTxByt, _ := json.Marshal(privateTx)

	encryptedPrivateTxHashByt := sha3.Sum512([]byte("encrypted pvt tx hash"))
	encryptedPrivateTxHash := common.BytesToEncryptedPayloadHash(encryptedPrivateTxHashByt[:])

	private.Ptm = &stubPrivateTransactionManager{
		responses: map[common.EncryptedPayloadHash]ptmResponse{
			encryptedPayloadHash: {
				body: []byte("private payload"), // equals to 0x70726976617465207061796c6f6164 after converting to bytes
				err:  nil,
			},
			encryptedPrivateTxHash: {
				body: privateTxByt,
				err:  nil},
		},
	}

	privacyMarkerTx := types.NewTransaction[nist.PublicKey](0, common.QuorumPrivacyPrecompileContractAddress(), big.NewInt(0), 0, big.NewInt(0), encryptedPrivateTxHash.Bytes())

	pmtQuery := &Transaction[nist.PrivateKey, nist.PublicKey]{tx: privacyMarkerTx, backend: &StubBackend[nist.PrivateKey, nist.PublicKey]{}}
	isPrivate, err := pmtQuery.IsPrivate(context.Background())
	if err != nil {
		t.Fatalf("Expect no error: %v", err)
	}
	if *isPrivate {
		t.Fatalf("Expect isPrivate to be false for public PMT")
	}
	privateInputData, err := pmtQuery.PrivateInputData(context.Background())
	if err != nil {
		t.Fatalf("Expect no error: %v", err)
	}
	if privateInputData.String() != "0x" {
		t.Fatalf("Expect privateInputData to be: \"0x\" for public PMT, actual: %v", privateInputData.String())
	}

	internalPrivateTxQuery, err := pmtQuery.PrivateTransaction(context.Background())
	if err != nil {
		t.Fatalf("Expect no error: %v", err)
	}
	if internalPrivateTxQuery == nil {
		t.Fatal("Expect PrivateTransaction to be non-nil for privacy precompile PMT, actual is nil")
	}
	isPrivate, err = internalPrivateTxQuery.IsPrivate(context.Background())
	if err != nil {
		t.Fatalf("Expect no error: %v", err)
	}
	if !*isPrivate {
		t.Fatalf("Expect isPrivate to be true for internal private TX")
	}
	privateInputData, err = internalPrivateTxQuery.PrivateInputData(context.Background())
	if err != nil {
		t.Fatalf("Expect no error: %v", err)
	}
	if privateInputData.String() != "0x70726976617465207061796c6f6164" {
		t.Fatalf("Expect privateInputData to be: \"0x70726976617465207061796c6f6164\" for internal private TX, actual: %v", privateInputData.String())
	}
	nestedInternalPrivateTxQuery, err := internalPrivateTxQuery.PrivateTransaction(context.Background())
	if err != nil {
		t.Fatalf("Expect no error: %v", err)
	}
	if nestedInternalPrivateTxQuery != nil {
		t.Fatalf("Expect PrivateTransaction to be nil for internal private tx, actual: %v", *nestedInternalPrivateTxQuery)
	}
	_, ok := internalPrivateTxQuery.receiptGetter.(*privateTransactionReceiptGetter[nist.PrivateKey,nist.PublicKey])
	if !ok {
		t.Fatalf("Expect internal private txs receiptGetter to be of type *graphql.privateTransactionReceiptGetter, actual: %T", internalPrivateTxQuery.receiptGetter)
	}
}

func TestQuorumTransaction_getReceipt_defaultReceiptGetter(t *testing.T) {
	graphqlTx := &Transaction[nist.PrivateKey,nist.PublicKey]{tx: &types.Transaction[nist.PublicKey]{}, backend: &StubBackend[nist.PrivateKey,nist.PublicKey]{}}

	if graphqlTx.receiptGetter != nil {
		t.Fatalf("Expect nil receiptGetter: actual %v", graphqlTx.receiptGetter)
	}

	_, _ = graphqlTx.getReceipt(context.Background())

	if graphqlTx.receiptGetter == nil {
		t.Fatalf("Expect default receiptGetter to have been set: actual nil")
	}

	if _, ok := graphqlTx.receiptGetter.(*transactionReceiptGetter[nist.PrivateKey,nist.PublicKey]); !ok {
		t.Fatalf("Expect default receiptGetter to be of type *graphql.transactionReceiptGetter: actual %T", graphqlTx.receiptGetter)
	}
}

type ptmResponse struct {
	body []byte
	err  error
}

type stubPrivateTransactionManager struct {
	notinuse.PrivateTransactionManager
	responses map[common.EncryptedPayloadHash]ptmResponse
}

func (spm *stubPrivateTransactionManager) HasFeature(f engine.PrivateTransactionManagerFeature) bool {
	return true
}

func (spm *stubPrivateTransactionManager) Receive(txHash common.EncryptedPayloadHash) (string, []string, []byte, *engine.ExtraMetadata, error) {
	res, ok := spm.responses[txHash]
	if !ok {
		return "", nil, nil, nil, nil
	}
	if res.err != nil {
		return "", nil, nil, nil, res.err
	}
	meta := &engine.ExtraMetadata{PrivacyFlag: engine.PrivacyFlagStandardPrivate}
	return "", nil, res.body, meta, nil
}

func (spm *stubPrivateTransactionManager) ReceiveRaw(hash common.EncryptedPayloadHash) ([]byte, string, *engine.ExtraMetadata, error) {
	_, sender, data, metadata, err := spm.Receive(hash)
	return data, sender[0], metadata, err
}

type StubBackend [T crypto.PrivateKey, P crypto.PublicKey]struct{}

var _ ethapi.Backend[nist.PrivateKey,nist.PublicKey] = &StubBackend[nist.PrivateKey,nist.PublicKey]{}

func (sb *StubBackend[T,P]) CurrentHeader() *types.Header[P] {
	panic("implement me")
}

func (sb *StubBackend[T,P]) Engine() consensus.Engine[P] {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SupportsMultitenancy(rpcCtx context.Context) (*proto.PreAuthenticatedAuthenticationToken, bool) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) AccountExtraDataStateGetterByNumber(context.Context, rpc.BlockNumber) (vm.AccountExtraDataStateGetter, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) IsAuthorized(authToken *proto.PreAuthenticatedAuthenticationToken, attributes ...*multitenancy.PrivateStateSecurityAttribute) (bool, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetEVM(ctx context.Context, msg core.Message, state vm.MinimalApiState, header *types.Header[P], vmconfig *vm.Config[P]) (*vm.EVM[P], func() error, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) CurrentBlock() *types.Block[nist.PublicKey] {
	panic("implement me")
}

func (sb *StubBackend[T,P]) Downloader() *downloader.Downloader[T,P] {
	panic("implement me")
}

func (sb *StubBackend[T,P]) ProtocolVersion() int {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SuggestPrice(ctx context.Context) (*big.Int, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) ChainDb() ethdb.Database {
	panic("implement me")
}

func (sb *StubBackend[T,P]) EventMux() *event.TypeMux {
	panic("implement me")
}

func (sb *StubBackend[T,P]) AccountManager() *accounts.Manager[P] {
	panic("implement me")
}

func (sb *StubBackend[T,P]) ExtRPCEnabled() bool {
	panic("implement me")
}

func (sb *StubBackend[T,P]) CallTimeOut() time.Duration {
	panic("implement me")
}

func (sb *StubBackend[T,P]) RPCTxFeeCap() float64 {
	panic("implement me")
}

func (sb *StubBackend[T,P]) RPCGasCap() uint64 {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SetHead(number uint64) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header[P], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header[P], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header[P], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block[nist.PublicKey], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block[nist.PublicKey], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block[nist.PublicKey], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (vm.MinimalApiState, *types.Header[P], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (vm.MinimalApiState, *types.Header[P], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts[P], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetTd(ctx context.Context, hash common.Hash) *big.Int {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribeChainEvent(ch chan<- core.ChainEvent[P]) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent[P]) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent[P]) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SendTx(ctx context.Context, signedTx *types.Transaction[nist.PublicKey]) error {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction[P], common.Hash, uint64, uint64, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetPoolTransactions() (types.Transactions[P], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetPoolTransaction(txHash common.Hash) *types.Transaction[P] {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) Stats() (pending int, queued int) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) TxPoolContent() (map[common.Address]types.Transactions[P], map[common.Address]types.Transactions[P]) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribeNewTxsEvent(chan<- core.NewTxsEvent[P]) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) BloomStatus() (uint64, uint64) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetLogs(ctx context.Context, blockHash common.Hash) ([][]*types.Log, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent[P]) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) ChainConfig() *params.ChainConfig {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) PSMR() mps.PrivateStateMetadataResolver {
	return &StubPSMR{}
}

func (sb *StubBackend[T,P]) IsPrivacyMarkerTransactionCreationEnabled() bool {
	panic("implement me")
}

func (sb *StubBackend[T,P]) UnprotectedAllowed() bool {
	panic("implement me")
}

type StubPSMR struct {
}

func (psmr *StubPSMR) ResolveForManagedParty(managedParty string) (*mps.PrivateStateMetadata, error) {
	panic("implement me")
}
func (psmr *StubPSMR) ResolveForUserContext(ctx context.Context) (*mps.PrivateStateMetadata, error) {
	return mps.DefaultPrivateStateMetadata, nil
}
func (psmr *StubPSMR) PSIs() []types.PrivateStateIdentifier {
	panic("implement me")
}
func (psmr *StubPSMR) NotIncludeAny(psm *mps.PrivateStateMetadata, managedParties ...string) bool {
	return false
}
