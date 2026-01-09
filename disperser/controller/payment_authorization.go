package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand/ondemandvalidation"
	"github.com/Layr-Labs/eigenda/core/payments/reservation/reservationvalidation"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	payments "github.com/Layr-Labs/eigenda/disperser/controller/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/prometheus/client_golang/prometheus"
)

// PaymentAuthorizationConfig contains configuration for building a payment authorization handler
type PaymentAuthorizationConfig struct {
	OnDemandConfig                 ondemandvalidation.OnDemandLedgerCacheConfig
	OnDemandMetererConfig          meterer.OnDemandMetererConfig
	ReservationConfig              reservationvalidation.ReservationLedgerCacheConfig
	EnablePerAccountPaymentMetrics bool
}

// Verify validates the PaymentAuthorizationConfig
func (c *PaymentAuthorizationConfig) Verify() error {
	if err := c.OnDemandConfig.Verify(); err != nil {
		return fmt.Errorf("on-demand config: %w", err)
	}
	if err := c.OnDemandMetererConfig.Verify(); err != nil {
		return fmt.Errorf("on-demand meterer config: %w", err)
	}
	if err := c.ReservationConfig.Verify(); err != nil {
		return fmt.Errorf("reservation config: %w", err)
	}
	return nil
}

// DefaultPaymentAuthorizationConfig returns a new PaymentAuthorizationConfig with default values
func DefaultPaymentAuthorizationConfig() *PaymentAuthorizationConfig {
	onDemandConfig := ondemandvalidation.OnDemandLedgerCacheConfig{
		MaxLedgers:        1024,
		OnDemandTableName: "",
		UpdateInterval:    30 * time.Second,
	}

	onDemandMetererConfig := meterer.OnDemandMetererConfig{
		RefreshInterval: meterer.DefaultOnDemandMetererRefreshInterval,
		FuzzFactor:      meterer.DefaultOnDemandMetererFuzzFactor,
	}

	reservationConfig := reservationvalidation.ReservationLedgerCacheConfig{
		MaxLedgers:           1024,
		BucketCapacityPeriod: 90 * time.Second,
		OverfillBehavior:     ratelimit.OverfillOncePermitted,
		UpdateInterval:       30 * time.Second,
	}

	return &PaymentAuthorizationConfig{
		OnDemandConfig:                 onDemandConfig,
		OnDemandMetererConfig:          onDemandMetererConfig,
		ReservationConfig:              reservationConfig,
		EnablePerAccountPaymentMetrics: true,
	}
}

// BuildPaymentAuthorizationHandler creates a payment authorization handler with the given configuration.
// If metricsRegistry is nil, metrics will be disabled (useful for tests).
func BuildPaymentAuthorizationHandler(
	ctx context.Context,
	logger logging.Logger,
	config PaymentAuthorizationConfig,
	contractDirectory *directory.ContractDirectory,
	ethClient common.EthClient,
	awsDynamoClient *awsdynamodb.Client,
	metricsRegistry *prometheus.Registry,
	userAccountRemapping map[string]string,
) (*payments.PaymentAuthorizationHandler, error) {
	paymentVaultAddress, err := contractDirectory.GetContractAddress(ctx, directory.PaymentVault)
	if err != nil {
		return nil, fmt.Errorf("get PaymentVault address: %w", err)
	}

	paymentVault, err := vault.NewPaymentVault(logger, ethClient, paymentVaultAddress)
	if err != nil {
		return nil, fmt.Errorf("create payment vault: %w", err)
	}

	// Create on-demand meterer (use nil metrics if registry is nil)
	var onDemandMetererMetrics *meterer.OnDemandMetererMetrics
	if metricsRegistry != nil {
		onDemandMetererMetrics = meterer.NewOnDemandMetererMetrics(
			metricsRegistry,
			"eigenda_controller",
			"authorize_payments",
		)
	}

	onDemandMeterer, err := meterer.NewOnDemandMeterer(
		ctx,
		paymentVault,
		time.Now,
		onDemandMetererMetrics,
		config.OnDemandMetererConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("create on-demand meterer: %w", err)
	}

	// Create on-demand validator (use nil metrics if registry is nil)
	var onDemandValidatorMetrics *ondemandvalidation.OnDemandValidatorMetrics
	var onDemandCacheMetrics *ondemandvalidation.OnDemandCacheMetrics
	if metricsRegistry != nil {
		onDemandValidatorMetrics = ondemandvalidation.NewOnDemandValidatorMetrics(
			metricsRegistry,
			"eigenda_controller",
			"authorize_payments",
			config.EnablePerAccountPaymentMetrics,
			userAccountRemapping,
		)
		onDemandCacheMetrics = ondemandvalidation.NewOnDemandCacheMetrics(
			metricsRegistry,
			"eigenda_controller",
			"authorize_payments",
		)
	}

	onDemandValidator, err := ondemandvalidation.NewOnDemandPaymentValidator(
		ctx,
		logger,
		config.OnDemandConfig,
		paymentVault,
		awsDynamoClient,
		onDemandValidatorMetrics,
		onDemandCacheMetrics,
	)
	if err != nil {
		return nil, fmt.Errorf("create on-demand payment validator: %w", err)
	}

	// Create reservation validator (use nil metrics if registry is nil)
	var reservationValidatorMetrics *reservationvalidation.ReservationValidatorMetrics
	var reservationCacheMetrics *reservationvalidation.ReservationCacheMetrics
	if metricsRegistry != nil {
		reservationValidatorMetrics = reservationvalidation.NewReservationValidatorMetrics(
			metricsRegistry,
			"eigenda_controller",
			"authorize_payments",
			config.EnablePerAccountPaymentMetrics,
			userAccountRemapping,
		)
		reservationCacheMetrics = reservationvalidation.NewReservationCacheMetrics(
			metricsRegistry,
			"eigenda_controller",
			"authorize_payments",
		)
	}

	reservationValidator, err := reservationvalidation.NewReservationPaymentValidator(
		ctx,
		logger,
		config.ReservationConfig,
		paymentVault,
		time.Now,
		reservationValidatorMetrics,
		reservationCacheMetrics,
	)
	if err != nil {
		return nil, fmt.Errorf("create reservation payment validator: %w", err)
	}

	return payments.NewPaymentAuthorizationHandler(
		onDemandMeterer,
		onDemandValidator,
		reservationValidator,
	), nil
}
