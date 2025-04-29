// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifierV2

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BN254G1Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

// BN254G2Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// BatchHeaderV2 is an auto generated low-level Go binding around an user-defined struct.
type BatchHeaderV2 struct {
	BatchRoot            [32]byte
	ReferenceBlockNumber uint32
}

// BlobCertificate is an auto generated low-level Go binding around an user-defined struct.
type BlobCertificate struct {
	BlobHeader BlobHeaderV2
	Signature  []byte
	RelayKeys  []uint32
}

// BlobCommitment is an auto generated low-level Go binding around an user-defined struct.
type BlobCommitment struct {
	Commitment       BN254G1Point
	LengthCommitment BN254G2Point
	LengthProof      BN254G2Point
	Length           uint32
}

// BlobHeaderV2 is an auto generated low-level Go binding around an user-defined struct.
type BlobHeaderV2 struct {
	Version           uint16
	QuorumNumbers     []byte
	Commitment        BlobCommitment
	PaymentHeaderHash [32]byte
}

// BlobInclusionInfo is an auto generated low-level Go binding around an user-defined struct.
type BlobInclusionInfo struct {
	BlobCertificate BlobCertificate
	BlobIndex       uint32
	InclusionProof  []byte
}

// EigenDACertV2 is an auto generated low-level Go binding around an user-defined struct.
type EigenDACertV2 struct {
	BatchHeader                 BatchHeaderV2
	BlobInclusionInfo           BlobInclusionInfo
	NonSignerStakesAndSignature NonSignerStakesAndSignature
	SignedQuorumNumbers         []byte
}

// NonSignerStakesAndSignature is an auto generated low-level Go binding around an user-defined struct.
type NonSignerStakesAndSignature struct {
	NonSignerQuorumBitmapIndices []uint32
	NonSignerPubkeys             []BN254G1Point
	QuorumApks                   []BN254G1Point
	ApkG2                        BN254G2Point
	Sigma                        BN254G1Point
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
}

// SecurityThresholds is an auto generated low-level Go binding around an user-defined struct.
type SecurityThresholds struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}

