package churner

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	bipMultiplier = big.NewInt(10000)
)

type ChurnRequest struct {
	OperatorAddress            gethcommon.Address
	OperatorToRegisterPubkeyG1 *core.G1Point
	OperatorToRegisterPubkeyG2 *core.G2Point
	OperatorRequestSignature   *core.Signature
	Salt                       [32]byte
	QuorumIDs                  []core.QuorumID
}

type SignatureWithSaltAndExpiry struct {
	Signature []byte
	Salt      [32]byte
	Expiry    *big.Int
}

type ChurnResponse struct {
	SignatureWithSaltAndExpiry *SignatureWithSaltAndExpiry
	OperatorsToChurn           []core.OperatorToChurn
}

type churner struct {
	mu          sync.Mutex
	Indexer     thegraph.IndexedChainState
	Transactor  core.Writer
	QuorumCount uint8

	privateKey            *ecdsa.PrivateKey
	logger                logging.Logger
	metrics               *Metrics
	churnApprovalInterval time.Duration
}

func NewChurner(
	config *Config,
	indexer thegraph.IndexedChainState,
	transactor core.Writer,
	logger logging.Logger,
	metrics *Metrics,
) (*churner, error) {
	privateKey, err := crypto.HexToECDSA(config.EthClientConfig.PrivateKeyString)
	if err != nil {
		return nil, err
	}

	logger.Info("Churner created with config", "ChurnApprovalInterval", config.ChurnApprovalInterval)

	return &churner{
		Indexer:     indexer,
		Transactor:  transactor,
		QuorumCount: 0,

		privateKey:            privateKey,
		logger:                logger.With("component", "Churner"),
		metrics:               metrics,
		churnApprovalInterval: config.ChurnApprovalInterval,
	}, nil
}

func (c *churner) VerifyRequestSignature(ctx context.Context, churnRequest *ChurnRequest) (gethcommon.Address, error) {
	operatorToRegisterAddress := churnRequest.OperatorAddress
	isEqual, err := churnRequest.OperatorToRegisterPubkeyG1.VerifyEquivalence(churnRequest.OperatorToRegisterPubkeyG2)
	if err != nil {
		return gethcommon.Address{}, err
	}
	if !isEqual {
		return gethcommon.Address{}, errors.New("operatorToRegisterPubkeyG1 and operatorToRegisterPubkeyG2 are not equivalent")
	}

	requestHash := CalculateRequestHash(churnRequest)
	ok := churnRequest.OperatorRequestSignature.Verify(churnRequest.OperatorToRegisterPubkeyG2, requestHash)
	if !ok {
		return gethcommon.Address{}, errors.New("operatorRequestSignature is invalid")
	}
	return operatorToRegisterAddress, nil
}

func (c *churner) ProcessChurnRequest(ctx context.Context, operatorToRegisterAddress gethcommon.Address, churnRequest *ChurnRequest) (*ChurnResponse, error) {
	operatorToRegisterId := churnRequest.OperatorToRegisterPubkeyG1.GetOperatorID()

	quorumBitmap, err := c.Transactor.GetCurrentQuorumBitmapByOperatorId(ctx, operatorToRegisterId)
	if err != nil {
		return nil, err
	}

	quorumIDsAlreadyRegisteredFor := eth.BitmapToQuorumIds(quorumBitmap)

	// check if the operator is already registered in the quorums
	for _, quorumID := range churnRequest.QuorumIDs {
		for _, quorumIDAlreadyRegisteredFor := range quorumIDsAlreadyRegisteredFor {
			if quorumIDAlreadyRegisteredFor == quorumID {
				return nil, api.NewErrorInvalidArg("operator is already registered in quorum")
			}
		}
	}

	return c.createChurnResponse(ctx, operatorToRegisterAddress, operatorToRegisterId, churnRequest.QuorumIDs)
}

func (c *churner) UpdateQuorumCount(ctx context.Context) error {
	currentBlock, err := c.Transactor.GetCurrentBlockNumber(ctx)
	if err != nil {
		return err
	}
	count, err := c.Transactor.GetQuorumCount(ctx, currentBlock)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.QuorumCount = count
	c.mu.Unlock()
	return nil
}

