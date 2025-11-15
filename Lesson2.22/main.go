package main

import (
	"context"
	"fmt"
	"lesson222/proto"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接 MySQL
	db, err := NewDB("root:uwycuge@tcp(localhost:3306)/bookstore_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("connect to db fail %v", err)
		return
	}

	ser := server{bs: &bookstore{db: db}}

	// gRPC TCP 服务
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("fail to create lis")
	}
	s := grpc.NewServer()
	proto.RegisterBookstoreServer(s, &ser)
	go func() {
		fmt.Println(s.Serve(lis)) // 启动 gRPC 服务
	}()

	// gRPC-Gateway: HTTP/JSON 转发
	conn, err := grpc.DialContext(
		context.Background(),
		"127.0.0.1:8972",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = proto.RegisterBookstoreHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// HTTP 服务监听 8090，访问 REST API
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}
	fmt.Println("grpc gateway qidong")
	gwServer.ListenAndServe()
}