// ContractEigenDACertVerifierV2MetaData contains all meta data concerning the ContractEigenDACertVerifierV2 contract.
var ContractEigenDACertVerifierV2MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistryV2\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDASignatureVerifierV2\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_registryCoordinatorV2\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_securityThresholdsV2\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"_quorumNumbersRequiredV2\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkDACert\",\"inputs\":[{\"name\":\"cert\",\"type\":\"tuple\",\"internalType\":\"structEigenDACertV2\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifierV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequiredV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinatorV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2\",\"inputs\":[{\"name\":\"cert\",\"type\":\"tuple\",\"internalType\":\"structEigenDACertV2\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"BlobQuorumsNotSubset\",\"inputs\":[{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInclusionProof\",\"inputs\":[{\"name\":\"blobIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidSecurityThresholds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RequiredQuorumsNotSubset\",\"inputs\":[{\"name\":\"requiredQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SecurityAssumptionsNotMet\",\"inputs\":[{\"name\":\"gamma\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"n\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	Bin: "0x60e06040523480156200001157600080fd5b506040516200215d3803806200215d833981016040819052620000349162000275565b816020015160ff16826000015160ff161162000063576040516308a6997560e01b815260040160405180910390fd5b6001600160a01b0380861660805284811660a052831660c05281516000805460208086015160ff9081166101000261ffff199093169416939093171790558151620000b59160019190840190620000c1565b50505050505062000395565b828054620000cf9062000358565b90600052602060002090601f016020900481019282620000f357600085556200013e565b82601f106200010e57805160ff19168380011785556200013e565b828001600101855582156200013e579182015b828111156200013e57825182559160200191906001019062000121565b506200014c92915062000150565b5090565b5b808211156200014c576000815560010162000151565b6001600160a01b03811681146200017d57600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620001c157620001c162000180565b604052919050565b805160ff81168114620001db57600080fd5b919050565b600082601f830112620001f257600080fd5b81516001600160401b038111156200020e576200020e62000180565b602062000224601f8301601f1916820162000196565b82815285828487010111156200023957600080fd5b60005b83811015620002595785810183015182820184015282016200023c565b838111156200026b5760008385840101525b5095945050505050565b600080600080600085870360c08112156200028f57600080fd5b86516200029c8162000167565b6020880151909650620002af8162000167565b6040880151909550620002c28162000167565b93506040605f1982011215620002d757600080fd5b50604080519081016001600160401b038082118383101715620002fe57620002fe62000180565b816040526200031060608a01620001c9565b83526200032060808a01620001c9565b602084015260a0890151929450808311156200033b57600080fd5b50506200034b88828901620001e0565b9150509295509295909350565b600181811c908216806200036d57607f821691505b602082108114156200038f57634e487b7160e01b600052602260045260246000fd5b50919050565b60805160a05160c051611d7e620003df60003960006101150152600081816087015281816101ba015261027901526000818160cb0152818161019901526102580152611d7e6000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c80635fafa4821161005b5780635fafa4821461011057806399d923e714610137578063b74d78711461014c578063ed0450ae1461016157600080fd5b8063154b9e861461008257806317f3578e146100c65780635653c730146100ed575b600080fd5b6100a97f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020015b60405180910390f35b6100a97f000000000000000000000000000000000000000000000000000000000000000081565b6101006100fb366004610fb2565b610191565b60405190151581526020016100bd565b6100a97f000000000000000000000000000000000000000000000000000000000000000081565b61014a610145366004610fb2565b610253565b005b6101546102e0565b6040516100bd9190611040565b6000546101779060ff8082169161010090041682565b6040805160ff9384168152929091166020830152016100bd565b60008061021e7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000856102116040805180820182526000808252602091820181905282518084019093525460ff8082168452610100909104169082015290565b61021961036e565b610400565b509050600081600481111561023557610235611053565b14156102445750600192915050565b50600092915050565b50919050565b6102dd7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000836102d06040805180820182526000808252602091820181905282518084019093525460ff8082168452610100909104169082015290565b6102d861036e565b610606565b50565b600180546102ed90611069565b80601f016020809104026020016040519081016040528092919081815260200182805461031990611069565b80156103665780601f1061033b57610100808354040283529160200191610366565b820191906000526020600020905b81548152906001019060200180831161034957829003601f168201915b505050505081565b60606001805461037d90611069565b80601f01602080910402602001604051908101604052809291908181526020018280546103a990611069565b80156103f65780601f106103cb576101008083540402835291602001916103f6565b820191906000526020600020905b8154815290600101906020018083116103d957829003601f168201915b5050505050905090565b6000606061043161041636879003870187611172565b61042360408801886111aa565b61042c90611391565b61062d565b9092509050600082600481111561044a5761044a611053565b14610454576105fc565b6104fe6001600160a01b038816632ecfe72b61047360408901896111aa565b61047d90806111aa565b610487908061155a565b610495906020810190611571565b6040516001600160e01b031960e084901b16815261ffff9091166004820152602401606060405180830381865afa1580156104d4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104f8919061158c565b85610706565b9092509050600082600481111561051757610517611053565b14610521576105fc565b600061052e8787876107e3565b91945092509050600083600481111561054957610549611053565b1461055457506105fc565b60006105c461056660408901896111aa565b61057090806111aa565b61057a908061155a565b610588906020810190611603565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525086925061099f915050565b9195509350905060008460048111156105df576105df611053565b146105eb5750506105fc565b6105f58582610a05565b9350935050505b9550959350505050565b6000806106168787878787610400565b915091506106248282610a69565b50505050505050565b6000606060006106408460000151610c48565b905060008160405160200161065791815260200190565b6040516020818303038152906040528051906020012090506000866000015190506000610694876040015183858a6020015163ffffffff16610c8d565b905080156106bb5760006040518060200160405280600081525095509550505050506106ff565b6020808801516040805163ffffffff909216928201929092529081018490526060810183905260019060800160405160208183030381529060405295509550505050505b9250929050565b60006060600083602001518460000151610720919061165f565b60ff1690506000856020015163ffffffff16866040015160ff1683620f42406107499190611698565b6107539190611698565b61075f906127106116ac565b61076991906116c3565b865190915060009061077d906127106116e2565b63ffffffff1690508082106107aa57600060405180602001604052806000815250945094505050506106ff565b604080516020810185905290810183905260608101829052600290608001604051602081830303815290604052945094505050506106ff565b6000606081806108006107fb36889003880188611172565b610ca5565b905060006001600160a01b038816636efb46368361082160808b018b611603565b61083160408d0160208e0161170e565b61083e60608e018e61172b565b6040518663ffffffff1660e01b815260040161085e9594939291906118d6565b600060405180830381865afa15801561087b573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526108a39190810190611ab6565b5090506000925060005b6108ba6080890189611603565b905081101561097f57866000015160ff16826020015182815181106108e1576108e1611b52565b60200260200101516108f39190611b68565b6001600160601b031660648360000151838151811061091457610914611b52565b60200260200101516001600160601b031661092f91906116c3565b1061096d5761096a8461094560808b018b611603565b8481811061095557610955611b52565b600192013560f81c9190911b91909117919050565b93505b8061097781611b8e565b9150506108ad565b505060408051602081019091526000808252945092505093509350939050565b6000606060006109ae85610cd0565b90508381168114156109d35760408051602081019091526000808252935091506109fe565b6040805160208101929092528181018590528051808303820181526060909201905260039250905060005b9250925092565b600060606000610a1485610cd0565b9050838116811415610a395750506040805160208101909152600080825291506106ff565b604080516020810183905290810185905260049060600160405160208183030381529060405292509250506106ff565b6000826004811115610a7d57610a7d611053565b1415610a87575050565b6001826004811115610a9b57610a9b611053565b1415610af157600080600083806020019051810190610aba9190611ba9565b60405163d54d727760e01b815260048101849052602481018390526044810182905292955090935091506064015b60405180910390fd5b6002826004811115610b0557610b05611053565b1415610b5957600080600083806020019051810190610b249190611ba9565b6040516001626dc9ad60e11b031981526004810184905260248101839052604481018290529295509093509150606401610ae8565b6003826004811115610b6d57610b6d611053565b1415610bb25760008082806020019051810190610b8a9190611bd7565b604051634a47030360e11b815260048101839052602481018290529193509150604401610ae8565b6004826004811115610bc657610bc6611053565b1415610c0b5760008082806020019051810190610be39190611bd7565b60405163114b085b60e21b815260048101839052602481018290529193509150604401610ae8565b60405162461bcd60e51b8152602060048201526012602482015271556e6b6e6f776e206572726f7220636f646560701b6044820152606401610ae8565b6000610c578260000151610e5d565b6020808401516040808601519051610c70949301611bfb565b604051602081830303815290604052805190602001209050919050565b600083610c9b868585610eaf565b1495945050505050565b600081604051602001610c7091908151815260209182015163ffffffff169181019190915260400190565b600061010082511115610d595760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610ae8565b8151610d6757506000919050565b60008083600081518110610d7d57610d7d611b52565b0160200151600160f89190911c81901b92505b8451811015610e5457848181518110610dab57610dab611b52565b0160200151600160f89190911c1b9150828211610e405760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610ae8565b91811791610e4d81611b8e565b9050610d90565b50909392505050565b6000816000015182602001518360400151604051602001610e8093929190611ca7565b60408051601f198184030181528282528051602091820120606080870151928501919091529183015201610c70565b600060208451610ebf9190611d1c565b15610f465760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a401610ae8565b8260205b85518111610fa957610f5d600285611d1c565b610f7e57816000528086015160205260406000209150600284049350610f97565b8086015160005281602052604060002091506002840493505b610fa2602082611d30565b9050610f4a565b50949350505050565b600060208284031215610fc457600080fd5b81356001600160401b03811115610fda57600080fd5b820160a08185031215610fec57600080fd5b9392505050565b6000815180845260005b8181101561101957602081850181015186830182015201610ffd565b8181111561102b576000602083870101525b50601f01601f19169290920160200192915050565b602081526000610fec6020830184610ff3565b634e487b7160e01b600052602160045260246000fd5b600181811c9082168061107d57607f821691505b6020821081141561024d57634e487b7160e01b600052602260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156110d6576110d661109e565b60405290565b604051606081016001600160401b03811182821017156110d6576110d661109e565b604051608081016001600160401b03811182821017156110d6576110d661109e565b604051601f8201601f191681016001600160401b03811182821017156111485761114861109e565b604052919050565b63ffffffff811681146102dd57600080fd5b803561116d81611150565b919050565b60006040828403121561118457600080fd5b61118c6110b4565b82358152602083013561119e81611150565b60208201529392505050565b60008235605e198336030181126111c057600080fd5b9190910192915050565b803561ffff8116811461116d57600080fd5b600082601f8301126111ed57600080fd5b81356001600160401b038111156112065761120661109e565b611219601f8201601f1916602001611120565b81815284602083860101111561122e57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f83011261125c57600080fd5b6112646110b4565b80604084018581111561127657600080fd5b845b81811015611290578035845260209384019301611278565b509095945050505050565b6000608082840312156112ad57600080fd5b604051604081018181106001600160401b03821117156112cf576112cf61109e565b6040529050806112df848461124b565b81526112ee846040850161124b565b60208201525092915050565b60006001600160401b038211156113135761131361109e565b5060051b60200190565b600082601f83011261132e57600080fd5b8135602061134361133e836112fa565b611120565b82815260059290921b8401810191818101908684111561136257600080fd5b8286015b8481101561138657803561137981611150565b8352918301918301611366565b509695505050505050565b6000606082360312156113a357600080fd5b6113ab6110dc565b82356001600160401b03808211156113c257600080fd5b8185019150606082360312156113d757600080fd5b6113df6110dc565b8235828111156113ee57600080fd5b8301368190036101c081121561140357600080fd5b61140b6110fe565b611414836111ca565b81526020808401358681111561142957600080fd5b611435368287016111dc565b82840152506040603f198401935061016084121561145257600080fd5b61145a6110fe565b8185121561146757600080fd5b61146f6110b4565b9450818601358552606086013583860152848152611490366080880161129b565b838201526114a236610100880161129b565b8282015261018086013594506114b785611150565b8460608201528082850152506101a08501356060840152828652818801359450868511156114e457600080fd5b6114f036868a016111dc565b828701528088013594508685111561150757600080fd5b61151336868a0161131d565b81870152858952611525828c01611162565b828a0152808b013597508688111561153c57600080fd5b61154836898d016111dc565b90890152509598975050505050505050565b600082356101be198336030181126111c057600080fd5b60006020828403121561158357600080fd5b610fec826111ca565b60006060828403121561159e57600080fd5b604051606081018181106001600160401b03821117156115c0576115c061109e565b60405282516115ce81611150565b815260208301516115de81611150565b6020820152604083015160ff811681146115f757600080fd5b60408201529392505050565b6000808335601e1984360301811261161a57600080fd5b8301803591506001600160401b0382111561163457600080fd5b6020019150368190038213156106ff57600080fd5b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff84168082101561167957611679611649565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b6000826116a7576116a7611682565b500490565b6000828210156116be576116be611649565b500390565b60008160001904831182151516156116dd576116dd611649565b500290565b600063ffffffff8083168185168183048111821515161561170557611705611649565b02949350505050565b60006020828403121561172057600080fd5b8135610fec81611150565b6000823561017e198336030181126111c057600080fd5b6000808335601e1984360301811261175957600080fd5b83016020810192503590506001600160401b0381111561177857600080fd5b8060051b36038313156106ff57600080fd5b8183526000602080850194508260005b858110156117c55781356117ad81611150565b63ffffffff168752958201959082019060010161179a565b509495945050505050565b6000808335601e198436030181126117e757600080fd5b83016020810192503590506001600160401b0381111561180657600080fd5b8060061b36038313156106ff57600080fd5b81835260208301925060008160005b8481101561184f57813586526020808301359087015260409586019590910190600101611827565b5093949350505050565b604081833760408201600081526040808301823750600060808301525050565b81835260006020808501808196508560051b810191508460005b878110156118c95782840389526118aa8288611742565b6118b586828461178a565b9a87019a9550505090840190600101611893565b5091979650505050505050565b85815260806020820152836080820152838560a0830137600060a085830101526000601f19601f860116820163ffffffff8516604084015260a08382030160608401526119238485611742565b6101808060a085015261193b6102208501838561178a565b925061194a60208801886117d0565b9250609f19808686030160c0870152611964858584611818565b945061197360408a018a6117d0565b94509150808686030160e087015261198c858584611818565b945061199f610100870160608b01611859565b6119b883870160e08b0180358252602090810135910152565b6119c66101208a018a611742565b9450925080868603016101c08701526119e085858561178a565b94506119f06101408a018a611742565b9450925080868603016101e0870152611a0a85858561178a565b9450611a1a6101608a018a611742565b9450925080868603016102008701525050611a36838383611879565b9b9a5050505050505050505050565b600082601f830112611a5657600080fd5b81516020611a6661133e836112fa565b82815260059290921b84018101918181019086841115611a8557600080fd5b8286015b848110156113865780516001600160601b0381168114611aa95760008081fd5b8352918301918301611a89565b60008060408385031215611ac957600080fd5b82516001600160401b0380821115611ae057600080fd5b9084019060408287031215611af457600080fd5b611afc6110b4565b825182811115611b0b57600080fd5b611b1788828601611a45565b825250602083015182811115611b2c57600080fd5b611b3888828601611a45565b602083015250809450505050602083015190509250929050565b634e487b7160e01b600052603260045260246000fd5b60006001600160601b038083168185168183048111821515161561170557611705611649565b6000600019821415611ba257611ba2611649565b5060010190565b600080600060608486031215611bbe57600080fd5b8351925060208401519150604084015190509250925092565b60008060408385031215611bea57600080fd5b505080516020909101519092909150565b83815260006020606081840152611c156060840186610ff3565b838103604085015284518082528286019183019060005b81811015611c4e57835163ffffffff1683529284019291840191600101611c2c565b509098975050505050505050565b8060005b6002811015611c7f578151845260209384019390910190600101611c60565b50505050565b611c90828251611c5c565b6020810151611ca26040840182611c5c565b505050565b60006101a061ffff86168352806020840152611cc581840186610ff3565b91505082518051604084015260208101516060840152506020830151611cee6080840182611c85565b506040830151611d02610100840182611c85565b5063ffffffff606084015116610180830152949350505050565b600082611d2b57611d2b611682565b500690565b60008219821115611d4357611d43611649565b50019056fea2646970667358221220c84d60f54ed7995780901b9d23b003cbcf860986f7b3a890db5933ce26f29aa464736f6c634300080c0033",
}

