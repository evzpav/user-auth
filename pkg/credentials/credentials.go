package credentials

import (
	"context"
)

const ClientCredentials = "client_credentials"

func GetClientCredentials(ctx context.Context) (string, bool) {
	cred := ctx.Value(ClientCredentials)
	if cred == nil {
		return "", false
	}
	return cred.(string), true
}
