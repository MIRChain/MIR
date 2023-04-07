package p2p

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/internal/testlog"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/p2p/enr"
)

func TestServerListenGost(t *testing.T) {
	// start the test server
	connected := make(chan *Peer[gost3410.PrivateKey, gost3410.PublicKey])
	remid := newkey[gost3410.PrivateKey]().Public()
	srv := startTestServer[gost3410.PrivateKey, gost3410.PublicKey](t, *remid, func(p *Peer[gost3410.PrivateKey, gost3410.PublicKey]) {
		if p.ID() != enode.PubkeyToIDV4(remid) {
			t.Error("peer func called with wrong node id")
		}
		connected <- p
	})
	defer close(connected)
	defer srv.Stop()

	// dial the test server
	conn, err := net.DialTimeout("tcp", srv.ListenAddr, 5*time.Second)
	if err != nil {
		t.Fatalf("could not dial: %v", err)
	}
	defer conn.Close()

	select {
	case peer := <-connected:
		if peer.LocalAddr().String() != conn.RemoteAddr().String() {
			t.Errorf("peer started with wrong conn: got %v, want %v",
				peer.LocalAddr(), conn.RemoteAddr())
		}
		peers := srv.Peers()
		if !reflect.DeepEqual(peers, []*Peer[gost3410.PrivateKey, gost3410.PublicKey]{peer}) {
			t.Errorf("Peers mismatch: got %v, want %v", peers, []*Peer[gost3410.PrivateKey, gost3410.PublicKey]{peer})
		}
	case <-time.After(1 * time.Second):
		t.Error("server did not accept within one second")
	}
}

func TestServerDialGost(t *testing.T) {
	// run a one-shot TCP server to handle the connection.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not setup listener: %v", err)
	}
	defer listener.Close()
	accepted := make(chan net.Conn, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		accepted <- conn
	}()

	// start the server
	connected := make(chan *Peer[gost3410.PrivateKey, gost3410.PublicKey])
	remid := newkey[gost3410.PrivateKey]().Public()
	srv := startTestServer[gost3410.PrivateKey, gost3410.PublicKey](t, *remid, func(p *Peer[gost3410.PrivateKey, gost3410.PublicKey]) { connected <- p })
	defer close(connected)
	defer srv.Stop()

	// tell the server to connect
	tcpAddr := listener.Addr().(*net.TCPAddr)
	node := enode.NewV4(*remid, tcpAddr.IP, tcpAddr.Port, 0)
	srv.AddPeer(node)

	select {
	case conn := <-accepted:
		defer conn.Close()

		select {
		case peer := <-connected:
			if peer.ID() != enode.PubkeyToIDV4(remid) {
				t.Errorf("peer has wrong id")
			}
			if peer.Name() != "test" {
				t.Errorf("peer has wrong name")
			}
			if peer.RemoteAddr().String() != conn.LocalAddr().String() {
				t.Errorf("peer started with wrong conn: got %v, want %v",
					peer.RemoteAddr(), conn.LocalAddr())
			}
			peers := srv.Peers()
			if !reflect.DeepEqual(peers, []*Peer[gost3410.PrivateKey, gost3410.PublicKey]{peer}) {
				t.Errorf("Peers mismatch: got %v, want %v", peers, []*Peer[gost3410.PrivateKey, gost3410.PublicKey]{peer})
			}

			// Test AddTrustedPeer/RemoveTrustedPeer and changing Trusted flags
			// Particularly for race conditions on changing the flag state.
			if peer := srv.Peers()[0]; peer.Info().Network.Trusted {
				t.Errorf("peer is trusted prematurely: %v", peer)
			}
			done := make(chan bool)
			go func() {
				srv.AddTrustedPeer(node)
				if peer := srv.Peers()[0]; !peer.Info().Network.Trusted {
					t.Errorf("peer is not trusted after AddTrustedPeer: %v", peer)
				}
				srv.RemoveTrustedPeer(node)
				if peer := srv.Peers()[0]; peer.Info().Network.Trusted {
					t.Errorf("peer is trusted after RemoveTrustedPeer: %v", peer)
				}
				done <- true
			}()
			// Trigger potential race conditions
			peer = srv.Peers()[0]
			_ = peer.Inbound()
			_ = peer.Info()
			<-done
		case <-time.After(1 * time.Second):
			t.Error("server did not launch peer within one second")
		}

	case <-time.After(1 * time.Second):
		t.Error("server did not connect within one second")
	}
}