// ContractEigenDACertVerifierV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDACertVerifierV2MetaData.ABI instead.
var ContractEigenDACertVerifierV2ABI = ContractEigenDACertVerifierV2MetaData.ABI

// ContractEigenDACertVerifierV2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDACertVerifierV2MetaData.Bin instead.
var ContractEigenDACertVerifierV2Bin = ContractEigenDACertVerifierV2MetaData.Bin

// DeployContractEigenDACertVerifierV2 deploys a new Ethereum contract, binding an instance of ContractEigenDACertVerifierV2 to it.
func DeployContractEigenDACertVerifierV2(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAThresholdRegistryV2 common.Address, _eigenDASignatureVerifierV2 common.Address, _registryCoordinatorV2 common.Address, _securityThresholdsV2 SecurityThresholds, _quorumNumbersRequiredV2 []byte) (common.Address, *types.Transaction, *ContractEigenDACertVerifierV2, error) {
	parsed, err := ContractEigenDACertVerifierV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDACertVerifierV2Bin), backend, _eigenDAThresholdRegistryV2, _eigenDASignatureVerifierV2, _registryCoordinatorV2, _securityThresholdsV2, _quorumNumbersRequiredV2)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDACertVerifierV2{ContractEigenDACertVerifierV2Caller: ContractEigenDACertVerifierV2Caller{contract: contract}, ContractEigenDACertVerifierV2Transactor: ContractEigenDACertVerifierV2Transactor{contract: contract}, ContractEigenDACertVerifierV2Filterer: ContractEigenDACertVerifierV2Filterer{contract: contract}}, nil
}

