package main

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/ejector"
)

const ejectorEnvVarPrefix = "EJECTOR"

func main() {

	cfg, err := config.Bootstrap(ejector.DefaultEjectorConfig, ejectorEnvVarPrefix)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Ejector configuration: %+v\n", cfg) // TODO

	// ctx := context.Background()

	// TODO initialize
	// ejector := ejector.NewEjector(
	// 	ctx context.Context,
	// 	logger logging.Logger,
	// 	ejectionManager *ejector.ThreadedEjectionManager,
	// 	signingRateLookupV1 ejector.SigningRateLookup,
	// 	signingRateLookupV2 ejector.SigningRateLookup,
	// 	period time.Duration,
	// 	ejectionCriteriaTimeWindow time.Duration,
	// 	validatorIDToAddressCache *eth.ValidatorIDToAddressCache,
	// )

	// TODO
}
