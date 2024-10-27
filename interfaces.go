package midtrans

import "context"

type PaymentAPI interface {
	//ChargeTransaction(ctx context.Context, input ChargeTransactionInput) (output ChargeTransactionOutput, err error)
}

type ChargeEWalletAPI interface {
	ChargeQRIS(ctx context.Context, input ChargeQRISInput) (output ChargeQRISOutput, err error)
}