// ContractEigenDACertVerifierV2 is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV2 struct {
	ContractEigenDACertVerifierV2Caller     // Read-only binding to the contract
	ContractEigenDACertVerifierV2Transactor // Write-only binding to the contract
	ContractEigenDACertVerifierV2Filterer   // Log filterer for contract events
}

// ContractEigenDACertVerifierV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDACertVerifierV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDACertVerifierV2Session struct {
	Contract     *ContractEigenDACertVerifierV2 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDACertVerifierV2CallerSession struct {
	Contract *ContractEigenDACertVerifierV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                        // Call options to use throughout this session
}

// ContractEigenDACertVerifierV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDACertVerifierV2TransactorSession struct {
	Contract     *ContractEigenDACertVerifierV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                        // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV2Raw struct {
	Contract *ContractEigenDACertVerifierV2 // Generic contract binding to access the raw methods on
}

// ContractEigenDACertVerifierV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV2CallerRaw struct {
	Contract *ContractEigenDACertVerifierV2Caller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDACertVerifierV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV2TransactorRaw struct {
	Contract *ContractEigenDACertVerifierV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDACertVerifierV2 creates a new instance of ContractEigenDACertVerifierV2, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV2(address common.Address, backend bind.ContractBackend) (*ContractEigenDACertVerifierV2, error) {
	contract, err := bindContractEigenDACertVerifierV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV2{ContractEigenDACertVerifierV2Caller: ContractEigenDACertVerifierV2Caller{contract: contract}, ContractEigenDACertVerifierV2Transactor: ContractEigenDACertVerifierV2Transactor{contract: contract}, ContractEigenDACertVerifierV2Filterer: ContractEigenDACertVerifierV2Filterer{contract: contract}}, nil
}

// NewContractEigenDACertVerifierV2Caller creates a new read-only instance of ContractEigenDACertVerifierV2, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV2Caller(address common.Address, caller bind.ContractCaller) (*ContractEigenDACertVerifierV2Caller, error) {
	contract, err := bindContractEigenDACertVerifierV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV2Caller{contract: contract}, nil
}

