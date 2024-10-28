package midtrans

import (
	"context"
	"github.com/go-resty/resty/v2"
)

type Tracing interface {
	StartTrace(ctx context.Context, traceName string) context.Context
	SetRestyTraceInfo(ctx context.Context, resp *resty.Response)
	EndTrace(ctx context.Context, err error, msg string)
}
