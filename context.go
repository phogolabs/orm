package orm

import (
	"context"
)

var (
	// GatewayCtxKey is the context.Context key to store the request gateway.
	GatewayCtxKey = &contextKey{"Gateway"}
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "orm/middleware context value " + k.name
}

// SetContext sets a log entry into the provided context
func SetContext(ctx context.Context, gateway *Gateway) context.Context {
	return context.WithValue(ctx, GatewayCtxKey, gateway)
}

// GetContext returns the log Entry found in the context,
// or a new Default log Entry if none is found
func GetContext(ctx context.Context) *Gateway {
	v := ctx.Value(GatewayCtxKey)

	if v == nil {
		return nil
	}

	return v.(*Gateway)
}
