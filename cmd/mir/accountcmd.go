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
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/accounts/keystore"
	"github.com/pavelkrolevets/MIR-pro/cmd/utils"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/log"
	"gopkg.in/urfave/cli.v1"
)

var (
	walletCommand = cli.Command{
		Name:      "wallet",
		Usage:     "Manage Ethereum presale wallets",
		ArgsUsage: "",
		Category:  "ACCOUNT COMMANDS",
		Description: `
    geth wallet import /path/to/my/presale.wallet

will prompt for your password and imports your ether presale account.
It can be used non-interactively with the --password option taking a
passwordfile as argument containing the wallet password in plaintext.`,
		Subcommands: []cli.Command{
			{

				Name:      "import",
				Usage:     "Import Ethereum presale wallet",
				ArgsUsage: "<keyFile>",
				Action:    utils.MigrateFlags(importWallet),
				Category:  "ACCOUNT COMMANDS",
				Flags: []cli.Flag{
					utils.DataDirFlag,
					utils.KeyStoreDirFlag,
					utils.PasswordFileFlag,
					utils.LightKDFFlag,
					utils.CryptoSwitchFlag,
					utils.CryptoGostCurveFlag,
				},
				Description: `
	geth wallet [options] /path/to/my/presale.wallet

will prompt for your password and imports your ether presale account.
It can be used non-interactively with the --password option taking a
passwordfile as argument containing the wallet password in plaintext.`,
			},
		},
	}

	accountCommand = cli.Command{
		Name:     "account",
		Usage:    "Manage accounts",
		Category: "ACCOUNT COMMANDS",
		Description: `

Manage accounts, list all existing accounts, import a private key into a new
account, create a new account or update an existing account.

It supports interactive mode, when you are prompted for password as well as
non-interactive mode where passwords are supplied via a given password file.
Non-interactive mode is only meant for scripted use on test networks or known
safe environments.

Make sure you remember the password you gave when creating a new account (with
either new or import). Without it you are not able to unlock your account.

Note that exporting your key in unencrypted format is NOT supported.

Keys are stored under <DATADIR>/keystore.
It is safe to transfer the entire directory or the individual keys therein
between ethereum nodes by simply copying.

Make sure you backup your keys regularly.`,
		Subcommands: []cli.Command{
			{
				Name:   "list",
				Usage:  "Print summary of existing accounts",
				Action: utils.MigrateFlags(accountList),
				Flags: []cli.Flag{
					utils.DataDirFlag,
					utils.KeyStoreDirFlag,
				},
				Description: `
Print a short summary of all accounts`,
			},
			{
				Name:   "new",
				Usage:  "Create a new account",
				Action: utils.MigrateFlags(accountCreate),
				Flags: []cli.Flag{
					utils.DataDirFlag,
					utils.KeyStoreDirFlag,
					utils.PasswordFileFlag,
					utils.LightKDFFlag,
					utils.CryptoSwitchFlag,
					utils.CryptoGostCurveFlag,
				},
				Description: `
    geth account new

Creates a new account and prints the address.

The account is saved in encrypted format, you are prompted for a password.

You must remember this password to unlock your account in the future.

For non-interactive use the password can be specified with the --password flag:

Note, this is meant to be used for testing only, it is a bad idea to save your
password to file or expose in any other way.
`,
			},
			{
				Name:      "update",
				Usage:     "Update an existing account",
				Action:    utils.MigrateFlags(accountUpdate),
				ArgsUsage: "<address>",
				Flags: []cli.Flag{
					utils.DataDirFlag,
					utils.KeyStoreDirFlag,
					utils.LightKDFFlag,
					utils.CryptoSwitchFlag,
					utils.CryptoGostCurveFlag,
				},
				Description: `
    geth account update <address>

Update an existing account.

The account is saved in the newest version in encrypted format, you are prompted
for a password to unlock the account and another to save the updated file.

This same command can therefore be used to migrate an account of a deprecated
format to the newest format or change the password for an account.

For non-interactive use the password can be specified with the --password flag:

    geth account update [options] <address>

Since only one password can be given, only format update can be performed,
changing your password is only possible interactively.
`,
			},
			{
				Name:   "import",
				Usage:  "Import a private key into a new account",
				Action: utils.MigrateFlags(accountImport),
				Flags: []cli.Flag{
					utils.DataDirFlag,
					utils.KeyStoreDirFlag,
					utils.PasswordFileFlag,
					utils.LightKDFFlag,
					utils.CryptoSwitchFlag,
					utils.CryptoGostCurveFlag,
				},
				ArgsUsage: "<keyFile>",
				Description: `
    geth account import <keyfile>

Imports an unencrypted private key from <keyfile> and creates a new account.
Prints the address.

The keyfile is assumed to contain an unencrypted private key in hexadecimal format.

The account is saved in encrypted format, you are prompted for a password.

You must remember this password to unlock your account in the future.

For non-interactive use the password can be specified with the -password flag:

    geth account import [options] <keyfile>

Note:
As you can directly copy your encrypted accounts to another ethereum instance,
this import mechanism is not needed when you transfer an account between
nodes.
`,
			},
			quorumAccountPluginCommands,
		},
	}
)

