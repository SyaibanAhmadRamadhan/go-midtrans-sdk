package midtrans

import (
	"context"
	"github.com/go-resty/resty/v2"
)

type PaymentAPI interface {
	//ChargeTransaction(ctx context.Context, input ChargeTransactionInput) (output ChargeTransactionOutput, err error)
}

type Tracing interface {
	StartTrace(ctx context.Context, traceName string) context.Context
	SetRestyTraceInfo(ctx context.Context, resp *resty.Response)
	SetRespBody(ctx context.Context, resp *resty.Response)
	EndTrace(ctx context.Context, err error, msg string)
}
