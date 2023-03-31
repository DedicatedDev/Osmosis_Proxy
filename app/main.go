package main

import (
	"context"
	"log"
	"net"

	tendermintv1beta1 "github.com/dedicatedDev/txbroker/pkg/cosmos/base/tendermint/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	tendermintv1beta1.UnimplementedServiceServer
	client tendermintv1beta1.ServiceClient
}

// Implements the GetNodeInfo RPC method from the tendermintv1beta1.ServiceServer interface
func (s *server) GetNodeInfo(ctx context.Context, req *tendermintv1beta1.GetNodeInfoRequest) (*tendermintv1beta1.GetNodeInfoResponse, error) {
	return s.client.GetNodeInfo(ctx, req)
}

// Implements the GetLatestBlock RPC method from the tendermintv1beta1.ServiceServer interface
func (s *server) GetLatestBlock(ctx context.Context, req *tendermintv1beta1.GetLatestBlockRequest) (*tendermintv1beta1.GetLatestBlockResponse, error) {
	return s.client.GetLatestBlock(ctx, req)
}

func main() {
	// Connect to the Osmosis gRPC server
	conn, err := grpc.Dial("grpc.osmosis.zone:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to osmosis rpc: %v", err)
	}
	defer conn.Close()

	// Create a client for the tendermintv1beta1.Service
	osmosisClient := tendermintv1beta1.NewServiceClient(conn)

	// Create a gRPC server and register the server struct with it
	grpcServer := grpc.NewServer()

	// Register the reflection service with the server
	reflection.Register(grpcServer)

	// Register the server struct with the server
	tendermintv1beta1.RegisterServiceServer(grpcServer, &server{client: osmosisClient})

	// Listen for incoming connections on localhost:9000
	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Start serving gRPC requests
	log.Println("Server is listening on localhost:9000")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
