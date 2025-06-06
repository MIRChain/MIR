package raft

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/coreos/etcd/raft/raftpb"
	mapset "github.com/deckarep/golang-set"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/hexutil"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/node"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/rlp"
)

const TEST_URL = "enode://3d9ca5956b38557aba991e31cf510d4df641dce9cc26bfeb7de082f0c07abb6ede3a58410c8f249dabeecee4ad3979929ac4c7c496ad20b8cfdd061b7401b4f5@127.0.0.1:21003?discport=0&raftport=50404"

func TestSignHeader(t *testing.T) {
	//create only what we need to test the seal
	var testRaftId uint16 = 5
	config := &node.Config[nist.PrivateKey,nist.PublicKey]{Name: "unit-test", DataDir: ""}
	nodeKey := config.NodeKey()

	raftProtocolManager := &ProtocolManager[nist.PrivateKey,nist.PublicKey]{raftId: testRaftId}
	raftService := &RaftService[nist.PrivateKey,nist.PublicKey]{nodeKey: nodeKey, raftProtocolManager: raftProtocolManager}
	minter := minter[nist.PrivateKey,nist.PublicKey]{eth: raftService}

	//create some fake header to sign
	fakeParentHash := common.HexToHash("0xc2c1dc1be8054808c69e06137429899d")

	header := &types.Header[nist.PublicKey]{
		ParentHash: fakeParentHash,
		Number:     big.NewInt(1),
		Difficulty: big.NewInt(1),
		GasLimit:   uint64(0),
		GasUsed:    uint64(0),
		Coinbase:   minter.coinbase,
		Time:       uint64(time.Now().UnixNano()),
	}

	headerHash := header.Hash()
	extraDataBytes := minter.buildExtraSeal(headerHash)
	var seal *extraSeal
	err := rlp.DecodeBytes(extraDataBytes[:], &seal)
	if err != nil {
		t.Fatalf("Unable to decode seal: %s", err.Error())
	}

	// Check raftId
	sealRaftId, err := hexutil.DecodeUint64("0x" + string(seal.RaftId)) //add the 0x prefix
	if err != nil {
		t.Errorf("Unable to get RaftId: %s", err.Error())
	}
	if sealRaftId != uint64(testRaftId) {
		t.Errorf("RaftID does not match. Expected: %d, Actual: %d", testRaftId, sealRaftId)
	}

	//Identify who signed it
	sig := seal.Signature
	pubKey, err := crypto.SigToPub[nist.PublicKey](headerHash.Bytes(), sig)
	if err != nil {
		t.Fatalf("Unable to get public key from signature: %s", err.Error())
	}

	//Compare derived public key to original public key
	if pubKey.X.Cmp(nodeKey.X) != 0 {
		t.Errorf("Signature incorrect!")
	}

}

// func TestSignHeaderCsp(t *testing.T) {
// 	//create only what we need to test the seal
// 	var testRaftId uint16 = 5
// 	crypto.CryptoAlg = crypto.GOST_CSP
// 	store, err := csp.SystemStore("My")
// 	if err != nil {
// 		t.Errorf("Store error: %s", err)
// 	}
// 	defer store.Close()
// 	crt, err := store.GetBySubjectId("43ad4195a67f95eea752861c96297045bb9ea5a7")
// 	if err != nil {
// 		t.Errorf("Get cert error: %s", err)
// 	}
// 	defer crt.Close()
// 	config := &node.Config[nist.PrivateKey,nist.PublicKey]{Name: "unit-test", DataDir: "", SignerCert: &crt}
// 	sigCert, err := config.GetSignerCert()
// 	if err != nil {
// 		t.Fatalf("Unable to get signer cert: %s", err.Error())
// 	}
// 	raftProtocolManager := &ProtocolManager{raftId: testRaftId}
// 	raftService := &RaftService{signerCert: sigCert, raftProtocolManager: raftProtocolManager}
// 	minter := minter{eth: raftService}

