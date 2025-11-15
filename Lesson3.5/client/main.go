package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "main.go/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // gRPC Consul resolver
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. Dial 连接 Consul 服务（忽略健康检查）
	conn, err := grpc.DialContext(
		ctx,
		"consul://127.0.0.1:8500/greeter?healthy=false", // healthy=false
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // round-robin
	)
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	// 3. 创建客户端
	client := pb.NewGreeterClient(conn)

	// 4. 循环调用 6 次，观察轮询效果
	for i := 0; i < 6; i++ {
		resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "R Mel"})
		if err != nil {
			log.Fatalf("SayHello failed: %v", err)
		}
		fmt.Println(resp.Reply)
		time.Sleep(time.Second)
	}
}
