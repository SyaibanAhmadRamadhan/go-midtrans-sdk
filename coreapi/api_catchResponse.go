package coreapi_midtrans

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-midtrans-sdk"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type output struct {
	errorBadReqResponse *midtrans.ErrorBadReqResponse
}

func (a *api) catchResponse(ctx context.Context, resp *resty.Response, respSuccess any) (*midtrans.ErrorBadReqResponse, error) {
	if resp.StatusCode() >= 500 {
		err := fmt.Errorf("%w: %s", midtrans.ErrInternalServerPaymentGatewayError, resp.String())
		a.tracing.EndTrace(ctx, err, fmt.Sprintf("request failed with status: %d", resp.StatusCode()))
		return nil, err
	}

	if resp.StatusCode() == http.StatusTooManyRequests {
		a.tracing.EndTrace(ctx, errors.New(resp.String()), fmt.Sprintf("request failed with status: %d", resp.StatusCode()))
		return nil, midtrans.ErrRateLimitExceeded
	}

	errorBadReqResponse := &midtrans.ErrorBadReqResponse{}
	err := json.Unmarshal(resp.Body(), &errorBadReqResponse)
	if err != nil {
		a.tracing.EndTrace(ctx, err, "failed to unmarshal ErrorBadReqResponse response")
		return nil, fmt.Errorf("%w:%w", midtrans.ErrUnMarshaller, err)
	}

	if errorBadReqResponse.StatusCode == "200" {
		err = json.Unmarshal(resp.Body(), respSuccess)
		if err != nil {
			a.tracing.EndTrace(ctx, err, "failed to unmarshal ErrorBadReqResponse response")
			return nil, fmt.Errorf("%w:%w", midtrans.ErrUnMarshaller, err)
		}
		a.tracing.EndTrace(ctx, nil, "request succeeded")
		return nil, nil
	}

	a.tracing.EndTrace(ctx, errors.New(resp.String()), fmt.Sprintf("request failed with status: %d", resp.StatusCode()))
	return errorBadReqResponse, nil
}
