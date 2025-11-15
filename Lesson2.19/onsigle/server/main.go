package main

import (
	"context"
	"fmt"
	"log"
	"main/onsigle/proto"

	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// 实现 Greeter 服务
type greeterServer struct {
	proto.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	reply := fmt.Sprintf("你好，%s！", req.GetName())
	return &proto.HelloResponse{Reply: reply}, nil
}

// 模拟 token 验证
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == "secret-token"
}

// 服务端一元拦截器
func unaryServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// 读取 metadata 并验证 token
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	if !valid(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	fmt.Printf("[Server Interceptor] 调用方法: %s\n", info.FullMethod)

	// 执行真正的 RPC 逻辑
	resp, err := handler(ctx, req)

	fmt.Printf("[Server Interceptor] 调用结束，用时: %v，错误: %v\n\n", time.Since(start), err)
	return resp, err
}

func main() {
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(unaryServerInterceptor),
	)

	proto.RegisterGreeterServer(s, &greeterServer{})

	fmt.Println("✅ gRPC 服务已启动，监听端口 :8972 ...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
