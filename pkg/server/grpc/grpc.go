package grpc

import (
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Grpc struct {
	name   string
	config Config

	server *grpc.Server
}

type Config struct {
	Port string `yaml:"port"`
}

func NewGrpc(name string) *Grpc {
	return &Grpc{
		name: name,
	}
}

func (t *Grpc) Init(cfg map[string]interface{}) error {
	if _, ok := cfg["middlewares"]; ok {
		delete(cfg, "middlewares")
	}

	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Port == "" {
		t.config.Port = "9090"
	}

	t.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	return nil

}

func (t *Grpc) Start() error {

	listenErr := make(chan error, 1)

	go func() {
		l, err := net.Listen("tcp", ":"+t.config.Port)
		if err != nil {
			listenErr <- err
			return
		}

		listenErr <- t.server.Serve(l)
	}()

	select {
	case err := <-listenErr:
		return err

	case <-time.After(time.Second):
		// pass
	}

	return nil
}

func (t *Grpc) Stop() error {
	t.server.GracefulStop()

	time.Sleep(time.Second)

	return nil
}

func (t *Grpc) Name() string {
	return t.name
}

func (t *Grpc) Type() string {
	return "server"
}

func (t *Grpc) RegisterService(desc *grpc.ServiceDesc, impl any) {
	t.server.RegisterService(desc, impl)
}
