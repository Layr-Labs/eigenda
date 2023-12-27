package dataapi

import (
	"context"
)

func (s *server) getDeregisterdOperatorsLast14days(ctx context.Context) ([]*DeregisteredOperatorInfo, error) {

	indexedDeregisteredOperatorState, err := s.subgraphClient.QueryIndexedDeregisteredOperatorsInTheLast14Days(ctx)
	if err != nil {
		return nil, err
	}

	// Convert the map to a slice.
	operators := indexedDeregisteredOperatorState.Operators
	deRegisteredOperators := make([]*DeregisteredOperatorInfo, 0, len(operators))

	for _, operatorInfo := range operators {
		deRegisteredOperators = append(deRegisteredOperators, operatorInfo)
	}

	return deRegisteredOperators, nil
}