// 	//create some fake header to sign
// 	fakeParentHash := common.HexToHash("0xc2c1dc1be8054808c69e06137429899d")

// 	header := &types.Header{
// 		ParentHash: fakeParentHash,
// 		Number:     big.NewInt(1),
// 		Difficulty: big.NewInt(1),
// 		GasLimit:   uint64(0),
// 		GasUsed:    uint64(0),
// 		Coinbase:   minter.coinbase,
// 		Time:       uint64(time.Now().UnixNano()),
// 	}

// 	headerHash := header.Hash()
// 	extraDataBytes := minter.buildExtraSeal(headerHash)
// 	var seal *extraSeal
// 	err = rlp.DecodeBytes(extraDataBytes[:], &seal)
// 	if err != nil {
// 		t.Fatalf("Unable to decode seal: %s", err.Error())
// 	}

// 	// Check raftId
// 	sealRaftId, err := hexutil.DecodeUint64("0x" + string(seal.RaftId)) //add the 0x prefix
// 	if err != nil {
// 		t.Errorf("Unable to get RaftId: %s", err.Error())
// 	}
// 	if sealRaftId != uint64(testRaftId) {
// 		t.Errorf("RaftID does not match. Expected: %d, Actual: %d", testRaftId, sealRaftId)
// 	}

// 	//Identify who signed it
// 	sig := seal.Signature
// 	pubKey, err := crypto.SigToPubCsp(headerHash.Bytes(), sig)
// 	if err != nil {
// 		t.Fatalf("Unable to get public key from signature: %s", err.Error())
// 	}

// 	//Compare derived public key to original public key
// 	if bytes.Compare(pubKey, sigCert.Info().PublicKeyBytes()[2:66]) != 0 {
// 		t.Errorf("Signature incorrect!")
// 	}
// }

func TestAddLearner_whenTypical(t *testing.T) {

	raftService := newTestRaftService[nist.PrivateKey,nist.PublicKey](t, 1, []uint64{1}, []uint64{})

	propPeer := func() {
		raftid, err := raftService.raftProtocolManager.ProposeNewPeer(TEST_URL, true)
		if err != nil {
			t.Errorf("propose new peer failed %v\n", err)
		}
		if raftid != raftService.raftProtocolManager.raftId+1 {
			t.Errorf("1. wrong raft id. expected %d got %d\n", raftService.raftProtocolManager.raftId+1, raftid)
		}
	}
	go propPeer()
	select {
	case confChange := <-raftService.raftProtocolManager.confChangeProposalC:
		if confChange.Type != raftpb.ConfChangeAddLearnerNode {
			t.Errorf("expected ConfChangeAddLearnerNode but got %s", confChange.Type.String())
		}
		if uint16(confChange.NodeID) != raftService.raftProtocolManager.raftId+1 {
			t.Errorf("2. wrong raft id. expected %d got %d\n", raftService.raftProtocolManager.raftId+1, uint16(confChange.NodeID))
		}
	case <-time.After(time.Millisecond * 200):
		t.Errorf("add learner conf change not received")
	}
}

func TestPromoteLearnerToPeer_whenTypical(t *testing.T) {
	learnerRaftId := uint16(3)
	raftService := newTestRaftService[nist.PrivateKey,nist.PublicKey](t, 2, []uint64{2}, []uint64{uint64(learnerRaftId)})
	promoteToPeer := func() {
		ok, err := raftService.raftProtocolManager.PromoteToPeer(learnerRaftId)
		if err != nil || !ok {
			t.Errorf("promote learner to peer failed %v\n", err)
		}
	}
	go promoteToPeer()
	select {
	case confChange := <-raftService.raftProtocolManager.confChangeProposalC:
		if confChange.Type != raftpb.ConfChangeAddNode {
			t.Errorf("expected ConfChangeAddNode but got %s", confChange.Type.String())
		}
		if uint16(confChange.NodeID) != learnerRaftId {
			t.Errorf("2. wrong raft id. expected %d got %d\n", learnerRaftId, uint16(confChange.NodeID))
		}
	case <-time.After(time.Millisecond * 200):
		t.Errorf("add learner conf change not received")
	}
}

