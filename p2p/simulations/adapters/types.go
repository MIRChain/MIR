// Copyright 2017 The go-ethereum Authors
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

package adapters

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"

	"github.com/docker/docker/pkg/reexec"
	"github.com/gorilla/websocket"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/node"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/p2p/enr"
	"github.com/pavelkrolevets/MIR-pro/rpc"
)

// Node represents a node in a simulation network which is created by a
// NodeAdapter, for example:
//
// * SimNode    - An in-memory node
// * ExecNode   - A child process node
// * DockerNode - A Docker container node
//
type Node interface {
	// Addr returns the node's address (e.g. an Enode URL)
	Addr() []byte

	// Client returns the RPC client which is created once the node is
	// up and running
	Client() (*rpc.Client, error)

	// ServeRPC serves RPC requests over the given connection
	ServeRPC(*websocket.Conn) error

	// Start starts the node with the given snapshots
	Start(snapshots map[string][]byte) error

	// Stop stops the node
	Stop() error

	// NodeInfo returns information about the node
	NodeInfo() *p2p.NodeInfo

	// Snapshots creates snapshots of the running services
	Snapshots() (map[string][]byte, error)
}

// NodeAdapter is used to create Nodes in a simulation network
type NodeAdapter  [T crypto.PrivateKey, P crypto.PublicKey]  interface {
	// Name returns the name of the adapter for logging purposes
	Name() string

	// NewNode creates a new node with the given configuration
	NewNode(config *NodeConfig[T,P]) (Node, error)
}

// NodeConfig is the configuration used to start a node in a simulation
// network
type NodeConfig  [T crypto.PrivateKey, P crypto.PublicKey]  struct {
	// ID is the node's ID which is used to identify the node in the
	// simulation network
	ID enode.ID

	// PrivateKey is the node's private key which is used by the devp2p
	// stack to encrypt communications
	PrivateKey T

	// Enable peer events for Msgs
	EnableMsgEvents bool

	// Name is a human friendly name for the node like "node01"
	Name string

	// Use an existing database instead of a temporary one if non-empty
	DataDir string

	// Lifecycles are the names of the service lifecycles which should be run when
	// starting the node (for SimNodes it should be the names of service lifecycles
	// contained in SimAdapter.lifecycles, for other nodes it should be
	// service lifecycles registered by calling the RegisterLifecycle function)
	Lifecycles []string

	// Properties are the names of the properties this node should hold
	// within running services (e.g. "bootnode", "lightnode" or any custom values)
	// These values need to be checked and acted upon by node Services
	Properties []string

	// ExternalSigner specifies an external URI for a clef-type signer
	ExternalSigner string

	// Enode
	node *enode.Node[P]

	// ENR Record with entries to overwrite
	Record enr.Record

	// function to sanction or prevent suggesting a peer
	Reachable func(id enode.ID) bool

	Port uint16

	// LogFile is the log file name of the p2p node at runtime.
	//
	// The default value is empty so that the default log writer
	// is the system standard output.
	LogFile string

	// LogVerbosity is the log verbosity of the p2p node at runtime.
	//
	// The default verbosity is INFO.
	LogVerbosity log.Lvl
}

// nodeConfigJSON is used to encode and decode NodeConfig as JSON by encoding
// all fields as strings
type nodeConfigJSON struct {
	ID              string   `json:"id"`
	PrivateKey      string   `json:"private_key"`
	Name            string   `json:"name"`
	Lifecycles      []string `json:"lifecycles"`
	Properties      []string `json:"properties"`
	EnableMsgEvents bool     `json:"enable_msg_events"`
	Port            uint16   `json:"port"`
	LogFile         string   `json:"logfile"`
	LogVerbosity    int      `json:"log_verbosity"`
}

// MarshalJSON implements the json.Marshaler interface by encoding the config
// fields as strings
func (n *NodeConfig[T,P]) MarshalJSON() ([]byte, error) {
	confJSON := nodeConfigJSON{
		ID:              n.ID.String(),
		Name:            n.Name,
		Lifecycles:      n.Lifecycles,
		Properties:      n.Properties,
		Port:            n.Port,
		EnableMsgEvents: n.EnableMsgEvents,
		LogFile:         n.LogFile,
		LogVerbosity:    int(n.LogVerbosity),
	}
	if !reflect.ValueOf(&n.PrivateKey).IsZero(){
		confJSON.PrivateKey = hex.EncodeToString(crypto.FromECDSA(n.PrivateKey))
	}
	return json.Marshal(confJSON)
}

// UnmarshalJSON implements the json.Unmarshaler interface by decoding the json
// string values into the config fields
func (n *NodeConfig[T,P]) UnmarshalJSON(data []byte) error {
	var confJSON nodeConfigJSON
	if err := json.Unmarshal(data, &confJSON); err != nil {
		return err
	}

	if confJSON.ID != "" {
		if err := n.ID.UnmarshalText([]byte(confJSON.ID)); err != nil {
			return err
		}
	}

	if confJSON.PrivateKey != "" {
		key, err := hex.DecodeString(confJSON.PrivateKey)
		if err != nil {
			return err
		}
		privKey, err := crypto.ToECDSA[T](key)
		if err != nil {
			return err
		}
		n.PrivateKey = privKey
	}

	n.Name = confJSON.Name
	n.Lifecycles = confJSON.Lifecycles
	n.Properties = confJSON.Properties
	n.Port = confJSON.Port
	n.EnableMsgEvents = confJSON.EnableMsgEvents
	n.LogFile = confJSON.LogFile
	n.LogVerbosity = log.Lvl(confJSON.LogVerbosity)

	return nil
}

