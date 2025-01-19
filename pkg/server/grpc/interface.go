package grpc

import (
	"google.golang.org/grpc"
)

// IGrpc defines the interface for gRPC server operations
type IGrpc interface {
	// RegisterService registers a service with the gRPC server
	RegisterService(desc *grpc.ServiceDesc, impl any)
}
