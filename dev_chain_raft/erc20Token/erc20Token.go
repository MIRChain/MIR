// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc20Token

import (
	"math/big"
	"strings"

	"github.com/pavelkrolevets/MIR-pro/common"
	ethereum "github.com/pavelkrolevets/MIR-pro"
	abi "github.com/pavelkrolevets/MIR-pro/accounts/abi"
	"github.com/pavelkrolevets/MIR-pro/event"

	"github.com/pavelkrolevets/MIR-pro/accounts/abi/bind"
	"github.com/pavelkrolevets/MIR-pro/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Erc20TokenABI is the input ABI used to generate the binding from.
const Erc20TokenABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

var Erc20TokenParsedABI, _ = abi.JSON(strings.NewReader(Erc20TokenABI))

// Erc20TokenBin is the compiled bytecode used for deploying new contracts.
var Erc20TokenBin = "0x60806040523480156200001157600080fd5b5060405162001e3b38038062001e3b833981810160405281019062000037919062000497565b848481600390816200004a91906200079e565b5080600490816200005c91906200079e565b5050506200008f818460ff16600a62000076919062000a08565b8462000083919062000a59565b6200009a60201b60201c565b505050505062000bc8565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036200010c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620001039062000b1b565b60405180910390fd5b62000120600083836200021260201b60201c565b806002600082825462000134919062000b3d565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546200018b919062000b3d565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051620001f2919062000bab565b60405180910390a36200020e600083836200021760201b60201c565b5050565b505050565b505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b62000285826200023a565b810181811067ffffffffffffffff82111715620002a757620002a66200024b565b5b80604052505050565b6000620002bc6200021c565b9050620002ca82826200027a565b919050565b600067ffffffffffffffff821115620002ed57620002ec6200024b565b5b620002f8826200023a565b9050602081019050919050565b60005b838110156200032557808201518184015260208101905062000308565b8381111562000335576000848401525b50505050565b6000620003526200034c84620002cf565b620002b0565b90508281526020810184848401111562000371576200037062000235565b5b6200037e84828562000305565b509392505050565b600082601f8301126200039e576200039d62000230565b5b8151620003b08482602086016200033b565b91505092915050565b600060ff82169050919050565b620003d181620003b9565b8114620003dd57600080fd5b50565b600081519050620003f181620003c6565b92915050565b6000819050919050565b6200040c81620003f7565b81146200041857600080fd5b50565b6000815190506200042c8162000401565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200045f8262000432565b9050919050565b620004718162000452565b81146200047d57600080fd5b50565b600081519050620004918162000466565b92915050565b600080600080600060a08688031215620004b657620004b562000226565b5b600086015167ffffffffffffffff811115620004d757620004d66200022b565b5b620004e58882890162000386565b955050602086015167ffffffffffffffff8111156200050957620005086200022b565b5b620005178882890162000386565b94505060406200052a88828901620003e0565b93505060606200053d888289016200041b565b9250506080620005508882890162000480565b9150509295509295909350565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620005b057607f821691505b602082108103620005c657620005c562000568565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620006307fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82620005f1565b6200063c8683620005f1565b95508019841693508086168417925050509392505050565b6000819050919050565b60006200067f620006796200067384620003f7565b62000654565b620003f7565b9050919050565b6000819050919050565b6200069b836200065e565b620006b3620006aa8262000686565b848454620005fe565b825550505050565b600090565b620006ca620006bb565b620006d781848462000690565b505050565b5b81811015620006ff57620006f3600082620006c0565b600181019050620006dd565b5050565b601f8211156200074e576200071881620005cc565b6200072384620005e1565b8101602085101562000733578190505b6200074b6200074285620005e1565b830182620006dc565b50505b505050565b600082821c905092915050565b6000620007736000198460080262000753565b1980831691505092915050565b60006200078e838362000760565b9150826002028217905092915050565b620007a9826200055d565b67ffffffffffffffff811115620007c557620007c46200024b565b5b620007d1825462000597565b620007de82828562000703565b600060209050601f83116001811462000816576000841562000801578287015190505b6200080d858262000780565b8655506200087d565b601f1984166200082686620005cc565b60005b82811015620008505784890151825560018201915060208501945060208101905062000829565b868310156200087057848901516200086c601f89168262000760565b8355505b6001600288020188555050505b505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008160011c9050919050565b6000808291508390505b60018511156200091357808604811115620008eb57620008ea62000885565b5b6001851615620008fb5780820291505b80810290506200090b85620008b4565b9450620008cb565b94509492505050565b6000826200092e576001905062000a01565b816200093e576000905062000a01565b8160018114620009575760028114620009625762000998565b600191505062000a01565b60ff84111562000977576200097662000885565b5b8360020a91508482111562000991576200099062000885565b5b5062000a01565b5060208310610133831016604e8410600b8410161715620009d25782820a905083811115620009cc57620009cb62000885565b5b62000a01565b620009e18484846001620008c1565b92509050818404811115620009fb57620009fa62000885565b5b81810290505b9392505050565b600062000a1582620003f7565b915062000a2283620003f7565b925062000a517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff84846200091c565b905092915050565b600062000a6682620003f7565b915062000a7383620003f7565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161562000aaf5762000aae62000885565b5b828202905092915050565b600082825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b600062000b03601f8362000aba565b915062000b108262000acb565b602082019050919050565b6000602082019050818103600083015262000b368162000af4565b9050919050565b600062000b4a82620003f7565b915062000b5783620003f7565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111562000b8f5762000b8e62000885565b5b828201905092915050565b62000ba581620003f7565b82525050565b600060208201905062000bc2600083018462000b9a565b92915050565b6112638062000bd86000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80633950935111610071578063395093511461016857806370a082311461019857806395d89b41146101c8578063a457c2d7146101e6578063a9059cbb14610216578063dd62ed3e14610246576100a9565b806306fdde03146100ae578063095ea7b3146100cc57806318160ddd146100fc57806323b872dd1461011a578063313ce5671461014a575b600080fd5b6100b6610276565b6040516100c39190610b1e565b60405180910390f35b6100e660048036038101906100e19190610bd9565b610308565b6040516100f39190610c34565b60405180910390f35b61010461032b565b6040516101119190610c5e565b60405180910390f35b610134600480360381019061012f9190610c79565b610335565b6040516101419190610c34565b60405180910390f35b610152610364565b60405161015f9190610ce8565b60405180910390f35b610182600480360381019061017d9190610bd9565b61036d565b60405161018f9190610c34565b60405180910390f35b6101b260048036038101906101ad9190610d03565b6103a4565b6040516101bf9190610c5e565b60405180910390f35b6101d06103ec565b6040516101dd9190610b1e565b60405180910390f35b61020060048036038101906101fb9190610bd9565b61047e565b60405161020d9190610c34565b60405180910390f35b610230600480360381019061022b9190610bd9565b6104f5565b60405161023d9190610c34565b60405180910390f35b610260600480360381019061025b9190610d30565b610518565b60405161026d9190610c5e565b60405180910390f35b60606003805461028590610d9f565b80601f01602080910402602001604051908101604052809291908181526020018280546102b190610d9f565b80156102fe5780601f106102d3576101008083540402835291602001916102fe565b820191906000526020600020905b8154815290600101906020018083116102e157829003601f168201915b5050505050905090565b60008061031361059f565b90506103208185856105a7565b600191505092915050565b6000600254905090565b60008061034061059f565b905061034d858285610770565b6103588585856107fc565b60019150509392505050565b60006012905090565b60008061037861059f565b905061039981858561038a8589610518565b6103949190610dff565b6105a7565b600191505092915050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6060600480546103fb90610d9f565b80601f016020809104026020016040519081016040528092919081815260200182805461042790610d9f565b80156104745780601f1061044957610100808354040283529160200191610474565b820191906000526020600020905b81548152906001019060200180831161045757829003601f168201915b5050505050905090565b60008061048961059f565b905060006104978286610518565b9050838110156104dc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104d390610ec7565b60405180910390fd5b6104e982868684036105a7565b60019250505092915050565b60008061050061059f565b905061050d8185856107fc565b600191505092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610616576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161060d90610f59565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610685576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161067c90610feb565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040516107639190610c5e565b60405180910390a3505050565b600061077c8484610518565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146107f657818110156107e8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107df90611057565b60405180910390fd5b6107f584848484036105a7565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361086b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610862906110e9565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036108da576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108d19061117b565b60405180910390fd5b6108e5838383610a7b565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508181101561096b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109629061120d565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546109fe9190610dff565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610a629190610c5e565b60405180910390a3610a75848484610a80565b50505050565b505050565b505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610abf578082015181840152602081019050610aa4565b83811115610ace576000848401525b50505050565b6000601f19601f8301169050919050565b6000610af082610a85565b610afa8185610a90565b9350610b0a818560208601610aa1565b610b1381610ad4565b840191505092915050565b60006020820190508181036000830152610b388184610ae5565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610b7082610b45565b9050919050565b610b8081610b65565b8114610b8b57600080fd5b50565b600081359050610b9d81610b77565b92915050565b6000819050919050565b610bb681610ba3565b8114610bc157600080fd5b50565b600081359050610bd381610bad565b92915050565b60008060408385031215610bf057610bef610b40565b5b6000610bfe85828601610b8e565b9250506020610c0f85828601610bc4565b9150509250929050565b60008115159050919050565b610c2e81610c19565b82525050565b6000602082019050610c496000830184610c25565b92915050565b610c5881610ba3565b82525050565b6000602082019050610c736000830184610c4f565b92915050565b600080600060608486031215610c9257610c91610b40565b5b6000610ca086828701610b8e565b9350506020610cb186828701610b8e565b9250506040610cc286828701610bc4565b9150509250925092565b600060ff82169050919050565b610ce281610ccc565b82525050565b6000602082019050610cfd6000830184610cd9565b92915050565b600060208284031215610d1957610d18610b40565b5b6000610d2784828501610b8e565b91505092915050565b60008060408385031215610d4757610d46610b40565b5b6000610d5585828601610b8e565b9250506020610d6685828601610b8e565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610db757607f821691505b602082108103610dca57610dc9610d70565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610e0a82610ba3565b9150610e1583610ba3565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115610e4a57610e49610dd0565b5b828201905092915050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b6000610eb1602583610a90565b9150610ebc82610e55565b604082019050919050565b60006020820190508181036000830152610ee081610ea4565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b6000610f43602483610a90565b9150610f4e82610ee7565b604082019050919050565b60006020820190508181036000830152610f7281610f36565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b6000610fd5602283610a90565b9150610fe082610f79565b604082019050919050565b6000602082019050818103600083015261100481610fc8565b9050919050565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b6000611041601d83610a90565b915061104c8261100b565b602082019050919050565b6000602082019050818103600083015261107081611034565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b60006110d3602583610a90565b91506110de82611077565b604082019050919050565b60006020820190508181036000830152611102816110c6565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b6000611165602383610a90565b915061117082611109565b604082019050919050565b6000602082019050818103600083015261119481611158565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b60006111f7602683610a90565b91506112028261119b565b604082019050919050565b60006020820190508181036000830152611226816111ea565b905091905056fea2646970667358221220e5117959fb3863b2a505df7bb0e74e01f0f1e1ca29b2d5e49cb9092063f0a67864736f6c634300080f0033"

// DeployErc20Token deploys a new Ethereum contract, binding an instance of Erc20Token to it.
func DeployErc20Token(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, decimals uint8, initialSupply *big.Int, owner common.Address) (common.Address, *types.Transaction, *Erc20Token, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc20TokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Erc20TokenBin), backend, name, symbol, decimals, initialSupply, owner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Erc20Token{Erc20TokenCaller: Erc20TokenCaller{contract: contract}, Erc20TokenTransactor: Erc20TokenTransactor{contract: contract}, Erc20TokenFilterer: Erc20TokenFilterer{contract: contract}}, nil
}

