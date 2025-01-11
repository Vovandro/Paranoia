package grpc_client

import (
	"errors"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	name   string
	config Config
	client *grpc.ClientConn
}

type Config struct {
	Url string `yaml:"url"`
}

func NewGrpcClient(name string) *GrpcClient {
	return &GrpcClient{
		name: name,
	}
}

func (t *GrpcClient) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Url == "" {
		return errors.New("url is required")
	}

	t.client, err = grpc.NewClient(
		t.config.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	return err
}

func (t *GrpcClient) Stop() error {
	return t.client.Close()
}

func (t *GrpcClient) Name() string {
	return t.name
}

func (t *GrpcClient) Type() string {
	return "client"
}

func (t *GrpcClient) GetClient() *grpc.ClientConn {
	return t.client
}
