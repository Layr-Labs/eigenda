package mock

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	_ "embed"

	"github.com/Layr-Labs/eigenda/indexer/test/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
)

//go:embed chain.json
var mockChainJson string

const gasLimit = 10000000

type (
	ContractSimulator struct {
		Client       SimulatedBackend
		WethAddr     common.Address
		DeployerPK   *ecdsa.PrivateKey
		DeployerAddr common.Address
	}

	MockChain struct {
		Chain []struct {
			Id   int  `json:"id"`
			Fork *int `json:"fork"`
		} `json:"chain"`
	}
)

func MustNewContractSimulator() *ContractSimulator {
	sb, deployerAddr, deployerPK := mustNewSimulatedBackend()
	wethAddress, err := mustDeployWethContract(sb, deployerPK)
	if err != nil {
		log.Fatal(err)
	}

	return &ContractSimulator{
		Client:       sb,
		WethAddr:     wethAddress,
		DeployerPK:   deployerPK,
		DeployerAddr: deployerAddr,
	}
}

func (cs *ContractSimulator) Start(blockWait time.Duration, cancel context.CancelFunc) {
	mockChain, err := parseChainJson()
	if err != nil {
		log.Fatal(err)
	}

	hashById := make(map[int]common.Hash)

	wethInstance, err := contracts.NewWeth(cs.WethAddr, cs.Client)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for _, c := range mockChain.Chain {
			if c.Fork != nil {
				fmt.Println("Forking to hash: ", hashById[*c.Fork])
				err = cs.Client.Fork(hashById[*c.Fork])
				if err != nil {
					log.Fatal(err)
				}
			}

			auth, err := GenerateTransactOpts(cs.Client, cs.DeployerPK)
			if err != nil {
				log.Fatal(err)
			}

			auth.Value = big.NewInt(int64(c.Id + 1))
			_, err = wethInstance.Deposit(auth)
			if err != nil {
				log.Fatal(err)
			}

			hash := cs.Client.Commit()
			hashById[c.Id] = hash
			if blockWait > 0 {
				time.Sleep(blockWait)
			}
		}
		// Sleep for a second to give indexer time to finish indexing the events before cancelling the context
		time.Sleep(1 * time.Second)
		cancel()
	}()
}

func (cs *ContractSimulator) DepositEvents() ([]*contracts.WethDeposit, error) {
	opts := &bind.FilterOpts{
		Start: 0,
		End:   nil,
	}

	wethInstance, err := contracts.NewWeth(cs.WethAddr, cs.Client)
	if err != nil {
		return nil, err
	}

	events, err := wethInstance.FilterDeposit(opts, []common.Address{})
	if err != nil {
		return nil, err
	}

	depositEvents := make([]*contracts.WethDeposit, 0, 5)
	for events.Next() {
		depositEvents = append(depositEvents, events.Event)
	}

	return depositEvents, nil
}

func mustNewSimulatedBackend() (client SimulatedBackend, deployerAddr common.Address, privateKey *ecdsa.PrivateKey) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}

	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth in wei

	deployerAddr = auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		deployerAddr: {
			Balance: balance,
		},
	}

	blockGasLimit := uint64(gasLimit)
	b := simulated.NewBackend(genesisAlloc, simulated.WithBlockGasLimit(blockGasLimit))
	client = &simulatedBackend{
		Backend: b,
		Client:  b.Client(),
	}
	return
}

func mustDeployWethContract(client SimulatedBackend, privateKey *ecdsa.PrivateKey) (address common.Address, err error) {
	auth, err := GenerateTransactOpts(client, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	address, tx, _, err := contracts.DeployWeth(auth, client)
	if err != nil {
		log.Fatal(err)
	}
	client.Commit()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = bind.WaitDeployed(ctx, client, tx)
	if err != nil {
		log.Fatal("Error deploying smart contract: ", err)
	}
	return
}

func GenerateTransactOpts(client SimulatedBackend, privateKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		return nil, err
	}

	fromAddress := auth.From
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)                         // in wei 1 eth
	auth.GasLimit = uint64(gasLimit)                   // in units
	auth.GasPrice = new(big.Int).SetUint64(1000000000) // Set gas price to 1000000000 wei when using the simulated backend.

	return auth, nil
}

func parseChainJson() (MockChain, error) {
	var data MockChain
	err := json.Unmarshal([]byte(mockChainJson), &data)
	if err != nil {
		return MockChain{}, err
	}
	return data, nil
}
