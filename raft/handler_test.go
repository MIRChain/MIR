package raft

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/coreos/etcd/wal"
	"github.com/coreos/etcd/wal/walpb"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/eth"
	"github.com/pavelkrolevets/MIR-pro/eth/ethconfig"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/miner"
	"github.com/pavelkrolevets/MIR-pro/node"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/params"
)

// pm.advanceAppliedIndex() and state updates are in different
// transaction boundaries hence there's a probablity that they are
// out of sync due to premature shutdown
func TestProtocolManager_whenAppliedIndexOutOfSync(t *testing.T) {
	logger := log.New()
	logger.SetHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(false)))
	tmpWorkingDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpWorkingDir)
	}()
	count := 3
	ports := make([]uint16, count)
	nodeKeys := make([]nist.PrivateKey, count)
	peers := make([]*enode.Node[nist.PublicKey], count)
	for i := 0; i < count; i++ {
		ports[i] = nextPort(t)
		nodeKeys[i] = mustNewNodeKey(t)
		peers[i] = enode.NewV4Hostname[nist.PublicKey](*nodeKeys[i].Public(), net.IPv4(127, 0, 0, 1).String(), 0, 0, int(ports[i]))
	}
	raftNodes := make([]*RaftService[nist.PrivateKey,nist.PublicKey], count)
	for i := 0; i < count; i++ {
		if s, err := startRaftNode(uint16(i+1), ports[i], tmpWorkingDir, nodeKeys[i], peers); err != nil {
			t.Fatal(err)
		} else {
			raftNodes[i] = s
		}
	}
	waitFunc := func() {
		for {
			time.Sleep(10 * time.Millisecond)
			for i := 0; i < count; i++ {
				if raftNodes[i].raftProtocolManager.role == minterRole {
					return
				}
			}
		}
	}
	waitFunc()
	logger.Debug("stop the cluster")
	for i := 0; i < count; i++ {
		if err := raftNodes[i].Stop(); err != nil {
			t.Fatal(err)
		}
		// somehow the wal dir is still being locked that causes failures in subsequent start
		// we need to check here to make sure everything is fully stopped
		for isWalDirStillLocked(fmt.Sprintf("%s/node%d/raft-wal", tmpWorkingDir, i+1)) {
			logger.Debug("sleep...", "i", i)
			time.Sleep(10 * time.Millisecond)
		}
		logger.Debug("node stopped", "id", i)
	}
	logger.Debug("update applied index")
	// update the index to mimic the issue (set applied index behind for node 0)
	if err := writeAppliedIndex(tmpWorkingDir, 0, 1); err != nil {
		t.Fatal(err)
	}
	//time.Sleep(3 * time.Second)
	logger.Debug("restart the cluster")
	for i := 0; i < count; i++ {
		if s, err := startRaftNode(uint16(i+1), ports[i], tmpWorkingDir, nodeKeys[i], peers); err != nil {
			t.Fatal(err)
		} else {
			raftNodes[i] = s
		}
	}
	waitFunc()
}

func isWalDirStillLocked(walDir string) bool {
	var snap walpb.Snapshot
	w, err := wal.Open(walDir, snap)
	if err != nil {
		return true
	}
	defer func() {
		_ = w.Close()
	}()
	return false
}

func writeAppliedIndex(workingDir string, node int, index uint64) error {
	db, err := openQuorumRaftDb(fmt.Sprintf("%s/node%d/quorum-raft-state", workingDir, node+1))
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, index)
	return db.Put(appliedDbKey, buf, noFsync)
}

func mustNewNodeKey(t *testing.T) nist.PrivateKey {
	k, err := crypto.GenerateKey[nist.PrivateKey]()
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func nextPort(t *testing.T) uint16 {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	return uint16(listener.Addr().(*net.TCPAddr).Port)
}

func prepareServiceContext(key nist.PrivateKey) (stack *node.Node[nist.PrivateKey,nist.PublicKey], cfg *node.Config[nist.PrivateKey,nist.PublicKey], err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
			stack = nil
			cfg = nil
		}
	}()
	cfg = &node.Config[nist.PrivateKey,nist.PublicKey]{
		P2P: p2p.Config[nist.PrivateKey,nist.PublicKey]{
			PrivateKey: key,
		},
	}
	stack, _ = node.New(cfg)
	return
}

func startRaftNode(id, port uint16, tmpWorkingDir string, key nist.PrivateKey, nodes []*enode.Node[nist.PublicKey]) (*RaftService[nist.PrivateKey,nist.PublicKey], error) {
	raftlogdir := fmt.Sprintf("%s/node%d", tmpWorkingDir, id)

	stack, _, err := prepareServiceContext(key)
	if err != nil {
		return nil, err
	}

	const testAddress = "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	e, err := eth.New(stack, &ethconfig.Config[nist.PublicKey]{
		Genesis: &core.Genesis[nist.PublicKey]{Config: params.QuorumTestChainConfig},
		Miner:   miner.Config{Etherbase: common.HexToAddress(testAddress)},
	})
	if err != nil {
		return nil, err
	}

	s, err := New(stack, params.QuorumTestChainConfig, id, port, false, 100*time.Millisecond, e, nodes, raftlogdir, false)
	if err != nil {
		return nil, err
	}

	if err := stack.Server().Start(); err != nil {
		return nil, fmt.Errorf("could not start: %v", err)
	}
	if err := s.Start(); err != nil {
		return nil, err
	}

	return s, nil
}
