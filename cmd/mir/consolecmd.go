// Copyright 2016 The go-ethereum Authors
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

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pavelkrolevets/MIR-pro/cmd/utils"
	"github.com/pavelkrolevets/MIR-pro/console"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/node"
	"github.com/pavelkrolevets/MIR-pro/plugin/security"
	"github.com/pavelkrolevets/MIR-pro/rpc"
	"gopkg.in/urfave/cli.v1"
)

var (
	consoleFlags = []cli.Flag{utils.JSpathFlag, utils.ExecFlag, utils.PreloadJSFlag}

	rpcClientFlags = []cli.Flag{utils.RPCClientToken, utils.RPCClientTLSCert, utils.RPCClientTLSCaCert, utils.RPCClientTLSCipherSuites, utils.RPCClientTLSInsecureSkipVerify}

	consoleCommand = cli.Command{
		Action:   utils.MigrateFlags(_localConsole),
		Name:     "console",
		Usage:    "Start an interactive JavaScript environment",
		Flags:    append(append(nodeFlags, rpcFlags...), consoleFlags...),
		Category: "CONSOLE COMMANDS",
		Description: `
The Geth console is an interactive shell for the JavaScript runtime environment
which exposes a node admin interface as well as the Ðapp JavaScript API.
See https://geth.ethereum.org/docs/interface/javascript-console.`,
	}

	attachCommand = cli.Command{
		Action:    utils.MigrateFlags(_remoteConsole),
		Name:      "attach",
		Usage:     "Start an interactive JavaScript environment (connect to node)",
		ArgsUsage: "[endpoint]",
		Flags:     append(append(consoleFlags, utils.DataDirFlag), rpcClientFlags...),
		Category:  "CONSOLE COMMANDS",
		Description: `
The Geth console is an interactive shell for the JavaScript runtime environment
which exposes a node admin interface as well as the Ðapp JavaScript API.
See https://geth.ethereum.org/docs/interface/javascript-console.
This command allows to open a console on a running geth node.`,
	}

	javascriptCommand = cli.Command{
		Action:    utils.MigrateFlags(_ephemeralConsole),
		Name:      "js",
		Usage:     "Execute the specified JavaScript files",
		ArgsUsage: "<jsfile> [jsfile...]",
		Flags:     append(nodeFlags, consoleFlags...),
		Category:  "CONSOLE COMMANDS",
		Description: `
The JavaScript VM exposes a node admin interface as well as the Ðapp
JavaScript API. See https://geth.ethereum.org/docs/interface/javascript-console`,
	}
)

// Quorum
//
// read tls client configuration from command line arguments
//
// only for HTTPS or WSS
func readTLSClientConfig(endpoint string, ctx *cli.Context) (*tls.Config, bool, error) {
	if !strings.HasPrefix(endpoint, "https://") && !strings.HasPrefix(endpoint, "wss://") {
		return nil, false, nil
	}
	hasCustomTls := false
	insecureSkipVerify := ctx.Bool(utils.RPCClientTLSInsecureSkipVerify.Name)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	}
	var certFile, caFile string
	if !insecureSkipVerify {
		var certPem, caPem []byte
		certFile, caFile = ctx.String(utils.RPCClientTLSCert.Name), ctx.String(utils.RPCClientTLSCaCert.Name)
		var err error
		if certFile != "" {
			if certPem, err = ioutil.ReadFile(certFile); err != nil {
				return nil, true, err
			}
		}
		if caFile != "" {
			if caPem, err = ioutil.ReadFile(caFile); err != nil {
				return nil, true, err
			}
		}
		if len(certPem) != 0 || len(caPem) != 0 {
			certPool, err := x509.SystemCertPool()
			if err != nil {
				certPool = x509.NewCertPool()
			}
			if len(certPem) != 0 {
				certPool.AppendCertsFromPEM(certPem)
			}
			if len(caPem) != 0 {
				certPool.AppendCertsFromPEM(caPem)
			}
			tlsConfig.RootCAs = certPool
			hasCustomTls = true
		}
	} else {
		hasCustomTls = true
	}
	cipherSuitesInput := ctx.String(utils.RPCClientTLSCipherSuites.Name)
	cipherSuitesStrings := strings.FieldsFunc(cipherSuitesInput, func(r rune) bool {
		return r == ','
	})
	if len(cipherSuitesStrings) > 0 {
		cipherSuiteList := make(security.CipherSuiteList, len(cipherSuitesStrings))
		for i, s := range cipherSuitesStrings {
			cipherSuiteList[i] = security.CipherSuite(strings.TrimSpace(s))
		}
		cipherSuites, err := cipherSuiteList.ToUint16Array()
		if err != nil {
			return nil, true, err
		}
		tlsConfig.CipherSuites = cipherSuites
		hasCustomTls = true
	}
	if !hasCustomTls {
		return nil, false, nil
	}
	return tlsConfig, hasCustomTls, nil
}

