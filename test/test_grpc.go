package main

import (
	"context"
	"log"
	"reflect"
	"time"

	tendermintv1beta1 "github.com/dedicatedDev/txproxy/pkg/cosmos/base/tendermint/v1beta1"
	"google.golang.org/grpc"
)

type TestCase struct {
	name    string
	request interface{}
}

func main() {
	// Connect to the local gRPC server
	localConn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to the local gRPC server: %v", err)
	}
	defer localConn.Close()
	localClient := tendermintv1beta1.NewServiceClient(localConn)

	// Connect to the Osmosis public RPC endpoint
	osmosisConn, err := grpc.Dial("grpc.osmosis.zone:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to the Osmosis RPC: %v", err)
	}
	defer osmosisConn.Close()
	osmosisClient := tendermintv1beta1.NewServiceClient(osmosisConn)

	testCases := []TestCase{
		{
			"Get Node Info",
			&tendermintv1beta1.GetNodeInfoRequest{},
		},
		{
			"Get Latest Block",
			&tendermintv1beta1.GetLatestBlockRequest{},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i, tc := range testCases {
		var localResp, osmosisResp interface{}
		var err error

		switch req := tc.request.(type) {
		case *tendermintv1beta1.GetNodeInfoRequest:
			lresp, err := localClient.GetNodeInfo(ctx, req)
			lresp.DefaultNodeInfo.DefaultNodeID = ""
			localResp = lresp
			if err != nil {
				log.Printf("Test %d: Local gRPC server request failed: %v", i, err)
				continue
			}
			oResp, err := osmosisClient.GetNodeInfo(ctx, req)
			if err != nil {
				log.Printf("Test %d: Local gRPC server request failed: %v", i, err)
				continue
			}
			oResp.DefaultNodeInfo.DefaultNodeID = ""
			osmosisResp = oResp
		case *tendermintv1beta1.GetLatestBlockRequest:
			localResp, err = localClient.GetLatestBlock(ctx, req)
			if err != nil {
				log.Printf("Test %d: Local gRPC server request failed: %v", i, err)
				continue
			}
			osmosisResp, err = osmosisClient.GetLatestBlock(ctx, req)
		default:
			log.Printf("Test %d: Unsupported request type", i)
			continue
		}

		if err != nil {
			log.Printf("Test %d: Osmosis RPC request failed: %v", i, err)
			continue
		}

		// Compare the responses
		isEqual := reflect.DeepEqual(localResp, osmosisResp)
		if isEqual {
			log.Printf("Test %d: PASSED", i)
		} else {
			log.Printf("Test %d: FAILED - Local response: %v, Osmosis response: %v", i, localResp, osmosisResp)
		}
	}
}
