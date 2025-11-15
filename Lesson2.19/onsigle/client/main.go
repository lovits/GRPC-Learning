package main

import (
	"context"
	"fmt"
	"log"
	"main/onsigle/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 客户端一元拦截器
func unaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	fmt.Printf("[Client Interceptor] 调用 RPC 方法: %s\n", method)

	// 给每次请求加上 token
	md := metadata.Pairs("authorization", "Bearer secret-token")
	ctx = metadata.NewOutgoingContext(ctx, md)

	// 执行实际的 RPC 调用
	err := invoker(ctx, method, req, reply, cc, opts...)

	fmt.Printf("[Client Interceptor] 调用结束，用时: %v，错误: %v\n\n", time.Since(start), err)
	return err
}

func main() {
	conn, err := grpc.Dial(
		"127.0.0.1:8972",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(unaryClientInterceptor),
	)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := proto.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &proto.HelloRequest{Name: "七米"})
	if err != nil {
		log.Fatalf("调用失败: %v", err)
	}

	fmt.Printf("✅ 服务端回复: %s\n", resp.GetReply())
}
