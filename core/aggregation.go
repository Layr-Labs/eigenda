package core

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sort"

	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	lru "github.com/hashicorp/golang-lru/v2"
)

const maxNumOperatorAddresses = 300

var (
	ErrPubKeysNotEqual     = errors.New("public keys are not equal")
	ErrInsufficientEthSigs = errors.New("insufficient eth signatures")
	ErrAggPubKeyNotValid   = errors.New("aggregated public key is not valid")
	ErrAggSigNotValid      = errors.New("aggregated signature is not valid")
)

type SigningMessage struct {
	Signature       *Signature
	Operator        OperatorID
	BatchHeaderHash [32]byte
	// Undefined if this value <= 0.
	AttestationLatencyMs float64
	Err                  error
}

// QuorumAttestation contains the results of aggregating signatures from a set of operators by quorums
// It also returns map of all signers across all quorums
type QuorumAttestation struct {
	// QuorumAggPubKeys contains the aggregated public keys for all of the operators each quorum,
	// including those that did not sign
	QuorumAggPubKey map[QuorumID]*G1Point
	// SignersAggPubKey is the aggregated public key for all of the operators that signed the message by each quorum
	SignersAggPubKey map[QuorumID]*G2Point
	// AggSignature is the aggregated signature for all of the operators that signed the message for each quorum,
	// mirroring the SignersAggPubKey.
	AggSignature map[QuorumID]*Signature
	// QuorumResults contains the quorum ID and the amount signed for each quorum
	QuorumResults map[QuorumID]*QuorumResult
	// SignerMap contains the operator IDs that signed the message
	SignerMap map[OperatorID]bool
}

// SignatureAggregation contains the results of aggregating signatures from a set of operators across multiple quorums
type SignatureAggregation struct {
	// NonSigners contains the public keys of the operators that did not sign the message
	NonSigners []*G1Point
	// QuorumAggPubKeys contains the aggregated public keys for all of the operators each quorum,
	// Including those that did not sign
	QuorumAggPubKeys map[QuorumID]*G1Point
	// AggPubKey is the aggregated public key for all of the operators that signed the message,
	// further aggregated across the quorums; operators signing for multiple quorums will be included in
	// the aggregation multiple times
	AggPubKey *G2Point
	// AggSignature is the aggregated signature for all of the operators that signed the message, mirroring the
	// AggPubKey.
	AggSignature *Signature
	// QuorumResults contains the quorum ID and the amount signed for each quorum
	QuorumResults map[QuorumID]*QuorumResult
}

// SignatureAggregator is an interface for aggregating the signatures returned by DA nodes so that they can be
// verified by the DA contract
type SignatureAggregator interface {
	// ReceiveSignatures blocks until it receives a response for each operator in the operator state via messageChan,
	// and then returns the attestation result by quorum.
	ReceiveSignatures(
		ctx context.Context,
		state *IndexedOperatorState,
		message [32]byte,
		messageChan chan SigningMessage,
	) (*QuorumAttestation, error)

	// AggregateSignatures takes attestation result by quorum and aggregates the signatures across them.
	// If the aggregated signature is invalid, an error is returned.
	AggregateSignatures(
		ctx context.Context,
		ics IndexedChainState,
		referenceBlockNumber uint,
		quorumAttestation *QuorumAttestation,
		quorumIDs []QuorumID,
	) (*SignatureAggregation, error)
}

type StdSignatureAggregator struct {
	Logger     logging.Logger
	Transactor Reader
	// OperatorAddresses contains the ethereum addresses of the operators corresponding to their operator IDs
	OperatorAddresses *lru.Cache[OperatorID, gethcommon.Address]
}

func NewStdSignatureAggregator(logger logging.Logger, transactor Reader) (*StdSignatureAggregator, error) {
	operatorAddrs, err := lru.New[OperatorID, gethcommon.Address](maxNumOperatorAddresses)
	if err != nil {
		return nil, err
	}

	return &StdSignatureAggregator{
		Logger:            logger.With("component", "SignatureAggregator"),
		Transactor:        transactor,
		OperatorAddresses: operatorAddrs,
	}, nil
}

var _ SignatureAggregator = (*StdSignatureAggregator)(nil)

// getQuorumIDs returns a sorted list quorum IDs from the state.
func (a *StdSignatureAggregator) getQuorumIDs(state *IndexedOperatorState) ([]QuorumID, error) {
	quorumIDs := make([]QuorumID, 0, len(state.AggKeys))
	for quorumID := range state.Operators {
		quorumIDs = append(quorumIDs, quorumID)
	}
	slices.Sort(quorumIDs)

	if len(quorumIDs) == 0 {
		return nil, errors.New("the number of quorums must be greater than zero")
	}

	// Ensure all quorums are found in state
	for _, id := range quorumIDs {
		_, found := state.Operators[id]
		if !found {
			return nil, fmt.Errorf("quorum %d not found", id)
		}
	}

	return quorumIDs, nil
}

