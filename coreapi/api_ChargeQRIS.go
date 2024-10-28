package coreapi_midtrans

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-midtrans-sdk"
)

// ChargeQRIS
// list error type: ErrMarshaller, ErrApiCall, ErrDuplicateOrderID, ErrRateLimitExceeded
func (a *api) ChargeQRIS(ctx context.Context, input ChargeQRISInput) (output ChargeQRISOutput, err error) {
	if a.usingOtel {
		if input.MetaData == nil {
			input.MetaData = make(map[string]string, 1)
		}
		input.MetaData["transparent"] = midtrans.GetTraceParent(ctx, a.traceParentKey)
	}
	ctx = a.tracing.StartTrace(ctx, "ChargeQRIS")

	req := midtrans.ChargeRequest{
		PaymentType:       "qris",
		TransactionDetail: input.TransactionDetail,
		CustomExpiry:      input.CustomExpiry,
		ItemDetails:       input.ItemDetails,
		CustomerDetail:    input.CustomerDetail,
		MetaData:          input.MetaData,
		QRIS: &midtrans.QRIS{
			Acquirer: input.Acquirer,
		},
	}

	err = a.v.Struct(req)
	if err != nil {
		a.tracing.EndTrace(ctx, err, "invalid request")
		return output, err
	}

	reqMarshal, err := json.Marshal(req)
	if err != nil {
		a.tracing.EndTrace(ctx, err, "failed to marshal request")
		return output, fmt.Errorf("%w:%w", midtrans.ErrMarshaller, err)
	}

	resp, err := a.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(a.serverKey, "").
		SetBody(reqMarshal).
		Post(a.baseURI + "/v2/charge")
	if err != nil {
		a.tracing.EndTrace(ctx, err, "HTTP request failed")
		return output, fmt.Errorf("%w:%w", midtrans.ErrApiCall, err)
	}

	a.tracing.SetRestyTraceInfo(ctx, resp)

	output.ErrorBadReqResponse, err = a.catchResponse(ctx, resp, &output.ResponseSuccess)
	if err != nil {
		return output, err
	}

	if output.ErrorBadReqResponse == nil {
		for _, actionResp := range output.ResponseSuccess.Actions {
			if actionResp.Name == "generate-qr-code" {
				output.ResponseSuccess.ActionGenerateQRCode = actionResp
			}
		}
	}
	return
}

type ChargeQRISInput struct {
	TransactionDetail midtrans.TransactionDetail
	ItemDetails       []midtrans.ItemDetail
	Acquirer          string
	CustomerDetail    *midtrans.CustomerDetail
	CustomExpiry      *midtrans.CustomExpiry
	MetaData          map[string]string
}

type ChargeQRISOutput struct {
	ErrorBadReqResponse *midtrans.ErrorBadReqResponse
	ResponseSuccess     midtrans.ChargeQRISResponse
}