func accountList(ctx *cli.Context) error {
	if ctx.GlobalString(utils.CryptoSwitchFlag.Name) != "" {
		if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost" {
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
			stack, _ := makeConfigNode[gost3410.PrivateKey, gost3410.PublicKey](ctx)
			var index int
			for _, wallet := range stack.AccountManager().Wallets() {
				for _, account := range wallet.Accounts() {
					fmt.Printf("Account #%d: {%x} %s\n", index, account.Address, &account.URL)
					index++
				}
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost_csp" {
			//TODO
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "nist" {
			stack, _ := makeConfigNode[nist.PrivateKey, nist.PublicKey](ctx)
			var index int
			for _, wallet := range stack.AccountManager().Wallets() {
				for _, account := range wallet.Accounts() {
					fmt.Printf("Account #%d: {%x} %s\n", index, account.Address, &account.URL)
					index++
				}
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "pqc" {
			//TODO
		} else {
			fmt.Errorf("wrong crypto flag")
		}
	}
	return nil
}

// tries unlocking the specified account a few times.
func unlockAccount[T crypto.PrivateKey, P crypto.PublicKey](ks *keystore.KeyStore[T,P], address string, i int, passwords []string) (accounts.Account, string) {
	account, err := utils.MakeAddress(ks, address)
	if err != nil {
		utils.Fatalf("Could not list accounts: %v", err)
	}
	for trials := 0; trials < 3; trials++ {
		prompt := fmt.Sprintf("Unlocking account %s | Attempt %d/%d", address, trials+1, 3)
		password := utils.GetPassPhraseWithList(prompt, false, i, passwords)
		err = ks.Unlock(account, password)
		if err == nil {
			log.Info("Unlocked account", "address", account.Address.Hex())
			return account, password
		}
		if err, ok := err.(*keystore.AmbiguousAddrError); ok {
			log.Info("Unlocked account", "address", account.Address.Hex())
			return ambiguousAddrRecovery(ks, err, password), password
		}
		if err != keystore.ErrDecrypt {
			// No need to prompt again if the error is not decryption-related.
			break
		}
	}
	// All trials expended to unlock account, bail out
	utils.Fatalf("Failed to unlock account %s (%v)", address, err)

	return accounts.Account{}, ""
}

func ambiguousAddrRecovery[T crypto.PrivateKey, P crypto.PublicKey](ks *keystore.KeyStore[T,P], err *keystore.AmbiguousAddrError, auth string) accounts.Account {
	fmt.Printf("Multiple key files exist for address %x:\n", err.Addr)
	for _, a := range err.Matches {
		fmt.Println("  ", a.URL)
	}
	fmt.Println("Testing your password against all of them...")
	var match *accounts.Account
	for _, a := range err.Matches {
		if err := ks.Unlock(a, auth); err == nil {
			match = &a
			break
		}
	}
	if match == nil {
		utils.Fatalf("None of the listed files could be unlocked.")
	}
	fmt.Printf("Your password unlocked %s\n", match.URL)
	fmt.Println("In order to avoid this warning, you need to remove the following duplicate key files:")
	for _, a := range err.Matches {
		if a != *match {
			fmt.Println("  ", a.URL)
		}
	}
	return *match
}

// accountCreate creates a new account into the keystore defined by the CLI flags.
func accountCreate(ctx *cli.Context) error {

	if ctx.GlobalString(utils.CryptoSwitchFlag.Name) != "" {
		if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost" {
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
			cfg := gethConfig[gost3410.PrivateKey,gost3410.PublicKey]{Node: defaultNodeConfig[gost3410.PrivateKey,gost3410.PublicKey]()}
				// Load config file.
			if file := ctx.GlobalString(configFileFlag.Name); file != "" {
				if err := loadConfig(file, &cfg); err != nil {
					utils.Fatalf("%v", err)
				}
			}
			utils.SetNodeConfig(ctx, &cfg.Node)
			scryptN, scryptP, keydir, err := cfg.Node.AccountConfig()

			if err != nil {
				utils.Fatalf("Failed to read configuration: %v", err)
			}
			password := utils.GetPassPhraseWithList("Your new account is locked with a password. Please give a password. Do not forget this password.", true, 0, utils.MakePasswordList(ctx))
			account, err := keystore.StoreKey[gost3410.PrivateKey,gost3410.PublicKey](keydir, password, scryptN, scryptP)
			fmt.Printf("\nYour new key was generated\n\n")
			fmt.Printf("Public address of the key:   %s\n", account.Address.Hex())
			fmt.Printf("Path of the secret key file: %s\n\n", account.URL.Path)
			if err != nil {
				utils.Fatalf("Failed to create account: %v", err)
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost_csp" {
			//TODO
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "nist" {
			cfg := gethConfig[nist.PrivateKey,nist.PublicKey]{Node: defaultNodeConfig[nist.PrivateKey,nist.PublicKey]()}
				// Load config file.
			if file := ctx.GlobalString(configFileFlag.Name); file != "" {
				if err := loadConfig(file, &cfg); err != nil {
					utils.Fatalf("%v", err)
				}
			}
			utils.SetNodeConfig(ctx, &cfg.Node)
			scryptN, scryptP, keydir, err := cfg.Node.AccountConfig()

			if err != nil {
				utils.Fatalf("Failed to read configuration: %v", err)
			}
			password := utils.GetPassPhraseWithList("Your new account is locked with a password. Please give a password. Do not forget this password.", true, 0, utils.MakePasswordList(ctx))
			account, err := keystore.StoreKey[nist.PrivateKey,nist.PublicKey](keydir, password, scryptN, scryptP)
			fmt.Printf("\nYour new key was generated\n\n")
			fmt.Printf("Public address of the key:   %s\n", account.Address.Hex())
			fmt.Printf("Path of the secret key file: %s\n\n", account.URL.Path)
			if err != nil {
				utils.Fatalf("Failed to create account: %v", err)
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "pqc" {
			//TODO
		} else {
			fmt.Errorf("wrong crypto flag")
		}
	}
	fmt.Printf("- You can share your public address with anyone. Others need it to interact with you.\n")
	fmt.Printf("- You must NEVER share the secret key with anyone! The key controls access to your funds!\n")
	fmt.Printf("- You must BACKUP your key file! Without the key, it's impossible to access account funds!\n")
	fmt.Printf("- You must REMEMBER your password! Without the password, it's impossible to decrypt the key!\n\n")
	return nil
}

// accountUpdate transitions an account from a previous format to the current
// one, also providing the possibility to change the pass-phrase.
func accountUpdate(ctx *cli.Context) error {
	if len(ctx.Args()) == 0 {
		utils.Fatalf("No accounts specified to update")
	}
	if ctx.GlobalString(utils.CryptoSwitchFlag.Name) != "" {
		if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost" {
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
			stack, _ := makeConfigNode[gost3410.PrivateKey,gost3410.PublicKey](ctx)
			ks := stack.AccountManager().Backends(reflect.TypeOf(&keystore.KeyStore[gost3410.PrivateKey,gost3410.PublicKey]{}))[0].(*keystore.KeyStore[gost3410.PrivateKey,gost3410.PublicKey])
		
			for _, addr := range ctx.Args() {
				account, oldPassword := unlockAccount(ks, addr, 0, nil)
				newPassword := utils.GetPassPhraseWithList("Please give a new password. Do not forget this password.", true, 0, nil)
				if err := ks.Update(account, oldPassword, newPassword); err != nil {
					utils.Fatalf("Could not update the account: %v", err)
				}
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost_csp" {
			//TODO
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "nist" {
			stack, _ := makeConfigNode[nist.PrivateKey,nist.PublicKey](ctx)
			ks := stack.AccountManager().Backends(reflect.TypeOf(&keystore.KeyStore[nist.PrivateKey,nist.PublicKey]{}))[0].(*keystore.KeyStore[nist.PrivateKey,nist.PublicKey])
		
			for _, addr := range ctx.Args() {
				account, oldPassword := unlockAccount(ks, addr, 0, nil)
				newPassword := utils.GetPassPhraseWithList("Please give a new password. Do not forget this password.", true, 0, nil)
				if err := ks.Update(account, oldPassword, newPassword); err != nil {
					utils.Fatalf("Could not update the account: %v", err)
				}
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "pqc" {
			//TODO
		} else {
			fmt.Errorf("wrong crypto flag")
		}
	}
	return nil
}

func importWallet(ctx *cli.Context) error {
		keyfile := ctx.Args().First()
		if len(keyfile) == 0 {
			utils.Fatalf("keyfile must be given as argument")
		}
		if ctx.GlobalString(utils.CryptoSwitchFlag.Name) != "" {
			if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost" {
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
				keyfile := ctx.Args().First()
				if len(keyfile) == 0 {
					utils.Fatalf("keyfile must be given as argument")
				}
				keyJSON, err := ioutil.ReadFile(keyfile)
				if err != nil {
					utils.Fatalf("Could not read wallet file: %v", err)
				}
		
				stack, _ := makeConfigNode[gost3410.PrivateKey,gost3410.PublicKey](ctx)
				passphrase := utils.GetPassPhraseWithList("", false, 0, utils.MakePasswordList(ctx))
		
				ks := stack.AccountManager().Backends(reflect.TypeOf(&keystore.KeyStore[gost3410.PrivateKey,gost3410.PublicKey]{}))[0].(*keystore.KeyStore[gost3410.PrivateKey,gost3410.PublicKey])
				acct, err := ks.ImportPreSaleKey(keyJSON, passphrase)
				if err != nil {
					utils.Fatalf("%v", err)
				}
				fmt.Printf("Address: {%x}\n", acct.Address)
			}
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "nist" { 
			keyfile := ctx.Args().First()
			if len(keyfile) == 0 {
				utils.Fatalf("keyfile must be given as argument")
			}
			keyJSON, err := ioutil.ReadFile(keyfile)
			if err != nil {
				utils.Fatalf("Could not read wallet file: %v", err)
			}

			stack, _ := makeConfigNode[nist.PrivateKey, nist.PublicKey](ctx)
			passphrase := utils.GetPassPhraseWithList("", false, 0, utils.MakePasswordList(ctx))

			ks := stack.AccountManager().Backends(reflect.TypeOf(&keystore.KeyStore[nist.PrivateKey, nist.PublicKey]{}))[0].(*keystore.KeyStore[nist.PrivateKey, nist.PublicKey])
			acct, err := ks.ImportPreSaleKey(keyJSON, passphrase)
			if err != nil {
				utils.Fatalf("%v", err)
			}
			fmt.Printf("Address: {%x}\n", acct.Address)
		}  else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost_csp" {
			//TODO
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "pqc" {
			//TODO
		} else {
			fmt.Errorf("wrong crypto flag")
		}
	return nil
}

func accountImport(ctx *cli.Context) error {
	keyfile := ctx.Args().First()
	if len(keyfile) == 0 {
		utils.Fatalf("keyfile must be given as argument")
	}
	if ctx.GlobalString(utils.CryptoSwitchFlag.Name) != "" {
		if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost" {
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

			key, err := crypto.LoadECDSA[gost3410.PrivateKey](keyfile)
			if err != nil {
				utils.Fatalf("Failed to load the private key: %v", err)
			}
			stack, _ := makeConfigNode[gost3410.PrivateKey,gost3410.PublicKey](ctx)
			passphrase := utils.GetPassPhraseWithList("Your new account is locked with a password. Please give a password. Do not forget this password.", true, 0, utils.MakePasswordList(ctx))

			ks := stack.AccountManager().Backends(reflect.TypeOf(&keystore.KeyStore[gost3410.PrivateKey,gost3410.PublicKey]{}))[0].(*keystore.KeyStore[gost3410.PrivateKey,gost3410.PublicKey])
			acct, err := ks.ImportECDSA(key, passphrase)
			if err != nil {
				utils.Fatalf("Could not create the account: %v", err)
			}
			fmt.Printf("Address: {%x}\n", acct.Address)
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "nist" { 
			key, err := crypto.LoadECDSA[nist.PrivateKey](keyfile)
			if err != nil {
				utils.Fatalf("Failed to load the private key: %v", err)
			}
			stack, _ := makeConfigNode[nist.PrivateKey,nist.PublicKey](ctx)
			passphrase := utils.GetPassPhraseWithList("Your new account is locked with a password. Please give a password. Do not forget this password.", true, 0, utils.MakePasswordList(ctx))

			ks := stack.AccountManager().Backends(reflect.TypeOf(&keystore.KeyStore[nist.PrivateKey,nist.PublicKey]{}))[0].(*keystore.KeyStore[nist.PrivateKey,nist.PublicKey])
			acct, err := ks.ImportECDSA(key, passphrase)
			if err != nil {
				utils.Fatalf("Could not create the account: %v", err)
			}
			fmt.Printf("Address: {%x}\n", acct.Address)
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost_csp" { 
			//TODO
		}else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "pqc" {
			//TODO
		} else {
			fmt.Errorf("wrong crypto flag")
		}
	}
	return nil
}