// NewContractEigenDACertVerifierV2Transactor creates a new write-only instance of ContractEigenDACertVerifierV2, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV2Transactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDACertVerifierV2Transactor, error) {
	contract, err := bindContractEigenDACertVerifierV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV2Transactor{contract: contract}, nil
}

// NewContractEigenDACertVerifierV2Filterer creates a new log filterer instance of ContractEigenDACertVerifierV2, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV2Filterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDACertVerifierV2Filterer, error) {
	contract, err := bindContractEigenDACertVerifierV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV2Filterer{contract: contract}, nil
}

// bindContractEigenDACertVerifierV2 binds a generic wrapper to an already deployed contract.
func bindContractEigenDACertVerifierV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDACertVerifierV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierV2.Contract.ContractEigenDACertVerifierV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV2.Contract.ContractEigenDACertVerifierV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV2.Contract.ContractEigenDACertVerifierV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV2.Contract.contract.Transact(opts, method, params...)
}

// CheckDACert is a free data retrieval call binding the contract method 0x5653c730.
//
// Solidity: function checkDACert(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes) cert) view returns(bool)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) CheckDACert(opts *bind.CallOpts, cert EigenDACertV2) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "checkDACert", cert)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckDACert is a free data retrieval call binding the contract method 0x5653c730.
//
// Solidity: function checkDACert(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes) cert) view returns(bool)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) CheckDACert(cert EigenDACertV2) (bool, error) {
	return _ContractEigenDACertVerifierV2.Contract.CheckDACert(&_ContractEigenDACertVerifierV2.CallOpts, cert)
}

