package midtrans

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel"
	"net/http"
	"os"
	"time"
)

var otelTracer = otel.Tracer("github.com/SyaibanAhmadRamadhan/go-midtrans-sdk")

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

func WithOtel() optApiFunc {
	return func(api *api) {
		api.usingOtel = true
	}
}

type api struct {
	v           *validator.Validate
	trans       ut.Translator
	baseURI     string
	serverKey   string
	restyClient *resty.Client

	// config
	usingOtel bool
}

func NewAPI(opts ...optApiFunc) *api {
	v, t := NewValidator()
	a := &api{
		v:         v,
		trans:     t,
		baseURI:   "https://api.sandbox.midtrans.com",
		serverKey: os.Getenv("MIDTRANS_API_KEY"),
		restyClient: resty.New().AddRetryCondition(func(response *resty.Response, err error) bool {
			return response.StatusCode() <= 500 || response.StatusCode() == http.StatusTooManyRequests
		}),
	}

	for _, opt := range opts {
		opt(a)
	}
	return a
}
