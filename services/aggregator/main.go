package main

import (
	"APIbenchmark/proto/pb"
	"context"
	"fmt"
	"log"
	"net"
	"slices"
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
	latencies []int64
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
		s.latencies = append(s.latencies, res.GetLatency())

		// Check if the request was succesful
		if(res.GetStatus() == 200){
			// Increase succesful requests counter
			s.succesfulRequests += 1
		}

	}
	// Calculate average latency
	var currentAverage int64 
    if s.totalRequests > 0 {
        currentAverage = s.latency / int64(s.totalRequests)
    }

	// Calculate P95 and P99 latencies
	var P95Latency int64
	var P99Latency int64
	if (len(s.latencies)) > 0 {
		slices.Sort(s.latencies)

		P95Index := int(float64(len(s.latencies)) * 0.95)
		P99Index := int(float64(len(s.latencies)) * 0.99)

		P95Latency = s.latencies[P95Index]
		P99Latency = s.latencies[P99Index]
	}

	// Log responses
	fmt.Printf("Received result, total requests: %d. Successful requests: %d. Average Latency: %dms. P95: %dms. P99: %dms\n", 
        s.totalRequests, 
        s.succesfulRequests, 
        currentAverage, 
        P95Latency, 
        P99Latency)
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