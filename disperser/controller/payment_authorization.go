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
	// Connfiguration for on-demand payment validation.
	OnDemand ondemandvalidation.OnDemandLedgerCacheConfig
	// Configuration for reservation payment validation.
	Reservation reservationvalidation.ReservationLedgerCacheConfig
	// If true, enable a metric per user account for payment validation and authorization.
	// Resulting metric may potentially have high cardinality.
	PerAccountMetrics bool
}

// Verify validates the PaymentAuthorizationConfig
func (c *PaymentAuthorizationConfig) Verify() error {
	if err := c.OnDemand.Verify(); err != nil {
		return fmt.Errorf("on-demand config: %w", err)
	}
	if err := c.Reservation.Verify(); err != nil {
		return fmt.Errorf("reservation config: %w", err)
	}
	return nil
}

// DefaultPaymentAuthorizationConfig returns a new PaymentAuthorizationConfig with default values
func DefaultPaymentAuthorizationConfig() PaymentAuthorizationConfig {
	onDemandConfig := ondemandvalidation.OnDemandLedgerCacheConfig{
		MaxLedgers:        1024,
		OnDemandTableName: "",
		UpdateInterval:    30 * time.Second,
	}

	reservationConfig := reservationvalidation.ReservationLedgerCacheConfig{
		MaxLedgers:           1024,
		BucketCapacityPeriod: 90 * time.Second,
		OverfillBehavior:     ratelimit.OverfillOncePermitted,
		UpdateInterval:       30 * time.Second,
	}

	return PaymentAuthorizationConfig{
		OnDemand:          onDemandConfig,
		Reservation:       reservationConfig,
		PerAccountMetrics: true,
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
			config.PerAccountMetrics,
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
		config.OnDemand,
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
			config.PerAccountMetrics,
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
		config.Reservation,
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
