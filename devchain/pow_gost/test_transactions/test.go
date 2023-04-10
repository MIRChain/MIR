package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/google/uuid"
	"github.com/pavelkrolevets/MIR-pro/accounts/abi/bind"
	"github.com/pavelkrolevets/MIR-pro/accounts/keystore"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/devchain/clique/test_transactions/simple"
	"github.com/pavelkrolevets/MIR-pro/ethclient"
)

func main() {
	DeploySipmleContract()
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
	jsonBytes, err := ioutil.ReadFile("../node1/keystore/UTC--2023-04-06T14-15-10.038483225Z--8ac1983a8e7656a10566c4d795f3509ee35a41c3")
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

	log.Println("Value set", tx.Hash().Hex())

	_value, err := contract.GetValue(&bind.CallOpts{Pending: true, Context: ctx})
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	log.Println("Value get: ",  uuid.UUID(_value))

}
