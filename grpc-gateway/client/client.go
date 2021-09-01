package main

import (
	"context"
	pb "github.com/maxlcoder/grpc-example/grpc-gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
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

	client := pb.NewHelloServiceClient(conn)
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{
		Request: "grpc",
	})
	if err != nil {
		log.Fatalf("client.Search failed: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())

}