// getOperatorAddress returns the ethereum address of the operator corresponding to the operator ID.
func (a *StdSignatureAggregator) getOperatorAddress(
	ctx context.Context,
	operatorIDHex string,
	signingMessage SigningMessage) gethcommon.Address {

	operatorAddr, ok := a.OperatorAddresses.Get(signingMessage.Operator)

	if !ok && a.Transactor != nil {
		var err error
		operatorAddr, err = a.Transactor.OperatorIDToAddress(ctx, signingMessage.Operator)
		if err != nil {
			a.Logger.Warn("failed to get operator address from registry", "operatorID", operatorIDHex)
			operatorAddr = gethcommon.Address{}
		} else {
			a.OperatorAddresses.Add(signingMessage.Operator, operatorAddr)
		}
	} else if !ok {
		operatorAddr = gethcommon.Address{}
	}

	return operatorAddr
}

// collectSignatureInfoForQuorum collects the signature information for the operator in the quorum, updating the
// attestation and stake signed.
func (a *StdSignatureAggregator) collectSignatureInfoForQuorum(
	state *IndexedOperatorState,
	operator *IndexedOperatorInfo,
	quorumID QuorumID,
	signingMessage SigningMessage,
	attestation *QuorumAttestation,
	stakeSigned map[QuorumID]*big.Int,
	sig *Signature,
	operatorQuorums []uint8) {

	// Get stake amounts for operator
	ops := state.Operators[quorumID]
	opInfo, ok := ops[signingMessage.Operator]
	// If operator is not in quorum, skip
	if !ok {
		return
	}
	operatorQuorums = append(operatorQuorums, quorumID)

	attestation.SignerMap[signingMessage.Operator] = true

	// Add to stake signed
	stakeSigned[quorumID].Add(stakeSigned[quorumID], opInfo.Stake)

	// Add to agg signature
	if attestation.AggSignature[quorumID] == nil {
		attestation.AggSignature[quorumID] = &Signature{sig.Clone()}
		attestation.SignersAggPubKey[quorumID] = operator.PubkeyG2.Clone()
	} else {
		attestation.AggSignature[quorumID].Add(sig.G1Point)
		attestation.SignersAggPubKey[quorumID].Add(operator.PubkeyG2)
	}
}

// processNextSignature is used to collect the next signature from the message channel and aggregate it into the
// attestation. It blocks until a signature is received. Once received, it verifies the signature and adds it to the
// attestation.
func (a *StdSignatureAggregator) processNextSignature(
	ctx context.Context,
	state *IndexedOperatorState,
	quorumIDs []QuorumID,
	message [32]byte,
	stakeSigned map[QuorumID]*big.Int,
	attestation *QuorumAttestation,
	messageChan chan SigningMessage) {

	signingMessage := <-messageChan
	if seen := attestation.SignerMap[signingMessage.Operator]; seen {
		a.Logger.Warn("duplicate signature received", "operatorID", signingMessage.Operator.Hex())
		return
	}

	operatorIDHex := signingMessage.Operator.Hex()
	operatorAddr := a.getOperatorAddress(ctx, operatorIDHex, signingMessage)
	socket := ""
	if op, ok := state.IndexedOperators[signingMessage.Operator]; ok {
		socket = op.Socket
	}
	batchHeaderHashHex := hex.EncodeToString(signingMessage.BatchHeaderHash[:])

	if signingMessage.Err != nil {
		a.Logger.Warn("error returned from messageChan",
			"operatorID", operatorIDHex,
			"operatorAddress", operatorAddr,
			"socket", socket,
			"batchHeaderHash", batchHeaderHashHex,
			"attestationLatencyMs", signingMessage.AttestationLatencyMs,
			"err", signingMessage.Err)
		return
	}

	operator, found := state.IndexedOperators[signingMessage.Operator]
	if !found {
		a.Logger.Error("Operator not found in state",
			"operatorID", operatorIDHex,
			"operatorAddress", operatorAddr,
			"socket", socket)
		return
	}

	// Verify Signature
	sig := signingMessage.Signature
	ok := sig.Verify(operator.PubkeyG2, message)
	if !ok {
		a.Logger.Error("signature is not valid",
			"operatorID", operatorIDHex,
			"operatorAddress", operatorAddr,
			"socket", socket,
			"pubkey", hexutil.Encode(operator.PubkeyG2.Serialize()))
		return
	}

	// Collect signature information for eqch quorum
	operatorQuorums := make([]uint8, 0, len(quorumIDs))
	for _, quorumID := range quorumIDs {
		a.collectSignatureInfoForQuorum(
			state, operator, quorumID, signingMessage, attestation, stakeSigned, sig, operatorQuorums)
	}
	a.Logger.Info("received signature from operator",
		"operatorID", operatorIDHex,
		"operatorAddress", operatorAddr,
		"socket", socket,
		"quorumIDs", fmt.Sprint(operatorQuorums),
		"batchHeaderHash", batchHeaderHashHex,
		"attestationLatencyMs", signingMessage.AttestationLatencyMs)
}

