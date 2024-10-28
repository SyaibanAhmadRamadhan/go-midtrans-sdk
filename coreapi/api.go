package coreapi_midtrans

import (
	"github.com/SyaibanAhmadRamadhan/go-midtrans-sdk"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
	"time"
)

type optApiFunc func(*api)

func WithRetryMechanism(retryCount int, retryWaitTime time.Duration) optApiFunc {
	return func(api *api) {
		api.restyClient.SetRetryCount(retryCount)
		api.restyClient.SetRetryWaitTime(retryWaitTime)
	}
}

func ServerKey(s string) optApiFunc {
	return func(api *api) {
		api.serverKey = s
	}
}

func ProductionLive() optApiFunc {
	return func(api *api) {
		api.baseURI = "https://api.midtrans.com"
	}
}

func WithTraceParentKey(k string) optApiFunc {
	return func(api *api) {
		api.traceParentKey = k
	}
}

func WithOtel() optApiFunc {
	return func(api *api) {
		api.usingOtel = true
		api.tracing = midtrans.NewOtelTracing()
	}
}

type api struct {
	v              *validator.Validate
	trans          ut.Translator
	baseURI        string
	serverKey      string
	restyClient    *resty.Client
	tracing        midtrans.Tracing
	traceParentKey string
	// config
	usingOtel bool
}

func NewAPI(opts ...optApiFunc) *api {
	v, t := midtrans.NewValidator()
	a := &api{
		v:              v,
		trans:          t,
		traceParentKey: "transparent",
		baseURI:        "https://api.sandbox.midtrans.com",
		serverKey:      os.Getenv("MIDTRANS_API_KEY"),
		restyClient: resty.New().AddRetryCondition(func(response *resty.Response, err error) bool {
			return response.StatusCode() >= 500 || response.StatusCode() == http.StatusTooManyRequests
		}),
	}

	for _, opt := range opts {
		opt(a)
	}

	if a.tracing == nil {
		a.tracing = midtrans.NewOtelTracing()
		a.usingOtel = true
	}
	return a
}