// CheckDACert is a free data retrieval call binding the contract method 0x5653c730.
//
// Solidity: function checkDACert(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes) cert) view returns(bool)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) CheckDACert(cert EigenDACertV2) (bool, error) {
	return _ContractEigenDACertVerifierV2.Contract.CheckDACert(&_ContractEigenDACertVerifierV2.CallOpts, cert)
}

// EigenDASignatureVerifierV2 is a free data retrieval call binding the contract method 0x154b9e86.
//
// Solidity: function eigenDASignatureVerifierV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) EigenDASignatureVerifierV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "eigenDASignatureVerifierV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDASignatureVerifierV2 is a free data retrieval call binding the contract method 0x154b9e86.
//
// Solidity: function eigenDASignatureVerifierV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) EigenDASignatureVerifierV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.EigenDASignatureVerifierV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// EigenDASignatureVerifierV2 is a free data retrieval call binding the contract method 0x154b9e86.
//
// Solidity: function eigenDASignatureVerifierV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) EigenDASignatureVerifierV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.EigenDASignatureVerifierV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// EigenDAThresholdRegistryV2 is a free data retrieval call binding the contract method 0x17f3578e.
//
// Solidity: function eigenDAThresholdRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) EigenDAThresholdRegistryV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "eigenDAThresholdRegistryV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistryV2 is a free data retrieval call binding the contract method 0x17f3578e.
//
// Solidity: function eigenDAThresholdRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) EigenDAThresholdRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.EigenDAThresholdRegistryV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// EigenDAThresholdRegistryV2 is a free data retrieval call binding the contract method 0x17f3578e.
//
// Solidity: function eigenDAThresholdRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) EigenDAThresholdRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.EigenDAThresholdRegistryV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// QuorumNumbersRequiredV2 is a free data retrieval call binding the contract method 0xb74d7871.
//
// Solidity: function quorumNumbersRequiredV2() view returns(bytes)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) QuorumNumbersRequiredV2(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "quorumNumbersRequiredV2")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequiredV2 is a free data retrieval call binding the contract method 0xb74d7871.
//
// Solidity: function quorumNumbersRequiredV2() view returns(bytes)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) QuorumNumbersRequiredV2() ([]byte, error) {
	return _ContractEigenDACertVerifierV2.Contract.QuorumNumbersRequiredV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// QuorumNumbersRequiredV2 is a free data retrieval call binding the contract method 0xb74d7871.
