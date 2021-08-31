package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	pb "github.com/maxlcoder/grpc-example/grpc-interceptor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

const (
	PORT = ":50051"
)

func main()  {
	cert, err := tls.LoadX509KeyPair("../cert/client.pem", "../cert/client.key")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile failed: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../cert/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName: "go-grpc",
		RootCAs: certPool,
	})

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
