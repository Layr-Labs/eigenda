package churner

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	bipMultiplier     = big.NewInt(10000)
	secondsTillExpiry = 90 * time.Second
	zeroAddressString = "0x0000000000000000000000000000000000000000"
)

type ChurnRequest struct {
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
	Indexer              thegraph.IndexedChainState
	Transactor           core.Transactor
	StakeRegistryAddress gethcommon.Address
	privateKey           *ecdsa.PrivateKey
	logger               common.Logger
	metrics              *Metrics
}

func NewChurner(
	config *Config,
	indexer thegraph.IndexedChainState,
	transactor core.Transactor,
	logger common.Logger,
	metrics *Metrics,
) (*churner, error) {
	stakeRegistryAddress, err := transactor.StakeRegistry(context.Background())
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(config.EthClientConfig.PrivateKeyString)
	if err != nil {
		return nil, err
	}

	return &churner{
		Indexer:              indexer,
		Transactor:           transactor,
		StakeRegistryAddress: stakeRegistryAddress,
		privateKey:           privateKey,
		logger:               logger,
		metrics:              metrics,
	}, nil
}

func (c *churner) VerifyRequestSignature(ctx context.Context, churnRequest *ChurnRequest) (gethcommon.Address, error) {
	operatorToRegisterAddress, err := c.Transactor.OperatorIDToAddress(ctx, churnRequest.OperatorToRegisterPubkeyG1.GetOperatorID())
	if err != nil {
		return gethcommon.Address{}, err
	}
	if operatorToRegisterAddress == gethcommon.HexToAddress(zeroAddressString) {
		return gethcommon.Address{}, errors.New("operatorToRegisterPubkey is not registered with bls pubkey compendium")
	}

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
				return nil, errors.New("operator is already registered in quorum")
			}
		}
	}

	return c.createChurnResponse(ctx, operatorToRegisterId, operatorToRegisterAddress, churnRequest.QuorumIDs)
}

func (c *churner) createChurnResponse(
	ctx context.Context,
	operatorToRegisterId core.OperatorID,
	operatorToRegisterAddress gethcommon.Address,
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

	signatureWithSaltAndExpiry, err := c.sign(ctx, operatorToRegisterId, operatorsToChurn)
	if err != nil {
		return nil, err
	}
	return &ChurnResponse{
		SignatureWithSaltAndExpiry: signatureWithSaltAndExpiry,
		OperatorsToChurn:           operatorsToChurn,
	}, nil
}

func (c *churner) getOperatorsToChurn(ctx context.Context, quorumIDs []uint8, operatorStakes [][]core.OperatorStake, operatorToRegisterAddress gethcommon.Address, currentBlockNumber uint32) ([]core.OperatorToChurn, error) {
	operatorsToChurn := make([]core.OperatorToChurn, 0)
	for i, quorumID := range quorumIDs {
		operatorSetParams, err := c.Transactor.GetOperatorSetParams(ctx, quorumID)
		if err != nil {
			return nil, nil
		}

		if operatorSetParams.MaxOperatorCount != uint32(len(operatorStakes[i])) {
			// quorum is not full, so we can continue
			continue
		}

		operatorToRegisterStake, err := c.Transactor.WeightOfOperatorForQuorum(ctx, quorumID, operatorToRegisterAddress)
		if err != nil {
			return nil, nil
		}

		// loop through operator stakes for the quorum and find the lowest one
		totalStake := big.NewInt(0)
		lowestStakeOperatorId := operatorStakes[i][0].OperatorID
		lowestStake := operatorStakes[i][0].Stake
		for _, operatorStake := range operatorStakes[i] {
			if operatorStake.Stake.Cmp(lowestStake) < 0 {
				lowestStake = operatorStake.Stake
				lowestStakeOperatorId = operatorStake.OperatorID
			}
			totalStake.Add(totalStake, operatorStake.Stake)
		}

		churnBIPsOfOperatorStake := big.NewInt(int64(operatorSetParams.ChurnBIPsOfOperatorStake))
		churnBIPsOfTotalStake := big.NewInt(int64(operatorSetParams.ChurnBIPsOfTotalStake))

		c.logger.Info("lowestStake", "lowestStake", lowestStake.String(), "operatorToRegisterStake", operatorToRegisterStake.String(), "totalStake", totalStake.String())

		// verify the lowest stake against the registering operator's stake
		// make sure that: lowestStake * churnBIPsOfOperatorStake < operatorToRegisterStake * bipMultiplier
		if new(big.Int).Mul(lowestStake, churnBIPsOfOperatorStake).Cmp(new(big.Int).Mul(operatorToRegisterStake, bipMultiplier)) >= 0 {
			c.metrics.IncrementFailedRequestNum("getOperatorsToChurn", "Insufficient stake: operator doesn't have enough stake")
			return nil, errors.New("registering operator has less than churnBIPsOfOperatorStake")
		}

		// verify the lowest stake against the total stake
		// make sure that: lowestStake * bipMultiplier < totalStake * churnBIPsOfTotalStake
		if new(big.Int).Mul(lowestStake, bipMultiplier).Cmp(new(big.Int).Mul(totalStake, churnBIPsOfTotalStake)) >= 0 {
			return nil, errors.New("operator to churn has less than churnBIPSOfTotalStake")
		}

		operatorToChurnAddress, err := c.Transactor.OperatorIDToAddress(ctx, lowestStakeOperatorId)
		if err != nil {
			return nil, err
		}

		operatorToChurnIndexedInfo, err := c.Indexer.GetIndexedOperatorInfoByOperatorId(ctx, lowestStakeOperatorId, currentBlockNumber)
		if err != nil {
			return nil, err
		}

		// add the operator to churn to the list
		operatorsToChurn = append(operatorsToChurn, core.OperatorToChurn{
			QuorumId: quorumIDs[i],
			Operator: operatorToChurnAddress,
			Pubkey:   operatorToChurnIndexedInfo.PubkeyG1,
		})
	}
	return operatorsToChurn, nil
}

func (c *churner) sign(ctx context.Context, operatorToRegisterId core.OperatorID, operatorsToChurn []core.OperatorToChurn) (*SignatureWithSaltAndExpiry, error) {
	now := time.Now()
	privateKeyBytes := crypto.FromECDSA(c.privateKey)
	saltKeccak256 := crypto.Keccak256([]byte("churn"), []byte(now.String()), operatorToRegisterId[:], privateKeyBytes)

	var salt [32]byte
	copy(salt[:], saltKeccak256)

	// set expiry to 90s in the future
	expiry := big.NewInt(now.Add(secondsTillExpiry).Unix())

	// sign and return signature
	hashToSign, err := c.Transactor.CalculateOperatorChurnApprovalDigestHash(ctx, operatorToRegisterId, operatorsToChurn, salt, expiry)
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
		churnRequest.OperatorToRegisterPubkeyG1.Serialize(),
		churnRequest.OperatorToRegisterPubkeyG2.Serialize(),
		churnRequest.Salt[:],
	)
	copy(requestHash[:], requestHashBytes)
	return requestHash
}
