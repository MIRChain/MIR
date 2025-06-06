// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// bootnode runs a bootstrap node for the Ethereum Discovery Protocol.
package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/pavelkrolevets/MIR-pro/cmd/utils"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/p2p/discover"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/p2p/nat"
	"github.com/pavelkrolevets/MIR-pro/p2p/netutil"
)

func main() {
	var (
		listenAddr  = flag.String("addr", ":30301", "listen address")
		genKey      = flag.String("genkey", "", "generate a node key")
		writeAddr   = flag.Bool("writeaddress", false, "write out the node's public key and quit")
		nodeKeyFile = flag.String("nodekey", "", "private key filename")
		nodeKeyHex  = flag.String("nodekeyhex", "", "private key as hex (for testing)")
		natdesc     = flag.String("nat", "none", "port mapping mechanism (any|none|upnp|pmp|extip:<IP>)")
		netrestrict = flag.String("netrestrict", "", "restrict network communication to the given IP networks (CIDR masks)")
		runv5       = flag.Bool("v5", false, "run a v5 topic discovery bootnode")
		verbosity   = flag.Int("verbosity", int(log.LvlInfo), "log verbosity (0-5)")
		vmodule     = flag.String("vmodule", "", "log verbosity pattern")
		cryptoType  = flag.String("crypto", "nist", "switch between differet tyoes of cryptography - NIST, GOST, PostQuantum")
		gostCurve   = flag.String("gostcurve", "id-GostR3410-2001-CryptoPro-A-ParamSet", "GOST ECDSA curve parameters")
		// nodeKey *ecdsa.PrivateKey
		err     error
	)
	flag.Parse()

	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(*verbosity))
	glogger.Vmodule(*vmodule)
	log.Root().SetHandler(glogger)

	if *cryptoType == "nist" || *cryptoType == "gost" || *cryptoType == "gost_csp" ||  *cryptoType == "pqc" {
		if *cryptoType == "nist"{
			runNode[nist.PrivateKey, nist.PublicKey](listenAddr, genKey, nodeKeyFile, nodeKeyHex, natdesc, netrestrict, writeAddr, runv5, err)
		}
		if *cryptoType == "gost"{
			if *gostCurve == "" || *gostCurve == "id-GostR3410-2001-CryptoPro-A-ParamSet"{
				gost3410.GostCurve = gost3410.CurveIdGostR34102001CryptoProAParamSet()
			} else if *gostCurve == "id-tc26-gost-3410-12-256-paramSetC" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetC()
			} else if *gostCurve == "id-tc26-gost-3410-12-256-paramSetB" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetB()
			} else if *gostCurve == "id-tc26-gost-3410-12-256-paramSetA" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetA()
			}
			runNode[gost3410.PrivateKey, gost3410.PublicKey](listenAddr, genKey, nodeKeyFile, nodeKeyHex, natdesc, netrestrict, writeAddr, runv5, err)
		}
	} else {
		panic("Crypto type should be set: nist, gost, gost_csp, pqc")
	}
	select {}
}

func printNotice[P crypto.PublicKey](nodeKey P, addr net.UDPAddr) {
	if addr.IP.IsUnspecified() {
		addr.IP = net.IP{127, 0, 0, 1}
	}
	n := enode.NewV4(nodeKey, addr.IP, 0, addr.Port)
	fmt.Println(n.URLv4())
	fmt.Println("Note: you're using cmd/bootnode, a developer tool.")
	fmt.Println("We recommend using a regular node as bootstrap node for production deployments.")
}

func runNode[T crypto.PrivateKey, P crypto.PublicKey](listenAddr, genKey, nodeKeyFile, nodeKeyHex, natdesc, netrestrict *string, writeAddr, runv5 *bool, err error){
	var nodeKey T
	natm, err := nat.Parse(*natdesc)
	if err != nil {
		utils.Fatalf("-nat: %v", err)
	}
	switch {
	case *genKey != "":
		nodeKey, err = crypto.GenerateKey[T]()
		if err != nil {
			utils.Fatalf("could not generate key: %v", err)
		}
		if err = crypto.SaveECDSA(*genKey, nodeKey); err != nil {
			utils.Fatalf("%v", err)
		}
		if !*writeAddr {
			return
		}
	case *nodeKeyFile == "" && *nodeKeyHex == "":
		utils.Fatalf("Use -nodekey or -nodekeyhex to specify a private key")
	case *nodeKeyFile != "" && *nodeKeyHex != "":
		utils.Fatalf("Options -nodekey and -nodekeyhex are mutually exclusive")
	case *nodeKeyFile != "":
		if nodeKey, err = crypto.LoadECDSA[T](*nodeKeyFile); err != nil {
			utils.Fatalf("-nodekey: %v", err)
		}
	case *nodeKeyHex != "":
		if nodeKey, err = crypto.HexToECDSA[T](*nodeKeyHex); err != nil {
			utils.Fatalf("-nodekeyhex: %v", err)
		}
	}
	var pub P
	switch t:=any(&nodeKey).(type) {
	case *nist.PrivateKey:
		p:=any(&pub).(*nist.PublicKey)
		*p = *t.Public()
	case *gost3410.PrivateKey:
		p:=any(&pub).(*gost3410.PublicKey)
		*p = *t.Public()
	}
	if *writeAddr {
		fmt.Printf("%x\n", crypto.FromECDSAPub(pub)[1:])
		os.Exit(0)
	}
	var restrictList *netutil.Netlist
	if *netrestrict != "" {
		restrictList, err = netutil.ParseNetlist(*netrestrict)
		if err != nil {
			utils.Fatalf("-netrestrict: %v", err)
		}
	}

	addr, err := net.ResolveUDPAddr("udp", *listenAddr)
	if err != nil {
		utils.Fatalf("-ResolveUDPAddr: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		utils.Fatalf("-ListenUDP: %v", err)
	}

	realaddr := conn.LocalAddr().(*net.UDPAddr)
	if natm != nil {
		if !realaddr.IP.IsLoopback() {
			go nat.Map(natm, nil, "udp", realaddr.Port, realaddr.Port, "ethereum discovery")
		}
		if ext, err := natm.ExternalIP(); err == nil {
			realaddr = &net.UDPAddr{IP: ext, Port: realaddr.Port}
		}
	}
	printNotice(pub, *realaddr)
	db, _ := enode.OpenDB[P]("")
	ln := enode.NewLocalNode(db, nodeKey)
	cfg := discover.Config[T,P]{
		PrivateKey:  nodeKey,
		NetRestrict: restrictList,
	}
	if *runv5 {
		if _, err := discover.ListenV5(conn, ln, cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	} else {
		if _, err := discover.ListenUDP(conn, ln, cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	}
}