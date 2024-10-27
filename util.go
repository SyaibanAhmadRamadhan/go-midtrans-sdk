package midtrans

import "context"

func getTraceParent(ctx context.Context) string {
	traceParent, ok := ctx.Value("traceparent").(string)
	if !ok || traceParent == "" {
		return ""
	}

	return traceParent
}