//
// Solidity: function quorumNumbersRequiredV2() view returns(bytes)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) QuorumNumbersRequiredV2() ([]byte, error) {
	return _ContractEigenDACertVerifierV2.Contract.QuorumNumbersRequiredV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// RegistryCoordinatorV2 is a free data retrieval call binding the contract method 0x5fafa482.
//
// Solidity: function registryCoordinatorV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) RegistryCoordinatorV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "registryCoordinatorV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinatorV2 is a free data retrieval call binding the contract method 0x5fafa482.
//
// Solidity: function registryCoordinatorV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) RegistryCoordinatorV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.RegistryCoordinatorV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// RegistryCoordinatorV2 is a free data retrieval call binding the contract method 0x5fafa482.
//
// Solidity: function registryCoordinatorV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) RegistryCoordinatorV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.RegistryCoordinatorV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// SecurityThresholdsV2 is a free data retrieval call binding the contract method 0xed0450ae.
//
// Solidity: function securityThresholdsV2() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) SecurityThresholdsV2(opts *bind.CallOpts) (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "securityThresholdsV2")

	outstruct := new(struct {
		ConfirmationThreshold uint8
		AdversaryThreshold    uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfirmationThreshold = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.AdversaryThreshold = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

// SecurityThresholdsV2 is a free data retrieval call binding the contract method 0xed0450ae.
//
// Solidity: function securityThresholdsV2() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) SecurityThresholdsV2() (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	return _ContractEigenDACertVerifierV2.Contract.SecurityThresholdsV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// SecurityThresholdsV2 is a free data retrieval call binding the contract method 0xed0450ae.
//
// Solidity: function securityThresholdsV2() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) SecurityThresholdsV2() (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	return _ContractEigenDACertVerifierV2.Contract.SecurityThresholdsV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x99d923e7.
//
// Solidity: function verifyDACertV2(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes) cert) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) VerifyDACertV2(opts *bind.CallOpts, cert EigenDACertV2) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "verifyDACertV2", cert)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x99d923e7.
//
// Solidity: function verifyDACertV2(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes) cert) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) VerifyDACertV2(cert EigenDACertV2) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertV2(&_ContractEigenDACertVerifierV2.CallOpts, cert)
}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x99d923e7.
//
// Solidity: function verifyDACertV2(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes) cert) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) VerifyDACertV2(cert EigenDACertV2) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertV2(&_ContractEigenDACertVerifierV2.CallOpts, cert)
}