// aggregateNonSigners aggregates the public keys of the operators that did not sign the message. It returns
// a list of non-signer public keys and operator IDs.
func (a *StdSignatureAggregator) aggregateNonSigners(
	state *IndexedOperatorState,
	attestation *QuorumAttestation) (nonSignerKeys []*G1Point, nonSignerOperatorIds []OperatorID) {

	nonSignerKeys = make([]*G1Point, 0)
	nonSignerOperatorIds = make([]OperatorID, 0)

	for id, op := range state.IndexedOperators {
		_, found := attestation.SignerMap[id]
		if !found {
			nonSignerKeys = append(nonSignerKeys, op.PubkeyG1)
			nonSignerOperatorIds = append(nonSignerOperatorIds, id)
		}
	}

	return nonSignerKeys, nonSignerOperatorIds
}

// processQuorumStatus processes the status of a quorum, verifying that the aggregated public key and
// signature are valid and computing the signing percentage. The attestation is updated with the results.
func (a *StdSignatureAggregator) processQuorumStatus(
	state *IndexedOperatorState,
	quorumID QuorumID,
	stakeSigned map[QuorumID]*big.Int,
	attestation *QuorumAttestation,
	message [32]byte,
	nonSignerKeys []*G1Point,
	nonSignerOperatorIds []OperatorID) error {

	// Check that quorum has sufficient stake
	percent := GetSignedPercentage(state.OperatorState, quorumID, stakeSigned[quorumID])
	attestation.QuorumResults[quorumID] = &QuorumResult{
		QuorumID:      quorumID,
		PercentSigned: percent,
	}

	if percent == 0 {
		a.Logger.Warn("no stake signed for quorum", "quorumID", quorumID)
		return nil
	}

	// Verify that the aggregated public key for the quorum matches the on-chain
	// quorum aggregate public key sans non-signers of the quorum
	quorumAggKey := state.AggKeys[quorumID]
	attestation.QuorumAggPubKey[quorumID] = quorumAggKey

	signersAggKey := quorumAggKey.Clone()
	for opInd, nsk := range nonSignerKeys {
		ops := state.Operators[quorumID]
		if _, ok := ops[nonSignerOperatorIds[opInd]]; ok {
			signersAggKey.Sub(nsk)
		}
	}

	if attestation.QuorumAggPubKey[quorumID] == nil {
		return ErrAggPubKeyNotValid
	}

	ok, err := signersAggKey.VerifyEquivalence(attestation.SignersAggPubKey[quorumID])
	if err != nil {
		return err
	}
	if !ok {
		return ErrPubKeysNotEqual
	}

	// Verify the aggregated signature for the quorum
	ok = attestation.AggSignature[quorumID].Verify(attestation.SignersAggPubKey[quorumID], message)
	if !ok {
		return ErrAggSigNotValid
	}

	return nil
}

func (a *StdSignatureAggregator) ReceiveSignatures(
	ctx context.Context,
	state *IndexedOperatorState,
	message [32]byte,
	messageChan chan SigningMessage) (*QuorumAttestation, error) {

	quorumIDs, err := a.getQuorumIDs(state)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve quorum IDs: %w", err)
	}

	stakeSigned := make(map[QuorumID]*big.Int, len(quorumIDs))
	for _, quorumID := range quorumIDs {
		stakeSigned[quorumID] = big.NewInt(0)
	}

	attestation := &QuorumAttestation{
		QuorumAggPubKey:  make(map[QuorumID]*G1Point, len(quorumIDs)),
		SignersAggPubKey: make(map[QuorumID]*G2Point, len(quorumIDs)),
		AggSignature:     make(map[QuorumID]*Signature, len(quorumIDs)),
		QuorumResults:    make(map[QuorumID]*QuorumResult),
		SignerMap:        make(map[OperatorID]bool),
	}

	// validate and aggregate signatures
	numOperators := len(state.IndexedOperators)
	for numReply := 0; numReply < numOperators; numReply++ {
		a.processNextSignature(ctx, state, quorumIDs, message, stakeSigned, attestation, messageChan)
	}

	// Aggregate non-signers
	nonSignerKeys, nonSignerOperatorIds := a.aggregateNonSigners(state, attestation)

	// Determine the status of each quorum
	for _, quorumID := range quorumIDs {
		err = a.processQuorumStatus(
			state, quorumID, stakeSigned, attestation, message, nonSignerKeys, nonSignerOperatorIds)
		if err != nil {
			return nil, fmt.Errorf("failed to process quorum status: %w", err)
		}
	}

	return attestation, nil
}

