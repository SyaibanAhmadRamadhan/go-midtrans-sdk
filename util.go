package midtrans

import "context"

func GetTraceParent(ctx context.Context, key string) string {
	traceParent, ok := ctx.Value(key).(string)
	if !ok || traceParent == "" {
		return ""
	}

	return traceParent
}

func GetGenerateQRCodeUrl(actionResp []ActionResponse) string {
	if actionResp == nil {
		return ""
	}

	for _, response := range actionResp {
		if response.Name == "generate-qr-code" {
			return response.URL
		}
	}

	return ""
}
