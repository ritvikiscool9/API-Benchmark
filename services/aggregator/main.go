package main

import (
	"APIbenchmark/proto/pb"
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// gRPC struct
type server struct{
	pb.UnimplementedAggregatorServer
	mutex sync.Mutex
	totalRequests int
	latency int64
}

func(s *server) SubmitResults(ctx context.Context, req *pb.Result) (empty *emptypb.Empty, err error){
	// Lock mutex to prevent race conditions
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Increment the request counter
	s.totalRequests += 1
	s.latency += req.GetLatency()

	// Log responses
	fmt.Printf("Received result, total requests so far: %d. Latency: %d", s.totalRequests, s.latency)

	return &emptypb.Empty{}, nil
}

func main(){
	// Open gRPC network
	ln, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Network connection failed")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAggregatorServer(grpcServer, &server{})

	log.Printf("Aggregator server listening on port 50051")

	grpcServer.Serve(ln)
}