// Node returns the node descriptor represented by the config.
func (n *NodeConfig[T,P]) Node() *enode.Node[P] {
	return n.node
}

// RandomNodeConfig returns node configuration with a randomly generated ID and
// PrivateKey
func RandomNodeConfig [T crypto.PrivateKey, P crypto.PublicKey] () *NodeConfig[T,P] {
	prvkey, err := crypto.GenerateKey[T]()
	if err != nil {
		panic("unable to generate key")
	}

	port, err := assignTCPPort()
	if err != nil {
		panic("unable to assign tcp port")
	}
	var pub P
	switch t:=any(&prvkey).(type) {
	case *nist.PrivateKey:
		p:=any(&pub).(*nist.PublicKey)
		*p = *t.Public()
	case *gost3410.PrivateKey:
		p:=any(&pub).(*gost3410.PublicKey)
		*p = *t.Public()
	}
	enodId := enode.PubkeyToIDV4(pub)
	return &NodeConfig[T,P]{
		PrivateKey:      prvkey,
		ID:              enodId,
		Name:            fmt.Sprintf("node_%s", enodId.String()),
		Port:            port,
		EnableMsgEvents: true,
		LogVerbosity:    log.LvlInfo,
	}
}

func assignTCPPort() (uint16, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	l.Close()
	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return 0, err
	}
	p, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint16(p), nil
}

// ServiceContext is a collection of options and methods which can be utilised
// when starting services
type ServiceContext  [T crypto.PrivateKey, P crypto.PublicKey]  struct {
	RPCDialer

	Config   *NodeConfig[T,P]
	Snapshot []byte
}

// RPCDialer is used when initialising services which need to connect to
// other nodes in the network (for example a simulated Swarm node which needs
// to connect to a Geth node to resolve ENS names)
type RPCDialer interface {
	DialRPC(id enode.ID) (*rpc.Client, error)
}

// LifecycleConstructor allows a Lifecycle to be constructed during node start-up.
// While the service-specific package usually takes care of Lifecycle creation and registration,
// for testing purposes, it is useful to be able to construct a Lifecycle on spot.
type LifecycleConstructor [T crypto.PrivateKey, P crypto.PublicKey]  func(ctx *ServiceContext[T,P], stack *node.Node[T,P]) (node.Lifecycle, error)

// LifecycleConstructors stores LifecycleConstructor functions to call during node start-up.
type LifecycleConstructors [T crypto.PrivateKey, P crypto.PublicKey]  map[string]LifecycleConstructor[T,P]

// lifecycleConstructorFuncs is a map of registered services which are used to boot devp2p
// nodes
var lifecycleConstructorFuncs = make(LifecycleConstructors)

// RegisterLifecycles registers the given Services which can then be used to
// start devp2p nodes using either the Exec or Docker adapters.
//
// It should be called in an init function so that it has the opportunity to
// execute the services before main() is called.
func RegisterLifecycles [T crypto.PrivateKey, P crypto.PublicKey] (lifecycles LifecycleConstructors[T,P]) {
	for name, f := range lifecycles {
		if _, exists := lifecycleConstructorFuncs[name]; exists {
			panic(fmt.Sprintf("node service already exists: %q", name))
		}
		lifecycleConstructorFuncs[name] = f
	}

	// now we have registered the services, run reexec.Init() which will
	// potentially start one of the services if the current binary has
	// been exec'd with argv[0] set to "p2p-node"
	if reexec.Init() {
		os.Exit(0)
	}
}

// adds the host part to the configuration's ENR, signs it
// creates and  the corresponding enode object to the configuration
func (n *NodeConfig[T,P]) initEnode(ip net.IP, tcpport int, udpport int) error {
	enrIp := enr.IP(ip)
	n.Record.Set(&enrIp)
	enrTcpPort := enr.TCP(tcpport)
	n.Record.Set(&enrTcpPort)
	enrUdpPort := enr.UDP(udpport)
	n.Record.Set(&enrUdpPort)

	err := enode.SignV4[T,P](&n.Record, n.PrivateKey)
	if err != nil {
		return fmt.Errorf("unable to generate ENR: %v", err)
	}
	nod, err := enode.New[P](enode.V4ID[P]{}, &n.Record)
	if err != nil {
		return fmt.Errorf("unable to create enode: %v", err)
	}
	log.Trace("simnode new", "record", n.Record)
	n.node = nod
	return nil
}

func (n *NodeConfig[T,P]) initDummyEnode() error {
	return n.initEnode(net.IPv4(127, 0, 0, 1), int(n.Port), 0)
}