// This test checks that RemovePeer disconnects the peer if it is connected.
func TestServerRemovePeerDisconnectGost(t *testing.T) {
	srv1 := &Server[gost3410.PrivateKey, gost3410.PublicKey]{Config: Config[gost3410.PrivateKey, gost3410.PublicKey]{
		PrivateKey:  newkey[gost3410.PrivateKey](),
		MaxPeers:    1,
		NoDiscovery: true,
		Logger:      testlog.Logger(t, log.LvlTrace).New("server", "1"),
	}}
	srv2 := &Server[gost3410.PrivateKey, gost3410.PublicKey]{Config: Config[gost3410.PrivateKey, gost3410.PublicKey]{
		PrivateKey:  newkey[gost3410.PrivateKey](),
		MaxPeers:    1,
		NoDiscovery: true,
		NoDial:      true,
		ListenAddr:  "127.0.0.1:0",
		Logger:      testlog.Logger(t, log.LvlTrace).New("server", "2"),
	}}
	srv1.Start()
	defer srv1.Stop()
	srv2.Start()
	defer srv2.Stop()

	if !syncAddPeer(srv1, srv2.Self()) {
		t.Fatal("peer not connected")
	}
	srv1.RemovePeer(srv2.Self())
	if srv1.PeerCount() > 0 {
		t.Fatal("removed peer still connected")
	}
}

// This test checks that connections are disconnected just after the encryption handshake
// when the server is at capacity. Trusted connections should still be accepted.
func TestServerAtCapGost(t *testing.T) {
	trustedNode := newkey[gost3410.PrivateKey]()
	trustedID := enode.PubkeyToIDV4(trustedNode.Public())
	srv := &Server[gost3410.PrivateKey, gost3410.PublicKey]{
		Config: Config[gost3410.PrivateKey, gost3410.PublicKey]{
			PrivateKey:   newkey[gost3410.PrivateKey](),
			MaxPeers:     10,
			NoDial:       true,
			NoDiscovery:  true,
			TrustedNodes: []*enode.Node[gost3410.PublicKey]{newNode[gost3410.PrivateKey, gost3410.PublicKey](trustedID, "")},
			Logger:       testlog.Logger(t, log.LvlTrace),
		},
	}
	if err := srv.Start(); err != nil {
		t.Fatalf("could not start: %v", err)
	}
	defer srv.Stop()

	newconn := func(id enode.ID) *conn[gost3410.PrivateKey, gost3410.PublicKey] {
		fd, _ := net.Pipe()
		tx := newTestTransport[gost3410.PrivateKey,gost3410.PublicKey](*trustedNode.Public(), fd, gost3410.PublicKey{})
		node := enode.SignNull[gost3410.PrivateKey, gost3410.PublicKey](new(enr.Record), id)
		return &conn[gost3410.PrivateKey, gost3410.PublicKey]{fd: fd, transport: tx, flags: inboundConn, node: node, cont: make(chan error)}
	}

	// Inject a few connections to fill up the peer set.
	for i := 0; i < 10; i++ {
		c := newconn(randomID())
		if err := srv.checkpoint(c, srv.checkpointAddPeer); err != nil {
			t.Fatalf("could not add conn %d: %v", i, err)
		}
	}
	// Try inserting a non-trusted connection.
	anotherID := randomID()
	c := newconn(anotherID)
	if err := srv.checkpoint(c, srv.checkpointPostHandshake); err != DiscTooManyPeers {
		t.Error("wrong error for insert:", err)
	}
	// Try inserting a trusted connection.
	c = newconn(trustedID)
	if err := srv.checkpoint(c, srv.checkpointPostHandshake); err != nil {
		t.Error("unexpected error for trusted conn @posthandshake:", err)
	}
	if !c.is(trustedConn) {
		t.Error("Server did not set trusted flag")
	}

	// Remove from trusted set and try again
	srv.RemoveTrustedPeer(newNode[gost3410.PrivateKey, gost3410.PublicKey](trustedID, ""))
	c = newconn(trustedID)
	if err := srv.checkpoint(c, srv.checkpointPostHandshake); err != DiscTooManyPeers {
		t.Error("wrong error for insert:", err)
	}

	// Add anotherID to trusted set and try again
	srv.AddTrustedPeer(newNode[gost3410.PrivateKey, gost3410.PublicKey](anotherID, ""))
	c = newconn(anotherID)
	if err := srv.checkpoint(c, srv.checkpointPostHandshake); err != nil {
		t.Error("unexpected error for trusted conn @posthandshake:", err)
	}
	if !c.is(trustedConn) {
		t.Error("Server did not set trusted flag")
	}
}

