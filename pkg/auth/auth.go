package auth

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// AuthMethod ...
type AuthMethod interface {
	Authenticate(ctx context.Context, md metadata.MD) (string, error)
}

// AuthConfig ...
type AuthConfig struct {
	Methods []AuthMethod
}

// NoAuth ... default auth-method
type NoAuth struct{}

// Authenticate ...
func (n *NoAuth) Authenticate(ctx context.Context, md metadata.MD) (string, error) {
	return "anonymous", nil // 返回默认用户名
}
