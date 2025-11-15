package main

import (
	"context"
	"lesson220/proto/helloworld"
	"log"

	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 1️⃣ 实现服务
type greeterServer struct {
	helloworld.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: req.Name + " world"}, nil
}

// 2️⃣ 启动服务器：同时支持 gRPC + HTTP
func main() {
	addr := ":8091"
	lis, err := net.Listen("tcp", addr)
	if err != nil {	
		log.Fatalf("Failed to listen: %v", err)
	}

	// gRPC Server
	grpcServer := grpc.NewServer()
	helloworld.RegisterGreeterServer(grpcServer, &greeterServer{})

	// HTTP Gateway
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = helloworld.RegisterGreeterHandlerFromEndpoint(context.Background(), gwmux, "127.0.0.1"+addr, opts)
	if err != nil {
		log.Fatalf("Failed to register handler: %v", err)
	}

	// HTTP mux
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	// 合并 gRPC + HTTP
	server := &http.Server{
		Addr:    addr,
		Handler: grpcHandlerFunc(grpcServer, mux),
	}

	log.Printf("🚀 Serving gRPC + HTTP on %s", addr)
	log.Fatal(server.Serve(lis))
}

// 判断请求类型（gRPC or HTTP）
func grpcHandlerFunc(grpcServer *grpc.Server, other http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			other.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
