package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/google/uuid"

	"github.com/MIRChain/MIR/accounts/abi"
	"github.com/MIRChain/MIR/accounts/abi/bind"
	"github.com/MIRChain/MIR/accounts/keystore"
	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/crypto/csp"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/devchain/clique_gost/test_transactions/simple"
	"github.com/MIRChain/MIR/ethclient"
)

func main() {
	DeploySipmleContractGOST()
	// TestTransferFromCSPtoGost()
	DeploySipmleContractCSP()
}

func TestTransfer() {
	back, err := ethclient.Dial[gost3410.PublicKey]("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	blockNumber, err := back.BlockNumber(ctx)
	if err != nil {
		panic(err)
	}
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
	key, err := keystore.DecryptKey[gost3410.PrivateKey, gost3410.PublicKey](jsonBytes, "12345678")
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

	toAddress := common.HexToAddress("0x07514B71D17497f81008319Fc5Eb8E69A07f2DBe")
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

	_, err = bind.WaitMined[gost3410.PublicKey](context.Background(), back, signedTx)
	if err != nil {
		log.Fatal(err)
	}
	balanceCheck, err := back.BalanceAt(context.Background(), toAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Account balance %s", balanceCheck.String())
}

func TestTransferFromCSPtoGost() {
	store, err := csp.SystemStore("My")
	if err != nil {
		log.Fatalf("Store error: %s", err)
	}
	defer store.Close()
	crt, err := store.GetBySubjectId("71732462bbc029d911e6d16a3ed00d9d1d772620")
	if err != nil {
		log.Fatalf("Get cert error: %s", err)
	}
	log.Printf("Cert pub key: %x", crt.Info().PublicKeyBytes()[2:66])
	defer crt.Close()
	back, err := ethclient.Dial[gost3410.PublicKey]("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	blockNumber, err := back.BlockNumber(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Block number ", blockNumber)
	nonce, err := back.PendingNonceAt(context.Background(), common.HexToAddress("0x07514B71D17497f81008319Fc5Eb8E69A07f2DBe"))
	if err != nil {
		log.Fatal(err)
	}
	value := big.NewInt(10000000) // in wei (1 eth)
	gasLimit := uint64(21000)     // in units
	gasPrice, err := back.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("gas price ", gasPrice.String())
	balanceCheck, err := back.BalanceAt(context.Background(), common.HexToAddress("0x37B3e5679521f5A652c2C0738EDaa0700C8cCC0a"), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Account balance From %s", balanceCheck.String())

	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	var data []byte
	tx := types.NewTransaction[gost3410.PublicKey](nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := back.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Chain ID ", chainID)

	// Sign TX low level
	signer := types.NewEIP155Signer[gost3410.PublicKey](chainID)
	txHash := signer.Hash(tx)
	// Sign with cert
	sig, err := crypto.Sign(txHash[:], crt)
	if err != nil {
		log.Fatal(err)
	}
	r, s, _ := crypto.RevertCSP(txHash[:], sig)
	resSig := make([]byte, 65)
	copy(resSig[:32], r.Bytes())
	copy(resSig[32:64], s.Bytes())
	resSig[64] = sig[64]

	// Get address which will be used at pure GO GOST network - recover to get the right value
	recoveredGostPub, err := crypto.Ecrecover[gost3410.PublicKey](txHash[:], resSig)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Account From %x", crypto.Keccak256[gost3410.PublicKey](recoveredGostPub[1:])[12:])
	resPub := recoveredGostPub[1:]
	pub := make([]byte, 64)
	copy(pub[:32], resPub[32:64])
	copy(pub[32:64], resPub[:32])
	reverse(pub)
	// compare recovered and crt.pub
	if !bytes.Equal(crt.Info().PublicKeyBytes()[2:66], pub) {
		log.Fatal("Wrong revovered pub key")
	}

	signedTx, err := tx.WithSignature(signer, resSig)
	if err != nil {
		log.Fatal(err)
	}
	err = back.SendTransaction(context.Background(), signedTx, bind.PrivateTxArgs{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = bind.WaitMined[gost3410.PublicKey](context.Background(), back, signedTx)
	if err != nil {
		log.Fatal(err)
	}
	balanceCheck, err = back.BalanceAt(context.Background(), toAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Account balance %s", balanceCheck.String())
}

func DeploySipmleContractGOST() {
	gost3410.GostCurve = gost3410.CurveIdGostR34102001CryptoProAParamSet()
	back, err := ethclient.Dial[gost3410.PublicKey]("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	blockNumber, err := back.BlockNumber(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Block number ", blockNumber)
	blockInfo, err := back.BlockByNumber(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Block hash ", blockInfo.Hash().Hex())
	jsonBytes, err := ioutil.ReadFile("../node1/keystore/UTC--2023-04-11T19-26-02.261401864Z--1f1a2f8231efe45f0ff0c2ed3c27eebf58fc175c")
	if err != nil {
		log.Fatal(err)
	}
	key, err := keystore.DecryptKey[gost3410.PrivateKey, gost3410.PublicKey](jsonBytes, "12345678")
	if err != nil {
		panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID[gost3410.PrivateKey, gost3410.PublicKey](key.PrivateKey, big.NewInt(1515))
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
	if err != nil {
		panic(err)
	}
	log.Println("Contract receipt block num : ", receipt.BlockNumber.String())

	contract, err := simple.NewSimple[gost3410.PublicKey](address, back)
	if err != nil {
		panic(err)
	}

	value, err := contract.GetValue(&bind.CallOpts{Pending: true, Context: ctx})
	if err != nil {
		panic(err)
	}
	log.Println("Value get: ", uuid.UUID(value))

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
	log.Println("Value get: ", uuid.UUID(_value))

}

func DeploySipmleContractCSP() {
	/*
		pragma solidity ^0.8.0;
		contract Simple {
		    uint256 public value;
		    function setValue(uint256 v) public {
		        value = v;
		    }
		    function getValue() public view returns (uint256) {
		        return value;
		    }
		}
	*/
	const contractAbi = "[{\"inputs\":[],\"name\":\"getValue\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"v\",\"type\":\"bytes16\"}],\"name\":\"setValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"value\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	const contractBin = "0x608060405234801561001057600080fd5b50610206806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806307ea02e1146100465780637d7ee50e14610064578063c7ff210314610082575b600080fd5b61004e61009e565b60405161005b919061012c565b60405180910390f35b61006c6100b4565b604051610079919061012c565b60405180910390f35b61009c60048036038101906100979190610178565b6100c5565b005b60008060009054906101000a900460801b905090565b60008054906101000a900460801b81565b806000806101000a8154816fffffffffffffffffffffffffffffffff021916908360801c021790555050565b60007fffffffffffffffffffffffffffffffff0000000000000000000000000000000082169050919050565b610126816100f1565b82525050565b6000602082019050610141600083018461011d565b92915050565b600080fd5b610155816100f1565b811461016057600080fd5b50565b6000813590506101728161014c565b92915050565b60006020828403121561018e5761018d610147565b5b600061019c84828501610163565b9150509291505056fea2646970667358221220a94d387a337f7be406d713b9d25400b18463d3def57c02777bfb3af704806ec864736f6c63782d302e382e31342d646576656c6f702e323032342e31312e32362b636f6d6d69742e38306434396633372e6d6f64005e"

	store, err := csp.SystemStore("My")
	if err != nil {
		log.Fatalf("Store error: %s", err)
	}
	defer store.Close()
	crt, err := store.GetBySubjectId("71732462bbc029d911e6d16a3ed00d9d1d772620")
	if err != nil {
		log.Fatalf("Get cert error: %s", err)
	}
	log.Printf("Cert pub key: %x", crt.Info().PublicKeyBytes()[2:66])
	defer crt.Close()

	back, err := ethclient.Dial[gost3410.PublicKey]("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	blockNumber, err := back.BlockNumber(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Block number ", blockNumber)
	blockInfo, err := back.BlockByNumber(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Block hash ", blockInfo.Hash().Hex())
	addrCSP := crypto.PubkeyToAddress[csp.PublicKey](*crt.Public())
	log.Printf("Cert address CSP %s", addrCSP.String())
	certPub := crt.Info().PublicKeyBytes()[2:66]
	reverse(certPub)
	pub := make([]byte, 65)
	copy(pub[33:65], certPub[:32])
	copy(pub[1:33], certPub[32:64])
	var addrFrom common.Address
	copy(addrFrom[:], crypto.Keccak256[gost3410.PublicKey](pub[1:])[12:])
	log.Printf("Cert address gost3410 %x", addrFrom.Bytes())

	var gas uint64 = 3000000
	nonce, err := back.NonceAt(ctx, addrFrom, nil)
	if err != nil {
		panic(err)
	}
	gasPrice, err := back.SuggestGasPrice(ctx)
	if err != nil {
		panic(err)
	}
	tx := types.NewContractCreation[gost3410.PublicKey](nonce, big.NewInt(0), gas, gasPrice, common.FromHex(contractBin))
	// Sign TX low level
	signer := types.NewEIP155Signer[gost3410.PublicKey](big.NewInt(1515))
	txHash := signer.Hash(tx)
	sig, err := crypto.Sign(txHash[:], crt)
	if err != nil {
		panic(err)
	}

	r, s, _ := crypto.RevertCSP(txHash[:], sig)
	resSig := make([]byte, 65)
	copy(resSig[:32], r.Bytes())
	copy(resSig[32:64], s.Bytes())
	resSig[64] = sig[64]

	// Get address which will be used at pure GO GOST network - recover to get the right value
	recoveredGostPub, err := crypto.Ecrecover[gost3410.PublicKey](txHash[:], resSig)
	if err != nil {
		panic(err)
	}

	recoveredPub := make([]byte, 64)
	copy(recoveredPub[:32], recoveredGostPub[33:65])
	copy(recoveredPub[32:64], recoveredGostPub[1:33])
	reverse(recoveredPub)
	// compare recovered and crt.pub
	if !bytes.Equal(crt.Info().PublicKeyBytes()[2:66], recoveredPub) {
		panic("Wrong recovered pub key")
	}

	signedTx, err := tx.WithSignature(signer, resSig)
	if err != nil {
		panic(err)
	}
	err = back.SendTransaction(context.Background(), signedTx, bind.PrivateTxArgs{})
	if err != nil {
		panic("error sending transaction")
	}

	txHash = signedTx.Hash()

	log.Println("Tx hash ", txHash.Hex())

	receipt, err := bind.WaitMined[gost3410.PublicKey](ctx, back, signedTx)
	if err != nil {
		panic(err)
	}
	log.Println("Contract receipt block num : ", receipt.ContractAddress.Hex())

	contract, err := simple.NewSimple[gost3410.PublicKey](receipt.ContractAddress, back)
	if err != nil {
		panic(err)
	}

	value, err := contract.GetValue(&bind.CallOpts{Pending: true, Context: ctx})
	if err != nil {
		panic(err)
	}
	log.Println("Value get: ", uuid.UUID(value))

	parsed, err := abi.JSON[gost3410.PublicKey](strings.NewReader(contractAbi))
	if err != nil {
		panic(err)
	}
	// Set value at contract
	input, err := parsed.Pack("setValue", uuid.New())
	if err != nil {
		panic(err)
	}
	nonce, err = back.NonceAt(ctx, addrFrom, nil)
	if err != nil {
		panic(err)
	}
	tx = types.NewTransaction[gost3410.PublicKey](nonce, receipt.ContractAddress, big.NewInt(0), gas, big.NewInt(1), input)
	txHash = signer.Hash(tx)
	sig, err = crypto.Sign(txHash[:], crt)
	if err != nil {
		panic(err)
	}

	r, s, _ = crypto.RevertCSP(txHash[:], sig)
	resSig = make([]byte, 65)
	copy(resSig[:32], r.Bytes())
	copy(resSig[32:64], s.Bytes())
	resSig[64] = sig[64]

	signedTx, err = tx.WithSignature(signer, resSig)
	if err != nil {
		panic(err)
	}

	err = back.SendTransaction(context.Background(), signedTx, bind.PrivateTxArgs{})
	if err != nil {
		panic("error sending transaction")
	}

	receipt, err = bind.WaitMined[gost3410.PublicKey](ctx, back, signedTx)
	if err != nil {
		panic(err)
	}
	log.Println("Value set receipt block num : ", receipt.BlockNumber.String())
	_value, err := contract.GetValue(&bind.CallOpts{Pending: true, Context: ctx})
	if err != nil {
		panic(err)
	}
	log.Println("Value get: ", uuid.UUID(_value))

}

func reverse(d []byte) {
	for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
		d[i], d[j] = d[j], d[i]
	}
}
