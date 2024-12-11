package auth

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

// JWT or Bearer Token
type BearerAuth struct {
	ValidateToken func(token string) (string, error) // Token 验证函数
}

func (b *BearerAuth) Authenticate(ctx context.Context, md metadata.MD) (string, error) {
	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return "", errors.New("missing authorization header")
	}
	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	return b.ValidateToken(token) // 验证 Token 并提取用户名
}
