package main

import (
	"context"
	"github.com/maxlcoder/grpc-example/grpc-zipkin/pkg/gtls"
	pb "github.com/maxlcoder/grpc-example/grpc-zipkin/proto"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"
	"log"
)

const (
	PORT = ":50051"
	ZIPKIN_SPAN_REPORT_URL = "http://localhost:9411/api/v2/spans"
)

func main()  {
	reporter := httpreporter.NewReporter(ZIPKIN_SPAN_REPORT_URL)
	defer reporter.Close()

	endpoint, err := zipkin.NewEndpoint("myClient", "localhost:50051")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// 初始化 tracer
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

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
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(c),
		grpc.WithStatsHandler(zipkingrpc.NewClientHandler(tracer)),
	}
	conn, err := grpc.Dial(PORT, opts...)
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
