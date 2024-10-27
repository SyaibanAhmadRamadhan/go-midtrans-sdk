package midtrans

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

// ChargeQRIS
// list error type: ErrMarshaller, ErrApiCall, ErrDuplicateOrderID, ErrRateLimitExceeded
func (a *api) ChargeQRIS(ctx context.Context, input ChargeQRISInput) (output ChargeQRISOutput, err error) {
	var span trace.Span
	if a.usingOtel {
		if input.MetaData == nil {
			input.MetaData = make(map[string]string, 1)
		}
		input.MetaData["transparent"] = getTraceParent(ctx)
		ctx, span = otelTracer.Start(ctx, "ChargeQRIS", trace.WithAttributes())
		defer span.End()
	}

	req := ChargeRequest{
		PaymentType:       "qris",
		TransactionDetail: input.TransactionDetail,
		CustomExpiry:      input.CustomExpiry,
		ItemDetails:       input.ItemDetails,
		CustomerDetail:    input.CustomerDetail,
		MetaData:          input.MetaData,
		QRIS: &QRIS{
			Acquirer: input.Acquirer,
		},
	}

	err = a.v.Struct(req)
	if err != nil {
		if a.usingOtel {
			span.SetStatus(codes.Error, "invalid request")
			span.RecordError(err)
		}
		return output, err
	}

	reqMarshal, err := json.Marshal(req)
	if err != nil {
		if a.usingOtel {
			span.SetStatus(codes.Error, "failed to marshal request")
			span.RecordError(err)
		}
		return output, fmt.Errorf("%w:%w", ErrMarshaller, err)
	}

	resp, err := a.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(a.serverKey, "").
		SetBody(reqMarshal).
		Post(a.baseURI + "/v2/charge")
	if err != nil {
		if a.usingOtel {
			span.SetStatus(codes.Error, "HTTP request failed")
			span.RecordError(err)
		}
		return output, fmt.Errorf("%w:%w", ErrApiCall, err)
	}

	if a.usingOtel {
		ti := resp.Request.TraceInfo()
		span.SetAttributes(
			semconv.RPCSystemKey.String("http"),
			attribute.String("dns_lookup", ti.DNSLookup.String()),
			attribute.String("conn_time", ti.ConnTime.String()),
			attribute.String("tcp_conn_time", ti.TCPConnTime.String()),
			attribute.String("tls_handshake", ti.TLSHandshake.String()),
			attribute.String("server_time", ti.ServerTime.String()),
			attribute.String("response_time", ti.ResponseTime.String()),
			attribute.String("total_time", ti.TotalTime.String()),
			attribute.Bool("is_conn_reused", ti.IsConnReused),
			attribute.Bool("is_conn_was_idle", ti.IsConnWasIdle),
			attribute.String("conn_idle_time", ti.ConnIdleTime.String()),
			attribute.Int("request_attempt", ti.RequestAttempt),
			attribute.String("remote_addr", ti.RemoteAddr.String()),
			attribute.Int("request_size", len(reqMarshal)),
		)
	}

	if resp.IsError() {
		if a.usingOtel {
			span.SetStatus(codes.Error, fmt.Sprintf("request failed with status: %d", resp.StatusCode()))
			span.RecordError(fmt.Errorf("API call returned an error: %s", resp.String()))
		}
		return output, fmt.Errorf("%w: %s", ErrInternalServerPaymentGatewayError, resp.String())
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		err = ErrDuplicateOrderID
	case http.StatusTooManyRequests:
		err = ErrRateLimitExceeded
	case http.StatusBadRequest:
		output = ChargeQRISOutput{
			Error400Response: &Error400Response{},
		}
		err = json.Unmarshal(resp.Body(), &output.Error400Response)
		if err != nil {
			if a.usingOtel {
				span.SetStatus(codes.Error, "failed to unmarshal 400 response")
				span.RecordError(err)
			}
			return output, fmt.Errorf("%w:%w", ErrUnMarshaller, err)
		}
		if a.usingOtel {
			span.SetStatus(codes.Error, fmt.Sprintf("request failed with status: %d", resp.StatusCode()))
			span.RecordError(errors.New(resp.String()))
		}
		return output, err
	case http.StatusUnauthorized:
		output = ChargeQRISOutput{
			Error401Response: &Error401Response{},
		}
		err = json.Unmarshal(resp.Body(), &output.Error401Response)
		if err != nil {
			if a.usingOtel {
				span.SetStatus(codes.Error, "failed to unmarshal 401 response")
				span.RecordError(err)
			}
			return output, fmt.Errorf("%w:%w", ErrUnMarshaller, err)
		}
		return output, err
	}

	if err != nil {
		if a.usingOtel {
			span.SetStatus(codes.Error, fmt.Sprintf("request failed with status: %d", resp.StatusCode()))
			span.RecordError(errors.New(resp.String()))
		}
		return output, err
	}

	err = json.Unmarshal(resp.Body(), &output.ChargeQRISResponse)
	if err != nil {
		if a.usingOtel {
			span.SetStatus(codes.Error, "failed to unmarshal response")
			span.RecordError(err)
		}
		return output, fmt.Errorf("%w:%w", ErrUnMarshaller, err)
	}

	if a.usingOtel {
		span.SetStatus(codes.Ok, "request succeeded")
	}

	return
}

type ChargeQRISInput struct {
	TransactionDetail TransactionDetail
	ItemDetails       []ItemDetail
	Acquirer          string
	CustomerDetail    *CustomerDetail
	CustomExpiry      *CustomExpiry
	MetaData          map[string]string
}

type ChargeQRISOutput struct {
	Error400Response   *Error400Response
	Error401Response   *Error401Response
	ChargeQRISResponse ChargeQRISResponse
}
