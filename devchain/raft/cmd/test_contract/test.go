package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/pavelkrolevets/MIR-pro/accounts/abi/bind"
	"github.com/pavelkrolevets/MIR-pro/accounts/keystore"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/devchain/raft/simple"
	"github.com/pavelkrolevets/MIR-pro/ethclient"
)

func main() {

	DeploySipmleContract()
}


func DeploySipmleContract(){


	back, err := ethclient.Dial[nist.PublicKey]("http://127.0.0.1:8545")
	if err != nil {
		log.Panic(err)
	}
	ctx := context.Background()
	blockNumber, err := back.BlockNumber(ctx)
	log.Println("Block number ", blockNumber)
	blockInfo, err := back.BlockByNumber(ctx, nil)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Block hash ", blockInfo.Hash().Hex())

	keyset := "../../node1/keystore/UTC--2022-08-17T13-11-15.885898644Z--25face3782641d4ad4c7eaf2ee5eb8a7dcfab465"
	jsonBytes, err := ioutil.ReadFile(keyset)
	if err != nil {
	log.Fatal(err)
	}
	password := "Gfdtk81,"
	key, err := keystore.DecryptKey[nist.PrivateKey,nist.PublicKey](jsonBytes, password)
	if err != nil {
		log.Panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID[nist.PrivateKey,nist.PublicKey](key.PrivateKey, big.NewInt(123456789))
	if err != nil {
		log.Panic(err)
	}
	

	address, tx, _, err := simple.DeploySimple(auth, back)
	if err != nil {
		panic(err)
	}
	
	log.Println("Contract deployed at addr: ", address.Hex())
	log.Println("Tx hash ",tx.Hash().Hex())
	
	receipt, err := bind.WaitMined[nist.PublicKey](ctx, back, tx)
	log.Println("Contract receipt block num : ", receipt.BlockNumber.String())

	contract, err := simple.NewSimple(address, back)
	if err != nil {
		log.Panic(err)
	}

	value, err := contract.GetValue(&bind.CallOpts{Pending: true, Context: ctx})
	if err != nil {
		log.Panic(err)
	}
	log.Println("Value get: ",  value.String())

	tx, err = contract.SetValue(auth, big.NewInt(777))
	if err != nil {
		log.Panic(err)
	}
	receipt, err = bind.WaitMined[nist.PublicKey](ctx, back, tx)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Value set receipt block num : ", receipt.BlockNumber.String())

	log.Println("Value set",  tx.Hash().Hex())

	value, err = contract.GetValue(&bind.CallOpts{Pending: true, Context: ctx})
	if err != nil {
		log.Panic(err)
	}
	log.Println("Value get: ",  value.String())

}