func TestServerPeerLimitsGost(t *testing.T) {
	srvkey := newkey[gost3410.PrivateKey]()
	clientkey := newkey[gost3410.PrivateKey]()
	clientnode := enode.NewV4(*clientkey.Public(), nil, 0, 0)

	var tp = &setupTransport[gost3410.PrivateKey,gost3410.PublicKey]{
		pubkey: *clientkey.Public(),
		phs: protoHandshake{
			ID: crypto.FromECDSAPub(*clientkey.Public())[1:],
			// Force "DiscUselessPeer" due to unmatching caps
			// Caps: []Cap{discard.cap()},
		},
	}

	srv := &Server[gost3410.PrivateKey, gost3410.PublicKey]{
		Config: Config[gost3410.PrivateKey, gost3410.PublicKey]{
			PrivateKey:  srvkey,
			MaxPeers:    0,
			NoDial:      true,
			NoDiscovery: true,
			Protocols:   []Protocol[gost3410.PrivateKey, gost3410.PublicKey]{
				Protocol[gost3410.PrivateKey,gost3410.PublicKey]{
				Name:   "discard",
				Length: 1,
				Run: func(p *Peer[gost3410.PrivateKey, gost3410.PublicKey], rw MsgReadWriter) error {
					for {
						msg, err := rw.ReadMsg()
						if err != nil {
							return err
						}
						fmt.Printf("discarding %d\n", msg.Code)
						if err = msg.Discard(); err != nil {
							return err
						}
					}
				},
			},
		},
			Logger:      testlog.Logger(t, log.LvlTrace),
		},
		newTransport: func(fd net.Conn, dialDest gost3410.PublicKey) transport[gost3410.PrivateKey, gost3410.PublicKey] { return tp },
	}
	if err := srv.Start(); err != nil {
		t.Fatalf("couldn't start server: %v", err)
	}
	defer srv.Stop()

	// Check that server is full (MaxPeers=0)
	flags := dynDialedConn
	dialDest := clientnode
	conn, _ := net.Pipe()
	srv.SetupConn(conn, flags, dialDest)
	if tp.closeErr != DiscTooManyPeers {
		t.Errorf("unexpected close error: %q", tp.closeErr)
	}
	conn.Close()

	srv.AddTrustedPeer(clientnode)

	// Check that server allows a trusted peer despite being full.
	conn, _ = net.Pipe()
	srv.SetupConn(conn, flags, dialDest)
	if tp.closeErr == DiscTooManyPeers {
		t.Errorf("failed to bypass MaxPeers with trusted node: %q", tp.closeErr)
	}

	if tp.closeErr != DiscUselessPeer {
		t.Errorf("unexpected close error: %q", tp.closeErr)
	}
	conn.Close()

	srv.RemoveTrustedPeer(clientnode)

	// Check that server is full again.
	conn, _ = net.Pipe()
	srv.SetupConn(conn, flags, dialDest)
	if tp.closeErr != DiscTooManyPeers {
		t.Errorf("unexpected close error: %q", tp.closeErr)
	}
	conn.Close()
}

