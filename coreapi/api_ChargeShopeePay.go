package coreapi_midtrans

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-midtrans-sdk"
)

func (a *api) ChargeShopeePay(ctx context.Context, input ChargeShopeePayInput) (output ChargeShopeePayOutput, err error) {
	if a.usingOtel {
		if input.MetaData == nil {
			input.MetaData = make(map[string]string, 1)
		}
		input.MetaData["transparent"] = midtrans.GetTraceParent(ctx, a.traceParentKey)
	}
	ctx = a.tracing.StartTrace(ctx, "ChargeGoPay")

	req := midtrans.ChargeRequest{
		PaymentType:       "shopeepay",
		TransactionDetail: input.TransactionDetail,
		CustomExpiry:      input.CustomExpiry,
		ItemDetails:       input.ItemDetails,
		CustomerDetail:    input.CustomerDetail,
		MetaData:          input.MetaData,
		ShopeePay:         input.ShopeePay,
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
			switch actionResp.Name {
			case "deeplink-redirect":
				output.ResponseSuccess.ActionDeepLinkRedirect = actionResp
			}
		}
	}
	return
}

type ChargeShopeePayInput struct {
	TransactionDetail midtrans.TransactionDetail
	ItemDetails       []midtrans.ItemDetail
	CustomerDetail    *midtrans.CustomerDetail
	CustomExpiry      *midtrans.CustomExpiry
	MetaData          map[string]string
	ShopeePay         *midtrans.ShopeePay
}

type ChargeShopeePayOutput struct {
	ErrorBadReqResponse *midtrans.ErrorBadReqResponse
	ResponseSuccess     midtrans.ShopeePayResponse
}
