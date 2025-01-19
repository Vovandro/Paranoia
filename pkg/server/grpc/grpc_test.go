package grpc

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/pkg/server/grpc/example"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		s := New("test")
		s.Init(map[string]interface{}{
			"port": "8091",
		})

		se := server{}

		s.RegisterService(&example.Example_ServiceDesc, &se)

		s.Start()
		defer s.Stop()

		client, err := grpc.NewClient(
			"localhost:8091",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		cs := example.NewExampleClient(client)

		resp, err := cs.Do(context.Background(), &example.Request{Message: "t1"})
		if err != nil {
			t1.Fatal(err)
		}

		if resp.Message != "t1" {
			t1.Fatal("resp.Message != \"t1\"")
		}
	})
}
