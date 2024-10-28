package coreapi_midtrans

import "context"

type PaymentAPI interface {
	ChargeEWalletAPI
}

type ChargeEWalletAPI interface {
	// ChargeQRIS
	// list error type: ErrMarshaller, ErrApiCall, ErrRateLimitExceeded
	ChargeQRIS(ctx context.Context, input ChargeQRISInput) (output ChargeQRISOutput, err error)

	// ChargeGoPay
	// list error type: ErrMarshaller, ErrApiCall, ErrRateLimitExceeded
	ChargeGoPay(ctx context.Context, input ChargeGoPayInput) (output ChargeGoPayOutput, err error)

	// ChargeShopeePay
	// list error type: ErrMarshaller, ErrApiCall, ErrRateLimitExceeded
	ChargeShopeePay(ctx context.Context, input ChargeShopeePayInput) (output ChargeShopeePayOutput, err error)
}
