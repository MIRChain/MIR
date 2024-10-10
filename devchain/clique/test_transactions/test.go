package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/MIRChain/MIR/accounts/abi/bind"
	"github.com/MIRChain/MIR/accounts/keystore"
	"github.com/MIRChain/MIR/crypto/nist"
	"github.com/MIRChain/MIR/devchain/clique/test_transactions/simple"
	"github.com/MIRChain/MIR/ethclient"
	"github.com/google/uuid"
)

func main() {
	DeploySipmleContract()
}

func DeploySipmleContract() {
	back, err := ethclient.Dial[nist.PublicKey]("http://127.0.0.1:8545")
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
	jsonBytes, err := ioutil.ReadFile("../node1/keystore/UTC--2018-06-24T06-41-08.065147879Z--b47f736b9b15dcc888ab790c38a6ad930217cbee")
	if err != nil {
		log.Fatal(err)
	}
	key, err := keystore.DecryptKey[nist.PrivateKey, nist.PublicKey](jsonBytes, "extreme8811")
	if err != nil {
		panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID[nist.PrivateKey, nist.PublicKey](key.PrivateKey, big.NewInt(1515))
	if err != nil {
		panic(err)
	}
	address, tx, _, err := simple.DeploySimple[nist.PublicKey](auth, back)
	if err != nil {
		panic(err)
	}

	log.Println("Contract deployed at addr: ", address.Hex())
	log.Println("Tx hash ", tx.Hash().Hex())

	receipt, err := bind.WaitMined[nist.PublicKey](ctx, back, tx)
	log.Println("Contract receipt block num : ", receipt.BlockNumber.String())

	contract, err := simple.NewSimple[nist.PublicKey](address, back)
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
	receipt, err = bind.WaitMined[nist.PublicKey](ctx, back, tx)
	if err != nil {
		panic(err)
	}
	log.Println("Value set receipt block num : ", receipt.BlockNumber.String())

	log.Println("Value set", tx.Hash().Hex())

	_value, err := contract.GetValue(&bind.CallOpts{Pending: true, Context: ctx})
	if err != nil {
		panic(err)
	}
	log.Println("Value get: ", uuid.UUID(_value))

}