func (c *churner) createChurnResponse(
	ctx context.Context,
	operatorToRegisterAddress gethcommon.Address,
	operatorToRegisterId core.OperatorID,
	quorumIDs []core.QuorumID,
) (*ChurnResponse, error) {
	currentBlockNumber, err := c.Transactor.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	// get the operator list for each quorum
	operatorStakes, err := c.Transactor.GetOperatorStakesForQuorums(ctx, quorumIDs, currentBlockNumber)
	if err != nil {
		return nil, err
	}

	// get the registering operator's stakes for each quorum
	operatorsToChurn, err := c.getOperatorsToChurn(ctx, quorumIDs, operatorStakes, operatorToRegisterAddress, currentBlockNumber)
	if err != nil {
		return nil, err
	}

	signatureWithSaltAndExpiry, err := c.sign(ctx, operatorToRegisterAddress, operatorToRegisterId, operatorsToChurn)
	if err != nil {
		return nil, err
	}
	return &ChurnResponse{
		SignatureWithSaltAndExpiry: signatureWithSaltAndExpiry,
		OperatorsToChurn:           operatorsToChurn,
	}, nil
}

func (c *churner) getOperatorsToChurn(ctx context.Context, quorumIDs []uint8, operatorStakes core.OperatorStakes, operatorToRegisterAddress gethcommon.Address, currentBlockNumber uint32) ([]core.OperatorToChurn, error) {
	operatorsToChurn := make([]core.OperatorToChurn, 0)
	for i, quorumID := range quorumIDs {
		operatorSetParams, err := c.Transactor.GetOperatorSetParams(ctx, quorumID)
		if err != nil {
			return nil, nil
		}

		if operatorSetParams.MaxOperatorCount == 0 {
			return nil, errors.New("maxOperatorCount is 0")
		}

		if uint32(len(operatorStakes[quorumID])) < operatorSetParams.MaxOperatorCount {
			// quorum is not full, so we leave out the operator for the quorum
			c.logger.Info("quorum is not full", "quorumID", quorumID, "maxOperatorCount", operatorSetParams.MaxOperatorCount, "numOperators", len(operatorStakes[quorumID]))
			operatorsToChurn = append(operatorsToChurn, core.OperatorToChurn{
				QuorumId: quorumIDs[i],
				Operator: gethcommon.Address{0},
				Pubkey:   nil,
			})
			continue
		}
		if len(operatorStakes[quorumID]) == 0 {
			c.logger.Info("no operators in quorum", "quorumID", quorumID)
			operatorsToChurn = append(operatorsToChurn, core.OperatorToChurn{
				QuorumId: quorumIDs[i],
				Operator: gethcommon.Address{0},
				Pubkey:   nil,
			})
			continue
		}

		operatorToRegisterStake, err := c.Transactor.WeightOfOperatorForQuorum(ctx, quorumID, operatorToRegisterAddress)
		if err != nil {
			return nil, nil
		}

		// loop through operator stakes for the quorum and find the lowest one
		totalStake := big.NewInt(0)
		lowestStakeOperatorId := operatorStakes[quorumID][0].OperatorID
		lowestStake := operatorStakes[quorumID][0].Stake
		for _, operatorStake := range operatorStakes[quorumID] {
			if operatorStake.Stake.Cmp(lowestStake) < 0 {
				lowestStake = operatorStake.Stake
				lowestStakeOperatorId = operatorStake.OperatorID
			}
			totalStake.Add(totalStake, operatorStake.Stake)
		}

		churnBIPsOfOperatorStake := big.NewInt(int64(operatorSetParams.ChurnBIPsOfOperatorStake))
		churnBIPsOfTotalStake := big.NewInt(int64(operatorSetParams.ChurnBIPsOfTotalStake))

		c.logger.Info("lowestStake", "lowestStake", lowestStake.String(), "operatorToRegisterStake", operatorToRegisterStake.String(), "totalStake", totalStake.String(), "operatorToRegisterAddress", operatorToRegisterAddress.Hex(), "lowestStakeOperatorId", lowestStakeOperatorId.Hex())

		// verify the lowest stake against the registering operator's stake
		// make sure that: lowestStake * churnBIPsOfOperatorStake < operatorToRegisterStake * bipMultiplier
		// This means the registering operator needs to have greater than
		// churnBIPsOfOperatorStake/10000 times the stake of lowest stake in order to
		// churn the lowest-stake operator out.
		// For example, when churnBIPsOfOperatorStake=11000, the operator trying to
		// register needs to have 1.1 times the stake of the lowest-stake operator.
		if new(big.Int).Mul(lowestStake, churnBIPsOfOperatorStake).Cmp(new(big.Int).Mul(operatorToRegisterStake, bipMultiplier)) >= 0 {
			c.metrics.IncrementFailedRequestNum("getOperatorsToChurn", FailReasonInsufficientStakeToRegister)
			msg := "registering operator must have %f%% more than the stake of the " +
				"lowest-stake operator. Block number used for this decision: %d, " +
				"registering operator address: %s, registering operator stake: %d, " +
				"stake of lowest-stake operator: %d, operatorId of lowest-stake operator: " +
				"%x, quorum ID: %d"
			return nil, api.NewErrorInvalidArg(fmt.Sprintf(msg, float64(operatorSetParams.ChurnBIPsOfOperatorStake)/100.0-100.0, currentBlockNumber, operatorToRegisterAddress.Hex(), operatorToRegisterStake, lowestStake, lowestStakeOperatorId, quorumID))
		}

		// verify the lowest stake against the total stake
		// make sure that: lowestStake * bipMultiplier < totalStake * churnBIPsOfTotalStake
		// For the lowest-stake operator to be churned out, it must have less than
		// churnBIPsOfTotalStake/10000 of the total stake.
		// For example, when churnBIPsOfTotalStake=1001, the operator to be churned out
		// (i.e. the lowest-stake operator) needs to have less than 10.01% of the total
		// stake.
		if new(big.Int).Mul(lowestStake, bipMultiplier).Cmp(new(big.Int).Mul(totalStake, churnBIPsOfTotalStake)) >= 0 {
			c.metrics.IncrementFailedRequestNum("getOperatorsToChurn", FailReasonInsufficientStakeToChurn)
			msg := "operator to churn out must have less than %f%% of the total stake. " +
				"Block number used for this decision: %d, operatorId of the operator " +
				"to churn: %x, stake of the operator to churn: %d, total stake in " +
				"quorum: %d, quorum ID: %d"
			return nil, api.NewErrorInvalidArg(fmt.Sprintf(msg, float64(operatorSetParams.ChurnBIPsOfTotalStake)/100.0, currentBlockNumber, lowestStakeOperatorId.Hex(), lowestStake, totalStake, quorumID))
		}

		operatorToChurnAddress, err := c.Transactor.OperatorIDToAddress(ctx, lowestStakeOperatorId)
		if err != nil {
			return nil, err
		}

		operatorToChurnIndexedInfo, err := c.Indexer.GetIndexedOperatorInfoByOperatorId(ctx, lowestStakeOperatorId, currentBlockNumber)
		if err != nil {
			return nil, err
		}

		// log the churn decision just made
		c.logger.Info("Churner made a churn decision", "address of operator churned out", operatorToChurnAddress.Hex(), "stake of operator churned out", lowestStake.String(), "address of operator churned in", operatorToRegisterAddress.Hex(), "stake of operator churned in", operatorToRegisterStake.String(), "block number", currentBlockNumber, "quorumID", quorumID)

		// add the operator to churn to the list
		operatorsToChurn = append(operatorsToChurn, core.OperatorToChurn{
			QuorumId: quorumIDs[i],
			Operator: operatorToChurnAddress,
			Pubkey:   operatorToChurnIndexedInfo.PubkeyG1,
		})
	}
	return operatorsToChurn, nil
}