func TestAddLearnerOrPeer_fromLearner(t *testing.T) {

	raftService := newTestRaftService[nist.PrivateKey,nist.PublicKey](t, 3, []uint64{2}, []uint64{3})

	_, err := raftService.raftProtocolManager.ProposeNewPeer(TEST_URL, true)

	if err == nil {
		t.Errorf("learner should not be allowed to add learner or peer")
	}

	if err != nil && !strings.Contains(err.Error(), "learner node can't add peer or learner") {
		t.Errorf("expect error message: propose new peer failed, got: %v\n", err)
	}

	_, err = raftService.raftProtocolManager.ProposeNewPeer(TEST_URL, false)

	if err == nil {
		t.Errorf("learner should not be allowed to add learner or peer")
	}

	if err != nil && !strings.Contains(err.Error(), "learner node can't add peer or learner") {
		t.Errorf("expect error message: propose new peer failed, got: %v\n", err)
	}

}

func TestPromoteLearnerToPeer_fromLearner(t *testing.T) {
	learnerRaftId := uint16(3)
	raftService := newTestRaftService[nist.PrivateKey,nist.PublicKey](t, 2, []uint64{1}, []uint64{2, uint64(learnerRaftId)})

	_, err := raftService.raftProtocolManager.PromoteToPeer(learnerRaftId)

	if err == nil {
		t.Errorf("learner should not be allowed to promote to peer")
	}

	if err != nil && !strings.Contains(err.Error(), "learner node can't promote to peer") {
		t.Errorf("expect error message: propose new peer failed, got: %v\n", err)
	}

}

func enodeId(id string, ip string, raftPort int) string {
	return fmt.Sprintf("enode://%s@%s?discport=0&raftport=%d", id, ip, raftPort)
}

func peerList[P crypto.PublicKey](url string) (error, []*enode.Node[P]) {
	var nodes []*enode.Node[P]
	node, err := enode.ParseV4[P](url)
	if err != nil {
		return fmt.Errorf("Node URL %s: %v\n", url, err), nil
	}
	nodes = append(nodes, node)
	return nil, nodes
}

func newTestRaftService[T crypto.PrivateKey, P crypto.PublicKey](t *testing.T, raftId uint16, nodes []uint64, learners []uint64) *RaftService[T,P] {
	//create only what we need to test add learner node
	config := &node.Config[T,P]{Name: "unit-test", DataDir: ""}
	// This will create a new node key, which is needed to set a stub p2p.Server and avoid `nil pointer dereference` when testing.
	nodeKey := config.NodeKey()
	mockp2pConfig := p2p.Config[T,P]{Name: "unit-test", ListenAddr: "30303", PrivateKey: nodeKey}
	mockp2p := &p2p.Server[T,P]{Config: mockp2pConfig}
	var pub P
	switch t:=any(&nodeKey).(type){
	case *nist.PrivateKey:
		p:= any(&pub).(*nist.PublicKey)
		*p=*t.Public()
	case *gost3410.PrivateKey:
		p:= any(&pub).(*gost3410.PublicKey)
		*p=*t.Public()
	}
	enodeIdStr := fmt.Sprintf("%x", crypto.FromECDSAPub[P](pub)[1:])
	url := enodeId(enodeIdStr, "127.0.0.1:21001", 50401)
	err, peers := peerList[P](url)
	if err != nil {
		t.Errorf("getting peers failed %v", err)
	}
	raftProtocolManager := &ProtocolManager[T,P]{
		raftId:              raftId,
		bootstrapNodes:      peers,
		confChangeProposalC: make(chan raftpb.ConfChange),
		removedPeers:        mapset.NewSet(),
		confState:           raftpb.ConfState{Nodes: nodes, Learners: learners},
		p2pServer:           mockp2p,
	}
	raftService := &RaftService[T,P]{nodeKey: nodeKey, raftProtocolManager: raftProtocolManager}
	return raftService
}
