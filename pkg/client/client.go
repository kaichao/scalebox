// Package client provides a gRPC client singleton for scalebox SDK packages.
// Use Default() for a shared connection across all pkg/* wrappers.
package client

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/kaichao/gopkg/errors"
	pb "github.com/kaichao/scalebox/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Config holds gRPC client connection parameters.
type Config struct {
	ServerAddr string
	TLSEnabled bool
	CAFile     string
	AuthToken  string
}

// LoadConfig reads gRPC client configuration from environment variables.
func LoadConfig() Config {
	return Config{
		ServerAddr: getServerAddr(),
		TLSEnabled: os.Getenv("GRPC_TLS_ENABLED") == "true",
		CAFile:     os.Getenv("GRPC_TLS_CA_FILE"),
		AuthToken:  os.Getenv("AUTH_TOKEN"),
	}
}

// Client wraps a single gRPC connection and available service clients.
type Client struct {
	conn *grpc.ClientConn

	Coordination pb.CoordinationServiceClient
	App          pb.AppServiceClient
}

// New creates a Client with the given config.
func New(cfg Config) (*Client, error) {
	opts := buildDialOpts(cfg)

	conn, err := grpc.NewClient(cfg.ServerAddr, opts...)
	if err != nil {
		return nil, errors.WrapE(err, "grpc.NewClient", "server", cfg.ServerAddr)
	}

	return &Client{
		conn:          conn,
		Coordination:  pb.NewCoordinationServiceClient(conn),
		App:           pb.NewAppServiceClient(conn),
	}, nil
}

// Close closes the underlying gRPC connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// ── Singleton ──────────────────────────────────────────────────────

var (
	defaultClient *Client
	defaultOnce   sync.Once
	defaultErr    error
)

// Default returns the singleton Client, creating it on first call.
func Default() (*Client, error) {
	defaultOnce.Do(func() {
		defaultClient, defaultErr = New(LoadConfig())
	})
	return defaultClient, defaultErr
}

// ResetDefault closes and clears the singleton. For testing.
func ResetDefault() {
	if defaultClient != nil {
		defaultClient.Close()
	}
	defaultClient = nil
	defaultOnce = sync.Once{}
	defaultErr = nil
}

// ── Internal helpers ───────────────────────────────────────────────

func buildDialOpts(cfg Config) []grpc.DialOption {
	var opts []grpc.DialOption

	if cfg.TLSEnabled && cfg.CAFile != "" {
		creds, err := credentials.NewClientTLSFromFile(cfg.CAFile, "")
		if err != nil {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(creds))
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if cfg.AuthToken != "" {
		opts = append(opts, grpc.WithPerRPCCredentials(&jwtCredentials{
			token:      cfg.AuthToken,
			tlsEnabled: cfg.TLSEnabled,
		}))
	}

	return opts
}

func getServerAddr() string {
	addr := os.Getenv("GRPC_SERVER")
	if addr == "" {
		addr = "localhost"
	}
	if !strings.Contains(addr, ":") {
		addr += ":50051"
	}
	return addr
}

type jwtCredentials struct {
	token      string
	tlsEnabled bool
}

func (c *jwtCredentials) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{"authorization": "Bearer " + c.token}, nil
}

func (c *jwtCredentials) RequireTransportSecurity() bool {
	return c.tlsEnabled
}
