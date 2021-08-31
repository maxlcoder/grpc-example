package main

import (
	"context"
	"github.com/maxlcoder/grpc-example/grpc-zipkin/pkg/gtls"
	pb "github.com/maxlcoder/grpc-example/grpc-zipkin/proto"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"runtime/debug"
)

const (
	PORT = ":50051"
	ZIPKIN_SPAN_REPORT_URL = "http://localhost:9411/api/v2/spans"
)

type SearchService struct {
	pb.UnimplementedSearchServiceServer
}

func (s *SearchService) Search(context context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{
		Response: "Search Key: " + r.GetRequest(),
	}, nil
}

// logging 拦截器
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

// recover 拦截器
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()
	return handler(ctx, req)
}

func main()  {

	reporter:= httpreporter.NewReporter(ZIPKIN_SPAN_REPORT_URL)
	defer reporter.Close()

	endpoint, err := zipkin.NewEndpoint("myService", "localhost:50051")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// 初始化 tracer
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	certFile := "../cert/server.pem"
	keyFile := "../cert/server.key"
	CaFile := "../cert/ca.pem"
	tlsServer := gtls.Server{
		KeyFile: keyFile,
		CertFile: certFile,
		CaFile: CaFile,
	}

	c, err := tlsServer.GetCredentialsByCA()
	if err != nil {
		log.Fatalf("tlsServer.GetCredentialsByCA err: %v", err)
	}

	// opt
	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpc.ChainUnaryInterceptor(
			RecoveryInterceptor,
			LoggingInterceptor,
		),
		grpc.StatsHandler(zipkingrpc.NewServerHandler(tracer)),
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(grpcServer, &SearchService{})

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("net listen tcp failed: %v", err)
	}
	grpcServer.Serve(lis)
}