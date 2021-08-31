package main

import (
	"context"
	pb "github.com/maxlcoder/grpc-example/grpc-deadline/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"time"
)

const (
	PORT = ":50051"
)

type SearchService struct {
	pb.UnimplementedSearchServiceServer
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	for i := 0; i< 5; i++ {
		if ctx.Err() == context.Canceled {
			return nil, status.Errorf(
				codes.Canceled, "SearchService.Search canceled")
		}
		time.Sleep(1 * time.Second)
	}
	return &pb.SearchResponse{
		Response: "Search Key: " + r.GetRequest(),
	}, nil
}

func main()  {
	cred, err := credentials.NewServerTLSFromFile("../cert/server.pem", "../cert/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile fail: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(cred))
	pb.RegisterSearchServiceServer(grpcServer, &SearchService{})

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("net listen tcp failed: %v", err)
	}
	grpcServer.Serve(lis)
}
