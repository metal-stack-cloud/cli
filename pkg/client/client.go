package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	v1 "github.com/metal-stack-cloud/api-server/api/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GRPCScheme the scheme to talk to the duros api endpoint, can be plaintext or https
type GRPCScheme string

const (
	// GRPC defines a plaitext communication
	GRPC GRPCScheme = "grpc"
	// GRPCS defines https protocol for the communication
	GRPCS GRPCScheme = "grpcs"

	defaultUserAgent = "api-server-go"
)

// DialConfig is the configuration to create a duros-api connection
type DialConfig struct {
	Endpoint        string
	Scheme          GRPCScheme
	Token           string
	Credentials     *Credentials
	ByteCredentials *ByteCredentials
	Log             *zap.SugaredLogger
	// UserAgent to use, if empty duros-go is used
	UserAgent string
}

// Credentials specify the TLS Certificate based authentication for the grpc connection
// If you provide credentials, provide either these or byte credentials but not both.
type Credentials struct {
	ServerName string
	Certfile   string
	Keyfile    string
	CAFile     string
}

// Credentials specify the TLS Certificate based authentication for the grpc connection
// without having to use certificate files.
// If you provide credentials, provide either these or file path credentials but not both.
type ByteCredentials struct {
	ServerName string
	Cert       []byte
	Key        []byte
	CA         []byte
}

// Client defines the client API
type Client interface {
	Version() v1.VersionServiceClient
	Cluster() v1.ClusterServiceClient
	IP() v1.IPServiceClient
	Volume() v1.VolumeServiceClient
	Stripe() v1.StripeServiceClient
	Close() error
}

// GRPCClient is a Client implementation with grpc transport.
type GRPCClient struct {
	conn *grpc.ClientConn
	log  *zap.SugaredLogger
}

// Close the underlying connection
func (c GRPCClient) Close() error {
	return c.conn.Close()
}

// Version client
func (c GRPCClient) Version() v1.VersionServiceClient {
	return v1.NewVersionServiceClient(c.conn)
}

// Version client
func (c GRPCClient) Cluster() v1.ClusterServiceClient {
	return v1.NewClusterServiceClient(c.conn)
}

// Version client
func (c GRPCClient) IP() v1.IPServiceClient {
	return v1.NewIPServiceClient(c.conn)
}

// Version client
func (c GRPCClient) Volume() v1.VolumeServiceClient {
	return v1.NewVolumeServiceClient(c.conn)
}

// Version client
func (c GRPCClient) Stripe() v1.StripeServiceClient {
	return v1.NewStripeServiceClient(c.conn)
}

func Dial(ctx context.Context, config DialConfig) (Client, error) {
	log := config.Log

	ua := defaultUserAgent
	if config.UserAgent != "" {
		ua = config.UserAgent
	}

	log.Infow("connecting...",
		"client", ua,
		"endpoint", config.Endpoint,
	)

	res := &GRPCClient{
		log: log,
	}

	zapOpts := []grpc_zap.Option{
		grpc_zap.WithLevels(grpcToZapLevel),
	}
	interceptors := []grpc.UnaryClientInterceptor{
		grpc_zap.UnaryClientInterceptor(log.Desugar(), zapOpts...),
		grpc_zap.PayloadUnaryClientInterceptor(log.Desugar(),
			func(context.Context, string) bool { return true },
		),
	}

	// these are broadly in line with the expected server SLOs:
	kal := keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	}

	dialBackoffConfig := backoff.Config{
		BaseDelay:  1.0 * time.Second,
		Multiplier: 1.2,
		Jitter:     0.1,
		MaxDelay:   7 * time.Second,
	}
	cp := grpc.ConnectParams{
		Backoff:           dialBackoffConfig,
		MinConnectTimeout: 6 * time.Second,
	}

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithDisableRetry(),
		grpc.WithUserAgent(ua),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...)),
		grpc.WithKeepaliveParams(kal),
		grpc.WithConnectParams(cp),
		grpc.WithPerRPCCredentials(tokenAuth{
			token: config.Token,
		}),
	}
	// Configure tls ca certificate based auth if credentials are given
	switch config.Scheme {
	case GRPC:
		log.Infof("connecting insecurely")
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	case GRPCS:
		log.Infof("connecting securely")
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	default:
		return nil, fmt.Errorf("unsupported scheme:%v", config.Scheme)
	}

	var err error
	res.conn, err = grpc.DialContext(
		ctx,
		config.Endpoint,
		opts...,
	)
	if err != nil {
		log.Errorw("failed to connect", "endpoint", config.Endpoint, "error", err.Error())
		return nil, err
	}

	return res, nil
}

type tokenAuth struct {
	token string
}

func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": "Bearer " + t.token,
	}, nil
}

func (tokenAuth) RequireTransportSecurity() bool {
	return true
}

func grpcToZapLevel(code codes.Code) zapcore.Level {
	switch code {
	case codes.OK,
		codes.Canceled,
		codes.DeadlineExceeded,
		codes.NotFound,
		codes.Unavailable:
		return zapcore.InfoLevel
	case codes.Aborted,
		codes.AlreadyExists,
		codes.FailedPrecondition,
		codes.InvalidArgument,
		codes.OutOfRange,
		codes.PermissionDenied,
		codes.ResourceExhausted,
		codes.Unauthenticated:
		return zapcore.WarnLevel
	case codes.DataLoss,
		codes.Internal,
		codes.Unimplemented,
		codes.Unknown:
		return zapcore.ErrorLevel
	default:
		return zapcore.ErrorLevel
	}
}
