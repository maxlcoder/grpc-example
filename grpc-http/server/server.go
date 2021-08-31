package main

import (
	"context"
	"github.com/maxlcoder/grpc-example/grpc-http/pkg/gtls"
	pb "github.com/maxlcoder/grpc-example/grpc-http/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
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

	grpcServer := grpc.NewServer(grpc.Creds(c))
	pb.RegisterSearchServiceServer(grpcServer, &SearchService{})

	mux := GetHTTPServeMux()

	http.ListenAndServeTLS(PORT,
		certFile,
		keyFile,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}
			return
		}),
	)

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("net listen tcp failed: %v", err)
	}
	grpcServer.Serve(lis)
}

func GetHTTPServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	return mux
}