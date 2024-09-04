package server

import (
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Grpc struct {
	Name   string
	Config GrpcConfig

	app    interfaces.IEngine
	server *grpc.Server
}

type GrpcConfig struct {
	Port string `yaml:"port"`
}

func NewGrpc(name string, cfg GrpcConfig) *Grpc {
	return &Grpc{
		Name:   name,
		Config: cfg,
	}
}

func (t *Grpc) Init(app interfaces.IEngine) error {
	t.app = app

	t.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	return nil

}

func (t *Grpc) Start() error {

	listenErr := make(chan error, 1)

	go func() {
		l, err := net.Listen("tcp", ":"+t.Config.Port)
		if err != nil {
			listenErr <- err
			return
		}

		listenErr <- t.server.Serve(l)
	}()

	select {
	case err := <-listenErr:
		t.app.GetLogger().Error(err)
		return err

	case <-time.After(time.Second):
		// pass
	}

	return nil
}

func (t *Grpc) Stop() error {
	t.server.GracefulStop()

	t.app.GetLogger().Info("grpc server gracefully stopped.")
	time.Sleep(time.Second)

	return nil
}

func (t *Grpc) String() string {
	return t.Name
}

func (t *Grpc) RegisterService(desc *grpc.ServiceDesc, impl any) {
	t.server.RegisterService(desc, impl)
}
