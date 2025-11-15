package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"main/Stream/proto"

	"os"
	"strings"
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
	md := metadata.Pairs("authorization", "Bearer secret-token")
	ctx = metadata.NewOutgoingContext(ctx, md)
	start := time.Now()
	fmt.Printf("[Unary Client Interceptor] 调用: %s\n", method)
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Printf("[Unary Client Interceptor] 耗时: %v, 错误: %v\n", time.Since(start), err)
	return err
}

// 包装流对象
type wrappedStream struct {
	grpc.ClientStream
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	fmt.Printf("[Stream Client Interceptor] 发送消息: %T, 时间: %v\n", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	fmt.Printf("[Stream Client Interceptor] 接收消息: %T, 时间: %v\n", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

// 客户端流拦截器
func streamClientInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	md := metadata.Pairs("authorization", "Bearer secret-token")
	ctx = metadata.NewOutgoingContext(ctx, md)
	fmt.Printf("[Stream Client Interceptor] 建立流: %s\n", method)
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return &wrappedStream{s}, nil
}

func main() {
	conn, err := grpc.Dial(
		"127.0.0.1:8972",
		grpc.WithInsecure(),
		//注册一元拦截器
		grpc.WithUnaryInterceptor(unaryClientInterceptor),
		//注册流拦截器
		grpc.WithStreamInterceptor(streamClientInterceptor),
	)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := proto.NewChatServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	//创建流，调用流拦截器
	stream, err := client.Chat(ctx)
	if err != nil {
		log.Fatalf("创建流失败: %v", err)
	}

	// 启动接收协程
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalf("接收失败: %v", err)
			}
			fmt.Printf("\n%s: %s\n", in.GetName(), in.GetContent())
		}
	}()

	// 从命令行读取输入
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("输入聊天内容 (输入 QUIT 退出):")
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if strings.ToUpper(text) == "QUIT" {
			break
		}
		if err := stream.Send(&proto.ChatMessage{
			Name:    "七米",
			Content: text,
		}); err != nil {
			log.Fatalf("发送失败: %v", err)
		}
	}
	stream.CloseSend()
}
