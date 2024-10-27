package midtrans

import (
	"context"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var otelTracer = otel.Tracer("github.com/SyaibanAhmadRamadhan/go-midtrans-sdk")

type otelTracing struct{}

func NewOtelTracing() *otelTracing {
	return &otelTracing{}
}

func (o *otelTracing) StartTrace(ctx context.Context, traceName string) context.Context {
	ctx, _ = otelTracer.Start(ctx, traceName, trace.WithAttributes())
	return ctx
}

func (o *otelTracing) SetRestyTraceInfo(ctx context.Context, resp *resty.Response) {
	span := trace.SpanFromContext(ctx)
	if resp == nil || !span.IsRecording() {
		return
	}
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
	)
}
func (o *otelTracing) EndTrace(ctx context.Context, err error, msg string) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	if err != nil {
		span.SetStatus(codes.Error, msg)
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Ok, msg)
	}

	span.End()
}

func (o *otelTracing) SetRespBody(ctx context.Context, resp *resty.Response) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	span.SetAttributes(attribute.String(
		"resp.body", string(resp.Body())),
	)
}