func (c *churner) sign(ctx context.Context, operatorToRegisterAddress gethcommon.Address, operatorToRegisterId core.OperatorID, operatorsToChurn []core.OperatorToChurn) (*SignatureWithSaltAndExpiry, error) {
	now := time.Now()
	privateKeyBytes := crypto.FromECDSA(c.privateKey)
	saltKeccak256 := crypto.Keccak256([]byte("churn"), []byte(now.String()), operatorToRegisterId[:], privateKeyBytes)

	var salt [32]byte
	copy(salt[:], saltKeccak256)

	// set expiry to ChurnApprovalInterval in the future
	expiry := big.NewInt(now.Add(c.churnApprovalInterval).Unix())

	// sign and return signature
	hashToSign, err := c.Transactor.CalculateOperatorChurnApprovalDigestHash(ctx, operatorToRegisterAddress, operatorToRegisterId, operatorsToChurn, salt, expiry)
	if err != nil {
		return nil, err
	}
	signature, err := crypto.Sign(hashToSign[:], c.privateKey)
	if err != nil {
		return nil, err
	}
	if signature[64] != 27 && signature[64] != 28 {
		signature[64] += 27
	}
	return &SignatureWithSaltAndExpiry{
		Signature: signature,
		Salt:      salt,
		Expiry:    expiry,
	}, nil
}

func CalculateRequestHash(churnRequest *ChurnRequest) [32]byte {
	var requestHash [32]byte
	requestHashBytes := crypto.Keccak256(
		[]byte("ChurnRequest"),
		[]byte(churnRequest.OperatorAddress.Hex()),
		churnRequest.OperatorToRegisterPubkeyG1.Serialize(),
		churnRequest.OperatorToRegisterPubkeyG2.Serialize(),
		churnRequest.Salt[:],
	)
	copy(requestHash[:], requestHashBytes)
	return requestHash
}
