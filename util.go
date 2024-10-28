package midtrans

import "context"

func GetTraceParent(ctx context.Context, key string) string {
	traceParent, ok := ctx.Value(key).(string)
	if !ok || traceParent == "" {
		return ""
	}

	return traceParent
}
