package grpc_client

import (
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Name   string
	Config GrpcClientConfig
	app    interfaces.IEngine
	client *grpc.ClientConn
}

type GrpcClientConfig struct {
	Url string `yaml:"url"`
}

func NewGrpcClient(name string, cfg GrpcClientConfig) *GrpcClient {
	return &GrpcClient{
		Name:   name,
		Config: cfg,
	}
}

func (t *GrpcClient) Init(app interfaces.IEngine) error {
	var err error
	t.app = app

	t.client, err = grpc.NewClient(
		t.Config.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	return err
}

func (t *GrpcClient) Stop() error {
	return t.client.Close()
}

func (t *GrpcClient) String() string {
	return t.Name
}

func (t *GrpcClient) GetClient() *grpc.ClientConn {
	return t.client
}
