package raft

import (
	"errors"

	"github.com/coreos/etcd/pkg/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

type RaftNodeInfo struct {
	ClusterSize    int        `json:"clusterSize"`
	Role           string     `json:"role"`
	Address        *Address   `json:"address"`
	PeerAddresses  []*Address `json:"peerAddresses"`
	RemovedPeerIds []uint16   `json:"removedPeerIds"`
	AppliedIndex   uint64     `json:"appliedIndex"`
	SnapshotIndex  uint64     `json:"snapshotIndex"`
}

type PublicRaftAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	raftService *RaftService[T,P]
}

func NewPublicRaftAPI[T crypto.PrivateKey, P crypto.PublicKey](raftService *RaftService[T,P]) *PublicRaftAPI[T,P] {
	return &PublicRaftAPI[T,P]{raftService}
}

func (s *PublicRaftAPI[T,P]) Role() string {
	if err := s.checkIfNodeInCluster(); err != nil {
		return ""
	}
	_, err := s.raftService.raftProtocolManager.LeaderAddress()
	if err != nil {
		return ""
	}
	return s.raftService.raftProtocolManager.NodeInfo().Role
}

// helper function to check if self node is part of cluster
func (s *PublicRaftAPI[T,P]) checkIfNodeInCluster() error {
	if s.raftService.raftProtocolManager.IsIDRemoved(uint64(s.raftService.raftProtocolManager.raftId)) {
		return errors.New("node not part of raft cluster. operations not allowed")
	}
	return nil
}

func (s *PublicRaftAPI[T,P]) AddPeer(enodeId string) (uint16, error) {
	if err := s.checkIfNodeInCluster(); err != nil {
		return 0, err
	}
	return s.raftService.raftProtocolManager.ProposeNewPeer(enodeId, false)
}

func (s *PublicRaftAPI[T,P]) AddLearner(enodeId string) (uint16, error) {
	if err := s.checkIfNodeInCluster(); err != nil {
		return 0, err
	}
	return s.raftService.raftProtocolManager.ProposeNewPeer(enodeId, true)
}

func (s *PublicRaftAPI[T,P]) PromoteToPeer(raftId uint16) (bool, error) {
	if err := s.checkIfNodeInCluster(); err != nil {
		return false, err
	}
	return s.raftService.raftProtocolManager.PromoteToPeer(raftId)
}

func (s *PublicRaftAPI[T,P]) RemovePeer(raftId uint16) error {
	if err := s.checkIfNodeInCluster(); err != nil {
		return err
	}
	return s.raftService.raftProtocolManager.ProposePeerRemoval(raftId)
}

func (s *PublicRaftAPI[T,P]) Leader() (string, error) {

	addr, err := s.raftService.raftProtocolManager.LeaderAddress()
	if err != nil {
		return "", err
	}
	return addr.NodeId.String(), nil
}

func (s *PublicRaftAPI[T,P]) Cluster() ([]ClusterInfo, error) {
	// check if the node has already been removed from cluster
	// if yes return nil
	if err := s.checkIfNodeInCluster(); err != nil {
		return []ClusterInfo{}, nil
	}

	nodeInfo := s.raftService.raftProtocolManager.NodeInfo()
	if nodeInfo.Role == "" {
		return []ClusterInfo{}, nil
	}

	noLeader := false
	leaderAddr, err := s.raftService.raftProtocolManager.LeaderAddress()
	if err != nil {
		noLeader = true
		if s.raftService.raftProtocolManager.NodeInfo().Role == "" {
			return []ClusterInfo{}, nil
		}
	}

	peerAddresses := append(nodeInfo.PeerAddresses, nodeInfo.Address)
	clustInfo := make([]ClusterInfo, len(peerAddresses))
	for i, a := range peerAddresses {
		role := ""
		if !noLeader {
			if a.RaftId == leaderAddr.RaftId {
				role = "minter"
			} else if s.raftService.raftProtocolManager.isLearner(a.RaftId) {
				role = "learner"
			} else if s.raftService.raftProtocolManager.isVerifier(a.RaftId) {
				role = "verifier"
			}
		}
		clustInfo[i] = ClusterInfo{*a, role, s.checkIfNodeIsActive(a.RaftId)}
	}
	return clustInfo, nil
}

// checkIfNodeIsActive checks if the raft node is active
// if the raft node is active ActiveSince returns non-zero time
func (s *PublicRaftAPI[T,P]) checkIfNodeIsActive(raftId uint16) bool {
	if raftId == s.raftService.raftProtocolManager.raftId {
		return true
	}
	activeSince := s.raftService.raftProtocolManager.transport.ActiveSince(types.ID(raftId))
	return !activeSince.IsZero()
}

func (s *PublicRaftAPI[T,P]) GetRaftId(enodeId string) (uint16, error) {
	return s.raftService.raftProtocolManager.FetchRaftId(enodeId)
}
