package main

import (
	"context"
	"github.com/maxlcoder/grpc-example/grpc-http/pkg/gtls"
	pb "github.com/maxlcoder/grpc-example/grpc-http/proto"
	"google.golang.org/grpc"
	"log"
)

const (
	PORT = ":50051"
)

func main()  {
	serverName := "go-grpc"
	certFile := "../cert/client.pem"
	keyFile := "../cert/client.key"
	CaFile := "../cert/ca.pem"
	tlsClient := gtls.Client{
		CertFile: certFile,
		KeyFile: keyFile,
		CaFile: CaFile,
		ServerName: serverName,
	}

	c, err := tlsClient.GetCredentialsByCA()
	if err != nil {
		log.Fatalf("tlsClient.GetCredentialsByCA fail: %v", err)
	}

	conn, err := grpc.Dial(PORT, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "Ca TLS",
	})
	if err != nil {
		log.Fatalf("client.Search failed: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())

}
