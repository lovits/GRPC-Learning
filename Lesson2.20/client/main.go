package main

import (
	"context"
	"flag"
	"log"
	"time"

	"lesson220/proto/helloworld" // ⚠️ 改成你项目 go.mod 中的 module 路径

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var defaultName = "hello"

var (
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// 1️⃣ 连接 gRPC 服务
	conn, err := grpc.Dial("127.0.0.1:8091",
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 无 TLS
	)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 2️⃣ 创建客户端
	client := helloworld.NewGreeterClient(conn)

	// 3️⃣ 调用远程方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &helloworld.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("调用失败: %v", err)
	}

	log.Printf("响应结果: %s", resp.Message)
}
