package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"main/Stream/proto"

	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ----------- 核心服务实现 -----------
type chatServer struct {
	proto.UnimplementedChatServiceServer
}

// 双向流式 RPC
func (s *chatServer) Chat(stream proto.ChatService_ChatServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil // 客户端关闭发送流
		}
		if err != nil {
			return err
		}

		reply := fmt.Sprintf("AI 回复 (%s): %s", msg.GetName(), magic(msg.GetContent()))
		if err := stream.Send(&proto.ChatMessage{
			Name:    "AI",
			Content: reply,
		}); err != nil {
			return err
		}
	}
}

// “魔法”文本处理函数
func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "你", "我")
	return s
}

// ----------- 拦截器实现 -----------

// 模拟 token 验证
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == "secret-token"
}

// 一元拦截器
func unaryServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	if !valid(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	fmt.Printf("[Unary Interceptor] 方法: %s\n", info.FullMethod)
	resp, err := handler(ctx, req)
	fmt.Printf("[Unary Interceptor] 用时: %v, err: %v\n", time.Since(start), err)
	return resp, err
}

// 自定义流包装器
type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	fmt.Printf("[Stream Interceptor] 收到消息类型: %T, 时间: %v\n", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	fmt.Printf("[Stream Interceptor] 发送消息类型: %T, 时间: %v\n", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

// 流式拦截器
func streamServerInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	md, _ := metadata.FromIncomingContext(ss.Context())
	if !valid(md["authorization"]) {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}
	fmt.Printf("[Stream Interceptor] 方法: %s\n", info.FullMethod)
	err := handler(srv, &wrappedStream{ss})
	fmt.Printf("[Stream Interceptor] 结束, err: %v\n\n", err)
	return err
}

func main() {
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}

	s := grpc.NewServer(
		//注册一元拦截器
		grpc.UnaryInterceptor(unaryServerInterceptor),
		//注册流拦截器
		grpc.StreamInterceptor(streamServerInterceptor),
	)
	//调用流拦截器
	proto.RegisterChatServiceServer(s, &chatServer{})

	fmt.Println("✅ Chat gRPC 服务启动中，端口 :8972 ...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
