package coreapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-midtrans-sdk"
	"net/http"
)

// ChargeQRIS
// list error type: ErrMarshaller, ErrApiCall, ErrDuplicateOrderID, ErrRateLimitExceeded
func (a *api) ChargeQRIS(ctx context.Context, input ChargeQRISInput) (output ChargeQRISOutput, err error) {
	if a.usingOtel {
		if input.MetaData == nil {
			input.MetaData = make(map[string]string, 1)
		}
		input.MetaData["transparent"] = midtrans.GetTraceParent(ctx)
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

	if resp.IsError() {
		err = fmt.Errorf("%w: %s", midtrans.ErrInternalServerPaymentGatewayError, resp.String())
		a.tracing.EndTrace(ctx, err, fmt.Sprintf("request failed with status: %d", resp.StatusCode()))
		return
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		a.tracing.EndTrace(ctx, errors.New(resp.String()), fmt.Sprintf("request failed with status: %d", resp.StatusCode()))
		return output, midtrans.ErrDuplicateOrderID
	case http.StatusTooManyRequests:
		a.tracing.EndTrace(ctx, errors.New(resp.String()), fmt.Sprintf("request failed with status: %d", resp.StatusCode()))
		return output, midtrans.ErrRateLimitExceeded
	case http.StatusBadRequest:
		output = ChargeQRISOutput{
			Error400Response: &midtrans.Error400Response{},
		}
		err = json.Unmarshal(resp.Body(), &output.Error400Response)
		if err != nil {
			a.tracing.EndTrace(ctx, err, "failed to unmarshal 400 response")
			return output, fmt.Errorf("%w:%w", midtrans.ErrUnMarshaller, err)
		}
		a.tracing.EndTrace(ctx, errors.New(resp.String()), fmt.Sprintf("request failed with status: %d", resp.StatusCode()))
		return output, err
	case http.StatusUnauthorized:
		output = ChargeQRISOutput{
			Error401Response: &midtrans.Error401Response{},
		}
		err = json.Unmarshal(resp.Body(), &output.Error401Response)
		if err != nil {
			a.tracing.EndTrace(ctx, err, "failed to unmarshal 401 response")
			return output, fmt.Errorf("%w:%w", midtrans.ErrUnMarshaller, err)
		}
		a.tracing.EndTrace(ctx, errors.New(resp.String()), fmt.Sprintf("request failed with status: %d", resp.StatusCode()))
		return output, err
	}

	err = json.Unmarshal(resp.Body(), &output.ChargeQRISResponse)
	if err != nil {
		a.tracing.EndTrace(ctx, err, "failed to unmarshal response")
		return output, fmt.Errorf("%w:%w", midtrans.ErrUnMarshaller, err)
	}

	a.tracing.SetRespBody(ctx, resp)
	a.tracing.EndTrace(ctx, nil, "request succeeded")
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
	Error400Response   *midtrans.Error400Response
	Error401Response   *midtrans.Error401Response
	ChargeQRISResponse midtrans.ChargeQRISResponse
}