func TestServerSetupConnGost(t *testing.T) {
	var (
		clientkey, srvkey = newkey[gost3410.PrivateKey](), newkey[gost3410.PrivateKey]()
		clientpub         = clientkey.Public()
		srvpub            = srvkey.Public()
	)
	tests := []struct {
		dontstart bool
		tt        *setupTransport[gost3410.PrivateKey,gost3410.PublicKey]
		flags     connFlag
		dialDest  *enode.Node[gost3410.PublicKey]

		wantCloseErr error
		wantCalls    string
	}{
		{
			dontstart:    true,
			tt:           &setupTransport[gost3410.PrivateKey,gost3410.PublicKey]{pubkey: *clientpub},
			wantCalls:    "close,",
			wantCloseErr: errServerStopped,
		},
		{
			tt:           &setupTransport[gost3410.PrivateKey,gost3410.PublicKey]{pubkey: *clientpub, encHandshakeErr: errors.New("read error")},
			flags:        inboundConn,
			wantCalls:    "doEncHandshake,close,",
			wantCloseErr: errors.New("read error"),
		},
		{
			tt:           &setupTransport[gost3410.PrivateKey,gost3410.PublicKey]{pubkey: *clientpub, phs: protoHandshake{ID: randomID().Bytes()}},
			dialDest:     enode.NewV4(*clientpub, nil, 0, 0),
			flags:        dynDialedConn,
			wantCalls:    "doEncHandshake,doProtoHandshake,close,",
			wantCloseErr: DiscUnexpectedIdentity,
		},
		{
			tt:           &setupTransport[gost3410.PrivateKey,gost3410.PublicKey]{pubkey: *clientpub, protoHandshakeErr: errors.New("foo")},
			dialDest:     enode.NewV4(*clientpub, nil, 0, 0),
			flags:        dynDialedConn,
			wantCalls:    "doEncHandshake,doProtoHandshake,close,",
			wantCloseErr: errors.New("foo"),
		},
		{
			tt:           &setupTransport[gost3410.PrivateKey,gost3410.PublicKey]{pubkey: *srvpub, phs: protoHandshake{ID: crypto.FromECDSAPub(*srvpub)[1:]}},
			flags:        inboundConn,
			wantCalls:    "doEncHandshake,close,",
			wantCloseErr: DiscSelf,
		},
		{
			tt:           &setupTransport[gost3410.PrivateKey,gost3410.PublicKey]{pubkey: *clientpub, phs: protoHandshake{ID: crypto.FromECDSAPub(*clientpub)[1:]}},
			flags:        inboundConn,
			wantCalls:    "doEncHandshake,doProtoHandshake,close,",
			wantCloseErr: DiscUselessPeer,
		},
	}

	for i, test := range tests {
		t.Run(test.wantCalls, func(t *testing.T) {
			cfg := Config[gost3410.PrivateKey,gost3410.PublicKey]{
				PrivateKey:  srvkey,
				MaxPeers:    10,
				NoDial:      true,
				NoDiscovery: true,
				Protocols:   []Protocol[gost3410.PrivateKey,gost3410.PublicKey]{
					Protocol[gost3410.PrivateKey,gost3410.PublicKey]{
						Name:   "discard",
						Length: 1,
						Run: func(p *Peer[gost3410.PrivateKey, gost3410.PublicKey], rw MsgReadWriter) error {
							for {
								msg, err := rw.ReadMsg()
								if err != nil {
									return err
								}
								fmt.Printf("discarding %d\n", msg.Code)
								if err = msg.Discard(); err != nil {
									return err
								}
							}
						},
					},
				},
				Logger:      testlog.Logger(t, log.LvlTrace),
			}
			srv := &Server[gost3410.PrivateKey,gost3410.PublicKey]{
				Config:       cfg,
				newTransport: func(fd net.Conn, dialDest gost3410.PublicKey) transport[gost3410.PrivateKey,gost3410.PublicKey] { return test.tt },
				log:          cfg.Logger,
			}
			if !test.dontstart {
				if err := srv.Start(); err != nil {
					t.Fatalf("couldn't start server: %v", err)
				}
				defer srv.Stop()
			}
			p1, _ := net.Pipe()
			srv.SetupConn(p1, test.flags, test.dialDest)
			if !reflect.DeepEqual(test.tt.closeErr, test.wantCloseErr) {
				t.Errorf("test %d: close error mismatch: got %q, want %q", i, test.tt.closeErr, test.wantCloseErr)
			}
			if test.tt.calls != test.wantCalls {
				t.Errorf("test %d: calls mismatch: got %q, want %q", i, test.tt.calls, test.wantCalls)
			}
		})
	}
}