// Erc20Token is an auto generated Go binding around an Ethereum contract.
type Erc20Token struct {
	Erc20TokenCaller     // Read-only binding to the contract
	Erc20TokenTransactor // Write-only binding to the contract
	Erc20TokenFilterer   // Log filterer for contract events
}

// Erc20TokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type Erc20TokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20TokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Erc20TokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20TokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Erc20TokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20TokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Erc20TokenSession struct {
	Contract     *Erc20Token       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc20TokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Erc20TokenCallerSession struct {
	Contract *Erc20TokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// Erc20TokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Erc20TokenTransactorSession struct {
	Contract     *Erc20TokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// Erc20TokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type Erc20TokenRaw struct {
	Contract *Erc20Token // Generic contract binding to access the raw methods on
}

// Erc20TokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Erc20TokenCallerRaw struct {
	Contract *Erc20TokenCaller // Generic read-only contract binding to access the raw methods on
}

// Erc20TokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Erc20TokenTransactorRaw struct {
	Contract *Erc20TokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewErc20Token creates a new instance of Erc20Token, bound to a specific deployed contract.
func NewErc20Token(address common.Address, backend bind.ContractBackend) (*Erc20Token, error) {
	contract, err := bindErc20Token(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erc20Token{Erc20TokenCaller: Erc20TokenCaller{contract: contract}, Erc20TokenTransactor: Erc20TokenTransactor{contract: contract}, Erc20TokenFilterer: Erc20TokenFilterer{contract: contract}}, nil
}

// NewErc20TokenCaller creates a new read-only instance of Erc20Token, bound to a specific deployed contract.
func NewErc20TokenCaller(address common.Address, caller bind.ContractCaller) (*Erc20TokenCaller, error) {
	contract, err := bindErc20Token(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20TokenCaller{contract: contract}, nil
}

// NewErc20TokenTransactor creates a new write-only instance of Erc20Token, bound to a specific deployed contract.
func NewErc20TokenTransactor(address common.Address, transactor bind.ContractTransactor) (*Erc20TokenTransactor, error) {
	contract, err := bindErc20Token(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20TokenTransactor{contract: contract}, nil
}

// NewErc20TokenFilterer creates a new log filterer instance of Erc20Token, bound to a specific deployed contract.
func NewErc20TokenFilterer(address common.Address, filterer bind.ContractFilterer) (*Erc20TokenFilterer, error) {
	contract, err := bindErc20Token(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Erc20TokenFilterer{contract: contract}, nil
}

// bindErc20Token binds a generic wrapper to an already deployed contract.
func bindErc20Token(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc20TokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20Token *Erc20TokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc20Token.Contract.Erc20TokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20Token *Erc20TokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20Token.Contract.Erc20TokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20Token *Erc20TokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20Token.Contract.Erc20TokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20Token *Erc20TokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc20Token.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20Token *Erc20TokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20Token.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20Token *Erc20TokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20Token.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Erc20Token *Erc20TokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Erc20Token.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Erc20Token *Erc20TokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Erc20Token.Contract.Allowance(&_Erc20Token.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Erc20Token *Erc20TokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Erc20Token.Contract.Allowance(&_Erc20Token.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Erc20Token *Erc20TokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Erc20Token.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Erc20Token *Erc20TokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Erc20Token.Contract.BalanceOf(&_Erc20Token.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Erc20Token *Erc20TokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Erc20Token.Contract.BalanceOf(&_Erc20Token.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Erc20Token *Erc20TokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Erc20Token.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Erc20Token *Erc20TokenSession) Decimals() (uint8, error) {
	return _Erc20Token.Contract.Decimals(&_Erc20Token.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Erc20Token *Erc20TokenCallerSession) Decimals() (uint8, error) {
	return _Erc20Token.Contract.Decimals(&_Erc20Token.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Erc20Token *Erc20TokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Erc20Token.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Erc20Token *Erc20TokenSession) Name() (string, error) {
	return _Erc20Token.Contract.Name(&_Erc20Token.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Erc20Token *Erc20TokenCallerSession) Name() (string, error) {
	return _Erc20Token.Contract.Name(&_Erc20Token.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Erc20Token *Erc20TokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Erc20Token.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Erc20Token *Erc20TokenSession) Symbol() (string, error) {
	return _Erc20Token.Contract.Symbol(&_Erc20Token.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Erc20Token *Erc20TokenCallerSession) Symbol() (string, error) {
	return _Erc20Token.Contract.Symbol(&_Erc20Token.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Erc20Token *Erc20TokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Erc20Token.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Erc20Token *Erc20TokenSession) TotalSupply() (*big.Int, error) {
	return _Erc20Token.Contract.TotalSupply(&_Erc20Token.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Erc20Token *Erc20TokenCallerSession) TotalSupply() (*big.Int, error) {
	return _Erc20Token.Contract.TotalSupply(&_Erc20Token.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Erc20Token *Erc20TokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Token.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Erc20Token *Erc20TokenSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Token.Contract.Approve(&_Erc20Token.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Erc20Token *Erc20TokenTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Token.Contract.Approve(&_Erc20Token.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Erc20Token *Erc20TokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Erc20Token.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Erc20Token *Erc20TokenSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Erc20Token.Contract.DecreaseAllowance(&_Erc20Token.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Erc20Token *Erc20TokenTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Erc20Token.Contract.DecreaseAllowance(&_Erc20Token.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Erc20Token *Erc20TokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Erc20Token.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Erc20Token *Erc20TokenSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Erc20Token.Contract.IncreaseAllowance(&_Erc20Token.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Erc20Token *Erc20TokenTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Erc20Token.Contract.IncreaseAllowance(&_Erc20Token.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_Erc20Token *Erc20TokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Token.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_Erc20Token *Erc20TokenSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Token.Contract.Transfer(&_Erc20Token.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_Erc20Token *Erc20TokenTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Token.Contract.Transfer(&_Erc20Token.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_Erc20Token *Erc20TokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Token.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_Erc20Token *Erc20TokenSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Token.Contract.TransferFrom(&_Erc20Token.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_Erc20Token *Erc20TokenTransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Token.Contract.TransferFrom(&_Erc20Token.TransactOpts, from, to, amount)
}

// Erc20TokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Erc20Token contract.
type Erc20TokenApprovalIterator struct {
	Event *Erc20TokenApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Erc20TokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20TokenApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Erc20TokenApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Erc20TokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20TokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20TokenApproval represents a Approval event raised by the Erc20Token contract.
type Erc20TokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Erc20Token *Erc20TokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*Erc20TokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Erc20Token.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &Erc20TokenApprovalIterator{contract: _Erc20Token.contract, event: "Approval", logs: logs, sub: sub}, nil
}

var ApprovalTopicHash = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Erc20Token *Erc20TokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *Erc20TokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Erc20Token.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20TokenApproval)
				if err := _Erc20Token.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Erc20Token *Erc20TokenFilterer) ParseApproval(log types.Log) (*Erc20TokenApproval, error) {
	event := new(Erc20TokenApproval)
	if err := _Erc20Token.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20TokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Erc20Token contract.
type Erc20TokenTransferIterator struct {
	Event *Erc20TokenTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Erc20TokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20TokenTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Erc20TokenTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Erc20TokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20TokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20TokenTransfer represents a Transfer event raised by the Erc20Token contract.
type Erc20TokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Erc20Token *Erc20TokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*Erc20TokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Erc20Token.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Erc20TokenTransferIterator{contract: _Erc20Token.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

var TransferTopicHash = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Erc20Token *Erc20TokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *Erc20TokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Erc20Token.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20TokenTransfer)
				if err := _Erc20Token.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Erc20Token *Erc20TokenFilterer) ParseTransfer(log types.Log) (*Erc20TokenTransfer, error) {
	event := new(Erc20TokenTransfer)
	if err := _Erc20Token.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