// localConsole starts a new geth node, attaching a JavaScript console to it at the
// same time.
func _localConsole(ctx *cli.Context) error {
	if ctx.GlobalString(utils.CryptoSwitchFlag.Name) != "" {
		if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost" {
			fmt.Println(`
			╔═╗┬─┐┬ ┬┌─┐┌┬┐┌─┐  ╔═╗╔═╗╔═╗╔╦╗
			║  ├┬┘└┬┘├─┘ │ │ │  ║ ╦║ ║╚═╗ ║ 
			╚═╝┴└─ ┴ ┴   ┴ └─┘  ╚═╝╚═╝╚═╝ ╩ `)
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-tc26-gost-3410-12-256-paramSetA" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetA()
			}
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-tc26-gost-3410-12-256-paramSetB" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetB()
			}
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-tc26-gost-3410-12-256-paramSetC" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetC()
			}
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-GostR3410-2001-CryptoPro-A-ParamSet" {
				gost3410.GostCurve = gost3410.CurveIdGostR34102001CryptoProAParamSet()
			}
			err := localConsole[gost3410.PrivateKey, gost3410.PublicKey](ctx)
			if err != nil {
				return err
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost_csp" {
			// Mir - check if signer cert is loaded
			// if stack.Config().SignerCert.Bytes() == nil {
			// 	return fmt.Errorf("signer cert cant be nil")
			// }
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "nist" {
			fmt.Println(`
			╔═╗┬─┐┬ ┬┌─┐┌┬┐┌─┐  ╔╗╔╦╔═╗╔╦╗
			║  ├┬┘└┬┘├─┘ │ │ │  ║║║║╚═╗ ║ 
			╚═╝┴└─ ┴ ┴   ┴ └─┘  ╝╚╝╩╚═╝ ╩ `)
			err := localConsole[nist.PrivateKey, nist.PublicKey](ctx)
			if err != nil {
				return err
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "pqc" {
			crypto.CryptoAlg = crypto.PQC
		} else {
			fmt.Errorf("wrong crypto flag")
		}
	} 	
	return nil		
}
func localConsole[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	// Create and start the node based on the CLI flags
	prepare(ctx)
	stack, backend := makeFullNode[T,P](ctx)
	startNode(ctx, stack, backend)
	defer stack.Close()

	// Attach to the newly started node and start the JavaScript console
	client, err := stack.Attach()
	if err != nil {
		utils.Fatalf("Failed to attach to the inproc geth: %v", err)
	}
	config := console.Config{
		DataDir: utils.MakeDataDir(ctx),
		DocRoot: ctx.GlobalString(utils.JSpathFlag.Name),
		Client:  client,
		Preload: utils.MakeConsolePreloads(ctx),
	}

	console, err := console.New(config)
	if err != nil {
		utils.Fatalf("Failed to start the JavaScript console: %v", err)
	}
	defer console.Stop(false)

	// If only a short execution was requested, evaluate and return
	if script := ctx.GlobalString(utils.ExecFlag.Name); script != "" {
		console.Evaluate(script)
		return nil
	}
	// Otherwise print the welcome screen and enter interactive mode
	console.Welcome()
	console.Interactive()

	return nil
}

func _remoteConsole(ctx *cli.Context) error {
	if ctx.GlobalString(utils.CryptoSwitchFlag.Name) != "" {
		if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost" {
			fmt.Println(`
			╔═╗┬─┐┬ ┬┌─┐┌┬┐┌─┐  ╔═╗╔═╗╔═╗╔╦╗
			║  ├┬┘└┬┘├─┘ │ │ │  ║ ╦║ ║╚═╗ ║ 
			╚═╝┴└─ ┴ ┴   ┴ └─┘  ╚═╝╚═╝╚═╝ ╩ `)
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-tc26-gost-3410-12-256-paramSetA" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetA()
			}
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-tc26-gost-3410-12-256-paramSetB" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetB()
			}
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-tc26-gost-3410-12-256-paramSetC" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetC()
			}
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-GostR3410-2001-CryptoPro-A-ParamSet" {
				gost3410.GostCurve = gost3410.CurveIdGostR34102001CryptoProAParamSet()
			}
			err := remoteConsole[gost3410.PrivateKey, gost3410.PublicKey](ctx)
			if err != nil {
				return err
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost_csp" {
			// Mir - check if signer cert is loaded
			// if stack.Config().SignerCert.Bytes() == nil {
			// 	return fmt.Errorf("signer cert cant be nil")
			// }
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "nist" {
			fmt.Println(`
			╔═╗┬─┐┬ ┬┌─┐┌┬┐┌─┐  ╔╗╔╦╔═╗╔╦╗
			║  ├┬┘└┬┘├─┘ │ │ │  ║║║║╚═╗ ║ 
			╚═╝┴└─ ┴ ┴   ┴ └─┘  ╝╚╝╩╚═╝ ╩ `)
			err := remoteConsole[nist.PrivateKey, nist.PublicKey](ctx)
			if err != nil {
				return err
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "pqc" {
			crypto.CryptoAlg = crypto.PQC
		} else {
			fmt.Errorf("wrong crypto flag")
		}
	} 	
	return nil	
}
// remoteConsole will connect to a remote geth instance, attaching a JavaScript
// console to it.
func remoteConsole[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	// Attach to a remotely running geth instance and start the JavaScript console
	endpoint := ctx.Args().First()
	if endpoint == "" {
		path := node.DefaultDataDir()
		if ctx.GlobalIsSet(utils.DataDirFlag.Name) {
			path = ctx.GlobalString(utils.DataDirFlag.Name)
		}
		if path != "" {
			if ctx.GlobalBool(utils.SoyuzFlag.Name) {
				// Maintain compatibility with older Mir configurations storing the
				// Soyuz database in `testnet` instead of `soyuz`.
				legacyPath := filepath.Join(path, "testnet")
				if _, err := os.Stat(legacyPath); !os.IsNotExist(err) {
					path = legacyPath
				} else {
					path = filepath.Join(path, "soyuz")
				}
			}
		}
		endpoint = fmt.Sprintf("%s/mir.ipc", path)
	}
	client, err := dialRPC[T,P](endpoint, ctx)
	if err != nil {
		utils.Fatalf("Unable to attach to remote mir: %v", err)
	}
	config := console.Config{
		DataDir: utils.MakeDataDir(ctx),
		DocRoot: ctx.GlobalString(utils.JSpathFlag.Name),
		Client:  client,
		Preload: utils.MakeConsolePreloads(ctx),
	}

	console, err := console.New(config)
	if err != nil {
		utils.Fatalf("Failed to start the JavaScript console: %v", err)
	}
	defer console.Stop(false)

	if script := ctx.GlobalString(utils.ExecFlag.Name); script != "" {
		console.Evaluate(script)
		return nil
	}

	// Otherwise print the welcome screen and enter interactive mode
	console.Welcome()
	console.Interactive()

	return nil
}

// dialRPC returns a RPC client which connects to the given endpoint.
// The check for empty endpoint implements the defaulting logic
// for "geth attach" and "geth monitor" with no argument.
//
// Quorum: passing the cli context to build security-aware client:
// 1. Custom TLS configuration
// 2. Access Token awareness via rpc.HttpCredentialsProviderFunc
// 3. PSI awareness from environment variable and endpoint query param
func dialRPC[T crypto.PrivateKey, P crypto.PublicKey](endpoint string, ctx *cli.Context) (*rpc.Client, error) {
	if endpoint == "" {
		endpoint = node.DefaultIPCEndpoint[T,P](clientIdentifier)
	} else if strings.HasPrefix(endpoint, "rpc:") || strings.HasPrefix(endpoint, "ipc:") {
		// Backwards compatibility with geth < 1.5 which required
		// these prefixes.
		endpoint = endpoint[4:]
	}
	var (
		client  *rpc.Client
		err     error
		dialCtx = context.Background()
	)
	tlsConfig, hasCustomTls, tlsReadErr := readTLSClientConfig(endpoint, ctx)
	if tlsReadErr != nil {
		return nil, tlsReadErr
	}
	if token := ctx.String(utils.RPCClientToken.Name); token != "" {
		var f rpc.HttpCredentialsProviderFunc = func(ctx context.Context) (string, error) {
			return token, nil
		}
		// it's important that f MUST BE OF TYPE rpc.HttpCredentialsProviderFunc
		dialCtx = rpc.WithCredentialsProvider(dialCtx, f)
	}
	if hasCustomTls {
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		switch u.Scheme {
		case "https":
			customHttpClient := &http.Client{
				Transport: http.DefaultTransport,
			}
			customHttpClient.Transport.(*http.Transport).TLSClientConfig = tlsConfig
			client, _ = rpc.DialHTTPWithClient(endpoint, customHttpClient)
		case "wss":
			client, _ = rpc.DialWebsocketWithCustomTLS(dialCtx, endpoint, "", tlsConfig)
		default:
			log.Warn("unsupported scheme for custom TLS which is only for HTTPS/WSS", "scheme", u.Scheme)
			client, _ = rpc.DialContext(dialCtx, endpoint)
		}
	} else {
		client, err = rpc.DialContext(dialCtx, endpoint)
	}
	if err != nil {
		return nil, err
	}
	// enrich clients with provider functions to populate HTTP request header
	if f := rpc.CredentialsProviderFromContext(dialCtx); f != nil {
		client = client.WithHTTPCredentials(f)
	}
	return client, nil
}

func _ephemeralConsole(ctx *cli.Context) error {
	if ctx.GlobalString(utils.CryptoSwitchFlag.Name) != "" {
		if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost" {
			fmt.Println(`
			╔═╗┬─┐┬ ┬┌─┐┌┬┐┌─┐  ╔═╗╔═╗╔═╗╔╦╗
			║  ├┬┘└┬┘├─┘ │ │ │  ║ ╦║ ║╚═╗ ║ 
			╚═╝┴└─ ┴ ┴   ┴ └─┘  ╚═╝╚═╝╚═╝ ╩ `)
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-tc26-gost-3410-12-256-paramSetA" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetA()
			}
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-tc26-gost-3410-12-256-paramSetB" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetB()
			}
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-tc26-gost-3410-12-256-paramSetC" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetC()
			}
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-GostR3410-2001-CryptoPro-A-ParamSet" {
				gost3410.GostCurve = gost3410.CurveIdGostR34102001CryptoProAParamSet()
			}
			err := ephemeralConsole[gost3410.PrivateKey, gost3410.PublicKey](ctx)
			if err != nil {
				return err
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost_csp" {
			// Mir - check if signer cert is loaded
			// if stack.Config().SignerCert.Bytes() == nil {
			// 	return fmt.Errorf("signer cert cant be nil")
			// }
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "nist" {
			fmt.Println(`
			╔═╗┬─┐┬ ┬┌─┐┌┬┐┌─┐  ╔╗╔╦╔═╗╔╦╗
			║  ├┬┘└┬┘├─┘ │ │ │  ║║║║╚═╗ ║ 
			╚═╝┴└─ ┴ ┴   ┴ └─┘  ╝╚╝╩╚═╝ ╩ `)
			err := ephemeralConsole[nist.PrivateKey, nist.PublicKey](ctx)
			if err != nil {
				return err
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "pqc" {
			crypto.CryptoAlg = crypto.PQC
		} else {
			fmt.Errorf("wrong crypto flag")
		}
	} 	
	return nil	
}

// ephemeralConsole starts a new geth node, attaches an ephemeral JavaScript
// console to it, executes each of the files specified as arguments and tears
// everything down.
func ephemeralConsole[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	// Create and start the node based on the CLI flags
	stack, backend := makeFullNode[T,P](ctx)
	startNode(ctx, stack, backend)
	defer stack.Close()

	// Attach to the newly started node and start the JavaScript console
	client, err := stack.Attach()
	if err != nil {
		utils.Fatalf("Failed to attach to the inproc geth: %v", err)
	}
	config := console.Config{
		DataDir: utils.MakeDataDir(ctx),
		DocRoot: ctx.GlobalString(utils.JSpathFlag.Name),
		Client:  client,
		Preload: utils.MakeConsolePreloads(ctx),
	}

	console, err := console.New(config)
	if err != nil {
		utils.Fatalf("Failed to start the JavaScript console: %v", err)
	}
	defer console.Stop(false)

	// Evaluate each of the specified JavaScript files
	for _, file := range ctx.Args() {
		if err = console.Execute(file); err != nil {
			utils.Fatalf("Failed to execute %s: %v", file, err)
		}
	}

	go func() {
		stack.Wait()
		console.Stop(false)
	}()
	console.Stop(true)

	return nil
}
