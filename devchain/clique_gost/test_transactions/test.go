package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/google/uuid"
	"github.com/pavelkrolevets/MIR-pro/accounts/abi/bind"
	"github.com/pavelkrolevets/MIR-pro/accounts/keystore"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto/csp"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/devchain/clique/test_transactions/simple"
	"github.com/pavelkrolevets/MIR-pro/ethclient"
)

func main() {
	// TestTransfer()
	// DeploySipmleContract()
	TestTransferFromCSP()
}


func TestTransfer(){
	back, err := ethclient.Dial[gost3410.PublicKey]("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	blockNumber, err := back.BlockNumber(ctx)
	log.Println("Block number ", blockNumber)
	blockInfo, err := back.BlockByNumber(ctx, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Block hash ", blockInfo.Hash().Hex())
	jsonBytes, err := ioutil.ReadFile("../node1/keystore/UTC--2023-04-11T19-26-02.261401864Z--1f1a2f8231efe45f0ff0c2ed3c27eebf58fc175c")
    if err != nil {
        log.Fatal(err)
    }
	key, err := keystore.DecryptKey[gost3410.PrivateKey,gost3410.PublicKey](jsonBytes, "12345678")
	if err != nil {
		panic(err)
	}
	nonce, err := back.PendingNonceAt(context.Background(), key.Address)
    if err != nil {
        log.Fatal(err)
    }
	value := big.NewInt(1000000000000000000) // in wei (1 eth)
    gasLimit := uint64(21000)                // in units
    gasPrice, err := back.SuggestGasPrice(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    toAddress := common.HexToAddress("0xE9D1230d60718938cE7bB47aE434b0caC9e354F6")
    var data []byte
    tx := types.NewTransaction[gost3410.PublicKey](nonce, toAddress, value, gasLimit, gasPrice, data)

    chainID, err := back.ChainID(context.Background())
    if err != nil {
        log.Fatal(err)
    }
	log.Println("Chain ID ", chainID)
    signedTx, err := types.SignTx[gost3410.PrivateKey, gost3410.PublicKey](tx, types.NewEIP155Signer[gost3410.PublicKey](chainID), key.PrivateKey)
    if err != nil {
        log.Fatal(err)
    }

    err = back.SendTransaction(context.Background(), signedTx, bind.PrivateTxArgs{})
    if err != nil {
        log.Fatal(err)
    }
	
	_, err =bind.WaitMined[gost3410.PublicKey](context.Background(), back, signedTx)
	if err != nil {
	log.Fatal(err)
	}
	balanceCheck, err := back.BalanceAt(context.Background(), toAddress, nil)
	if err != nil {
		log.Fatal(err)
		}
	log.Printf("Account balance %s", balanceCheck.String())
}

func TestTransferFromCSP(){
	back, err := ethclient.Dial[csp.PublicKey]("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	blockNumber, err := back.BlockNumber(ctx)
	if err != nil {
        log.Fatal(err)
    }
	log.Println("Block number ", blockNumber)
	ks := keystore.NewPlaintextKeyStore[csp.Cert, csp.PublicKey]("../node1/keystore")
	a, err := ks.NewAccountCsp("4ac93fc08bc0efd24180b0fa47f7309c257e8c85")
	if err != nil {
        log.Fatal(err)
    }
	nonce, err := back.PendingNonceAt(context.Background(), a.Address)
    if err != nil {
        log.Fatal(err)
    }
	value := big.NewInt(1000000000000000000) // in wei (1 eth)
    gasLimit := uint64(21000)                // in units
    gasPrice, err := back.SuggestGasPrice(context.Background())
    if err != nil {
        log.Fatal(err)
    }
	balanceCheck, err := back.BalanceAt(context.Background(), a.Address, nil)
	if err != nil {
		log.Fatal(err)
		}
	log.Printf("Account balance From %s", balanceCheck.String())

    toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
    var data []byte
    tx := types.NewTransaction[csp.PublicKey](nonce, toAddress, value, gasLimit, gasPrice, data)

    chainID, err := back.ChainID(context.Background())
    if err != nil {
        log.Fatal(err)
    }
	log.Println("Chain ID ", chainID)
	
    signedTx, err :=ks.SignTxWithPassphraseCsp(a, "4ac93fc08bc0efd24180b0fa47f7309c257e8c85", "", tx, chainID)
    if err != nil {
        log.Fatal(err)
    }

    err = back.SendTransaction(context.Background(), signedTx, bind.PrivateTxArgs{})
    if err != nil {
        log.Fatal(err)
    }
	
	_, err =bind.WaitMined[csp.PublicKey](context.Background(), back, signedTx)
	if err != nil {
	log.Fatal(err)
	}
	balanceCheck, err = back.BalanceAt(context.Background(), toAddress, nil)
	if err != nil {
		log.Fatal(err)
		}
	log.Printf("Account balance %s", balanceCheck.String())
}

func DeploySipmleContract(){
	gost3410.GostCurve = gost3410.CurveIdGostR34102001CryptoProAParamSet()
	back, err := ethclient.Dial[gost3410.PublicKey]("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	blockNumber, err := back.BlockNumber(ctx)
	log.Println("Block number ", blockNumber)
	blockInfo, err := back.BlockByNumber(ctx, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Block hash ", blockInfo.Hash().Hex())
	jsonBytes, err := ioutil.ReadFile("../node1/keystore/UTC--2023-04-11T19-26-02.261401864Z--1f1a2f8231efe45f0ff0c2ed3c27eebf58fc175c")
    if err != nil {
        log.Fatal(err)
    }
	key, err := keystore.DecryptKey[gost3410.PrivateKey,gost3410.PublicKey](jsonBytes, "12345678")
	if err != nil {
		panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID[gost3410.PrivateKey,gost3410.PublicKey](key.PrivateKey, big.NewInt(1515))
	if err != nil {
		panic(err)
	}
	address, tx, _, err := simple.DeploySimple[gost3410.PublicKey](auth, back)
	if err != nil {
		panic(err)
	}

	log.Println("Contract deployed at addr: ", address.Hex())
	log.Println("Tx hash ", tx.Hash().Hex())

	receipt, err := bind.WaitMined[gost3410.PublicKey](ctx, back, tx)
	log.Println("Contract receipt block num : ", receipt.BlockNumber.String())

	contract, err := simple.NewSimple[gost3410.PublicKey](address, back)
	if err != nil {
		panic(err)
	}

	value, err := contract.GetValue(&bind.CallOpts{Pending: true, Context: ctx})
	if err != nil {
		panic(err)
	}
	log.Println("Value get: ", value[:])

	tx, err = contract.SetValue(auth, uuid.New())
	if err != nil {
		panic(err)
	}
	receipt, err = bind.WaitMined[gost3410.PublicKey](ctx, back, tx)
	if err != nil {
		panic(err)
	}
	log.Println("Value set receipt block num : ", receipt.BlockNumber.String())

	_value, err := contract.GetValue(&bind.CallOpts{Pending: true, Context: ctx})
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	log.Println("Value get: ",  uuid.UUID(_value))

}
