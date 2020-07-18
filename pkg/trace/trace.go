package trace

import (
	"context"
)

const RequestIDHeader = "X-Request-Id"

func GetRequestID(ctx context.Context) (string, bool) {
	reqIDHeader := ctx.Value(RequestIDHeader)
	if reqIDHeader == nil {
		return "", false
	}
	return reqIDHeader.(string), true
}
