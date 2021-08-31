package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	pb "github.com/maxlcoder/grpc-example/grpc-interceptor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"net"
	"runtime/debug"
)

const (
	PORT = ":50051"
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
	cert, err := tls.LoadX509KeyPair("../cert/server.pem", "../cert/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile fail: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../cert/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile fail: $v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM failed")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: certPool,
	})

	// opt
	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpc.ChainUnaryInterceptor(
			RecoveryInterceptor,
			LoggingInterceptor,
		),
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(grpcServer, &SearchService{})

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("net listen tcp failed: %v", err)
	}
	grpcServer.Serve(lis)
}
