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
	succesfulRequests int
}

func(s *server) SubmitResults(ctx context.Context, req *pb.BatchPayload) (empty *emptypb.Empty, err error){
	// Lock mutex to prevent race conditions
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Loop over every result in the payload
	for _, res:= range req.Results{
		// Increment the request counter
		s.totalRequests += 1
		s.latency += res.GetLatency()

		// Check if the request was succesful
		if(res.GetStatus() == 200){
			// Increase succesful requests counter
			s.succesfulRequests += 1
		}

	}
	// Calculate average latency
	var currentAverage int64 = 0
	if s.totalRequests > 0 {
		currentAverage = s.latency / int64(s.totalRequests)
	}

	// Log responses
	fmt.Printf("Received result, total requests: %d. Succesful requests: %d. Current Average Latency: %d\n", s.totalRequests, s.succesfulRequests, currentAverage)
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