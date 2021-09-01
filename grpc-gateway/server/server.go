package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/maxlcoder/grpc-example/grpc-gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
)

const (
	PORT = ":50051"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Response: "Hello " + r.GetRequest(),
	}, nil
}

func main()  {
	// 这里是用 tls 加密通信，没有使用基于 ca，这样 服务端和客户端共用一套密钥
	certFile := "../cert/server.pem"
	keyFile := "../cert/server.key"
	//certName := "go-grpc"

	// 生成服务端 tls 证书
	serverCred, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}
	// 创建 grpc 服务实例
	grpcServer := grpc.NewServer(grpc.Creds(serverCred))
	// 将服务注册到对应的服务实例
	pb.RegisterHelloServiceServer(grpcServer, &server{})

	// 生成客户端 tls 证书
	clientCred, err := credentials.NewClientTLSFromFile("../cert/server.pem", "go-grpc")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("net listen tcp failed: %v", err)
	}
	// 开启 grpc 服务，这里将监听 50051 端口，也就是这个端口的请求均认为是 grpc 的调用
	go func() {
		log.Println("Serving gRPC on http://0.0.0.0:50051")
		grpcServer.Serve(lis)
	}()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(clientCred),
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterHelloServiceHandlerFromEndpoint(context.Background(), gwmux, PORT, opts)
	if err != nil {
		log.Fatalf("pb.RegisterHelloServiceHandler fail: %v", err)
	}

	gwServer := &http.Server{
		Addr: ":8090",
		Handler: gwmux,
	}
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServeTLS(certFile, keyFile))

}
