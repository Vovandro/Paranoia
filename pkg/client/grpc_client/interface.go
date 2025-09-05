package grpc_client

import "google.golang.org/grpc"

type IGrpcClient interface {
	// GetClient returns the underlying gRPC client connection.
	GetClient() *grpc.ClientConn
}
