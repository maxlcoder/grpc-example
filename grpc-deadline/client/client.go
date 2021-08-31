package main

import (
	"context"
	pb "github.com/maxlcoder/grpc-example/grpc-deadline/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

const (
	PORT = ":50051"
)

func main()  {
	cred, err := credentials.NewClientTLSFromFile("../cert/server.pem", "go-grpc")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile failed: %v", err)
	}
	conn, err := grpc.Dial(PORT, grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(5 * time.Second)))
	defer cancel()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(ctx, &pb.SearchRequest{
		Request: "TLS",
	})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok  {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Fatalln("client.Search err: deadline")
			}
		}
		log.Fatalf("client.Search failed: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())
}
