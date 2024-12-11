package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

type BasicAuth struct {
	Users map[string]string // 用户名-密码对
}

func (b *BasicAuth) Authenticate(ctx context.Context, md metadata.MD) (string, error) {
	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return "", errors.New("missing authorization header")
	}
	// Basic auth 格式: Basic base64(username:password)
	encoded := strings.TrimPrefix(authHeader[0], "Basic ")

	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	decoded := string(decodedBytes)

	parts := strings.SplitN(decoded, ":", 2)
	if len(parts) != 2 || b.Users[parts[0]] != parts[1] {
		return "", errors.New("invalid username or password")
	}
	return parts[0], nil
}
