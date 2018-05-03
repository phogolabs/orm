package oak

import (
	"context"
	"net/http"
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
	return "oak/middleware context value " + k.name
}

// GatewayHandler is a middleware that sets a given gateway in a HTTP context chain.
func GatewayHandler(gateway *Gateway) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = WithGateway(r, gateway)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// GetGateway returns the in-context Gateway for a request.
func GetGateway(r *http.Request) *Gateway {
	entry, _ := r.Context().Value(GatewayCtxKey).(*Gateway)
	return entry
}

// WithGateway sets the in-context LogEntry for a request.
func WithGateway(r *http.Request, gateway *Gateway) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), GatewayCtxKey, gateway))
	return r
}
