package server

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/client/grpc-client"
	"gitlab.com/devpro_studio/Paranoia/framework"
	"gitlab.com/devpro_studio/Paranoia/logger"
	"gitlab.com/devpro_studio/Paranoia/server/example"
	"testing"
)

type server struct {
	example.UnimplementedExampleServer
}

func (t *server) Do(c context.Context, r *example.Request) (*example.Response, error) {
	resp := &example.Response{}

	resp.Message = r.Message

	return resp, nil
}

func TestGrpc_RegisterService(t1 *testing.T) {
	t1.Run("base test", func(t1 *testing.T) {
		app := framework.New("test", nil, &logger.Mock{})

		s := NewGrpc("test", GrpcConfig{Port: "8091"})
		s.Init(app)

		se := server{}

		s.RegisterService(&example.Example_ServiceDesc, &se)

		s.Start()
		defer s.Stop()

		c := grpc_client.NewGrpcClient("test", grpc_client.GrpcClientConfig{Url: "localhost:8091"})
		c.Init(app)
		defer c.Stop()

		cs := example.NewExampleClient(c.GetClient())

		resp, err := cs.Do(context.Background(), &example.Request{Message: "t1"})
		if err != nil {
			t1.Fatal(err)
		}

		if resp.Message != "t1" {
			t1.Fatal("resp.Message != \"t1\"")
		}
	})
}
