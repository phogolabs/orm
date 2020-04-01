package orm

import (
	"context"
)

var (
	// GatewayCtxKey is the context.Context key to store the request gateway.
	GatewayCtxKey = &contextKey{"gateway"}
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "orm: context value " + k.name
}

// NewGatewayContext sets a log entry into the provided context
func NewGatewayContext(ctx context.Context, gateway *Gateway) context.Context {
	return context.WithValue(ctx, GatewayCtxKey, gateway)
}

// GatewayFromContext returns the log Entry found in the context,
// or a new Default log Entry if none is found
func GatewayFromContext(ctx context.Context) *Gateway {
	v := ctx.Value(GatewayCtxKey)

	if v == nil {
		return nil
	}

	return v.(*Gateway)
}
