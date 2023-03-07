package main

import (

	"context"
	"io/ioutil"
	"log"
	"math/big"


	"github.com/pavelkrolevets/MIR-pro/accounts/abi/bind"
	"github.com/pavelkrolevets/MIR-pro/accounts/keystore"
	"github.com/pavelkrolevets/MIR-pro/dev_chain/simple"
	"github.com/pavelkrolevets/MIR-pro/dev_chain/token"
	"github.com/pavelkrolevets/MIR-pro/dev_chain/erc20Token"
	"github.com/pavelkrolevets/MIR-pro/ethclient"
)

func main() {

	DeploySipmleContract()
	DeployFactory()
}


func DeploySipmleContract(){


	back, err := ethclient.Dial("http://127.0.0.1:8545")
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
	key, err := keystore.DecryptKey(jsonBytes, password)
	if err != nil {
		log.Panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, big.NewInt(123456789))
	if err != nil {
		log.Panic(err)
	}
	

	address, tx, _, err := simple.DeploySimple(auth, back)
	if err != nil {
		panic(err)
	}
	
	log.Println("Contract deployed at addr: ", address.Hex())
	log.Println("Tx hash ",tx.Hash().Hex())
	
	receipt, err := bind.WaitMined(ctx, back, tx)
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
	receipt, err = bind.WaitMined(ctx, back, tx)
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


func DeployFactory(){


	keyset := "../../node1/keystore/UTC--2022-08-17T13-11-15.885898644Z--25face3782641d4ad4c7eaf2ee5eb8a7dcfab465"

	jsonBytes, err := ioutil.ReadFile(keyset)
	if err != nil {
	log.Fatal(err)
	}
	password := ""
	key, err := keystore.DecryptKey(jsonBytes, password)
	if err != nil {
		log.Panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, big.NewInt(123456789))
	if err != nil {
		log.Panic(err)
	}

	back, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Panic(err)
	}

	address, tx, instance, err := mch_factory.DeployMchFactory(auth, back)
	if err != nil {
		panic(err)
	}
	
	log.Println("Contract deployed at addr: ", address.Hex())
	log.Println("Tx hash ",tx.Hash().Hex())
	
	receipt, err := bind.WaitMined(context.Background(), back, tx)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Contract receipt block num : ", receipt.BlockNumber.String())

	tx, err =instance.DeployNewERC20Token(auth, "ERC20Factory", "MCH20", 2, big.NewInt(100000))
	if err != nil {
		log.Panic(err)
	}

	sink := make(chan *mch_factory.MchFactoryERC20TokenCreated)
	event, err := instance.WatchERC20TokenCreated(&bind.WatchOpts{Context: context.Background()}, sink)
	if err != nil {
		log.Panic(err)
	}
	for {
		select {
		case err := <-event.Err():
		  log.Fatal(err)
		case erc20 := <-sink:
			log.Println("Contract ERC20 deployed : ", erc20.TokenAddress.Hex())
			erc20instance, err := erc20Token.NewErc20Token(erc20.TokenAddress, back)
			if err != nil {
				log.Panic(err)
			}
			balance, err := erc20instance.BalanceOf(&bind.CallOpts{Pending: true, Context: context.Background()}, key.Address)
			if err != nil {
				log.Panic(err)
			}
			log.Println("Balance of deployed ERC20 : ", balance.String())
		}
	  }
}