func (a *StdSignatureAggregator) AggregateSignatures(
	ctx context.Context,
	ics IndexedChainState,
	referenceBlockNumber uint,
	quorumAttestation *QuorumAttestation,
	quorumIDs []QuorumID) (*SignatureAggregation, error) {

	// Aggregate the aggregated signatures. We reuse the first aggregated signature as the accumulator
	var aggSig *Signature
	for _, quorumID := range quorumIDs {
		if quorumAttestation.AggSignature[quorumID] == nil {
			a.Logger.Error("cannot aggregate signature for quorum because aggregated signature is nil",
				"quorumID", quorumID)
			continue
		}
		sig := quorumAttestation.AggSignature[quorumID]
		if aggSig == nil {
			aggSig = &Signature{sig.G1Point.Clone()}
		} else {
			aggSig.Add(sig.G1Point)
		}
	}

	// Aggregate the aggregated public keys. We reuse the first aggregated public key as the accumulator
	var aggPubKey *G2Point
	for _, quorumID := range quorumIDs {
		if quorumAttestation.SignersAggPubKey[quorumID] == nil {
			a.Logger.Error("cannot aggregate public key for quorum because signers aggregated public key is nil",
				"quorumID", quorumID)
			continue
		}
		apk := quorumAttestation.SignersAggPubKey[quorumID]
		if aggPubKey == nil {
			aggPubKey = apk.Clone()
		} else {
			aggPubKey.Add(apk)
		}
	}

	nonSignerKeys := make([]*G1Point, 0)
	indexedOperatorState, err := ics.GetIndexedOperatorState(ctx, referenceBlockNumber, quorumIDs)
	if err != nil {
		return nil, err
	}
	for id, op := range indexedOperatorState.IndexedOperators {
		_, found := quorumAttestation.SignerMap[id]
		if !found {
			nonSignerKeys = append(nonSignerKeys, op.PubkeyG1)
		}
	}

	// sort non signer keys according to how it's checked onchain
	// ref: https://github.com/Layr-Labs/eigenlayer-middleware/blob/m2-mainnet/src/BLSSignatureChecker.sol#L99
	sort.Slice(nonSignerKeys, func(i, j int) bool {
		hash1 := nonSignerKeys[i].Hash()
		hash2 := nonSignerKeys[j].Hash()
		// sort in ascending order
		return bytes.Compare(hash1[:], hash2[:]) == -1
	})

	quorumAggKeys := make(map[QuorumID]*G1Point, len(quorumIDs))
	for _, quorumID := range quorumIDs {
		if quorumAttestation.QuorumAggPubKey[quorumID] == nil {
			a.Logger.Error("cannot aggregate public key for quorum because aggregated public key is nil",
				"quorumID", quorumID)
			continue
		}
		quorumAggKeys[quorumID] = quorumAttestation.QuorumAggPubKey[quorumID]
	}

	quorumResults := make(map[QuorumID]*QuorumResult, len(quorumIDs))
	for _, quorumID := range quorumIDs {
		quorumResults[quorumID] = quorumAttestation.QuorumResults[quorumID]
	}

	return &SignatureAggregation{
		NonSigners:       nonSignerKeys,
		QuorumAggPubKeys: quorumAggKeys,
		AggPubKey:        aggPubKey,
		AggSignature:     aggSig,
		QuorumResults:    quorumResults,
	}, nil

}

func GetStakeThreshold(state *OperatorState, quorum QuorumID, quorumThreshold uint8) *big.Int {

	// Get stake threshold
	quorumThresholdBig := new(big.Int).SetUint64(uint64(quorumThreshold))
	stakeThreshold := new(big.Int)
	stakeThreshold.Mul(quorumThresholdBig, state.Totals[quorum].Stake)
	stakeThreshold = RoundUpDivideBig(stakeThreshold, new(big.Int).SetUint64(percentMultiplier))

	return stakeThreshold
}

func GetSignedPercentage(state *OperatorState, quorum QuorumID, stakeAmount *big.Int) uint8 {
	totalStake := state.Totals[quorum].Stake
	if totalStake.Cmp(big.NewInt(0)) == 0 {
		return 0
	}

	stakeAmount = stakeAmount.Mul(stakeAmount, new(big.Int).SetUint64(percentMultiplier))
	quorumThresholdBig := stakeAmount.Div(stakeAmount, totalStake)

	quorumThreshold := uint8(quorumThresholdBig.Uint64())

	return quorumThreshold
}
