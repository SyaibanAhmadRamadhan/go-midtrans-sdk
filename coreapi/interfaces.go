package coreapi

import "context"

type ChargeEWalletAPI interface {
	ChargeQRIS(ctx context.Context, input ChargeQRISInput) (output ChargeQRISOutput, err error)
}
