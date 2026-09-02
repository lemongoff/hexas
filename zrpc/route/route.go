// Package route defines per-call RPC endpoint routing directives.
package route

import (
	"context"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// Mode controls how a client behaves when the requested endpoint is not ready.
type Mode uint8

const (
	// Prefer uses the requested endpoint when available and otherwise falls back
	// to the configured load-balancing policy.
	Prefer Mode = iota
	// Require fails the call when the requested endpoint is not ready.
	Require
)

// Target is a per-call endpoint routing directive.
type Target struct {
	InstanceID string
	Mode       Mode
}

type targetContextKey struct{}
type instanceIDAttributeKey struct{}

// WithTarget returns a context that routes the RPC to instanceID.
func WithTarget(ctx context.Context, instanceID string, mode Mode) context.Context {
	return context.WithValue(ctx, targetContextKey{}, Target{InstanceID: instanceID, Mode: mode})
}

// TargetFromContext returns the routing directive carried by ctx.
func TargetFromContext(ctx context.Context) (Target, bool) {
	if ctx == nil {
		return Target{}, false
	}
	target, ok := ctx.Value(targetContextKey{}).(Target)
	return target, ok
}

// SetInstanceID attaches a stable service instance identity to an address.
func SetInstanceID(address resolver.Address, instanceID string) resolver.Address {
	if address.Attributes == nil {
		address.Attributes = attributes.New(instanceIDAttributeKey{}, instanceID)
	} else {
		address.Attributes = address.Attributes.WithValue(instanceIDAttributeKey{}, instanceID)
	}
	return address
}

// InstanceID returns the stable service instance identity attached to address.
func InstanceID(address resolver.Address) string {
	if address.Attributes == nil {
		return ""
	}
	instanceID, _ := address.Attributes.Value(instanceIDAttributeKey{}).(string)
	return instanceID
}
