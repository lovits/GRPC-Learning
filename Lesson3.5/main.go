package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "main.go/proto" // 生成的 proto Go 包

	"github.com/hashicorp/consul/api" // Consul 客户端
	"google.golang.org/grpc"          // gRPC
)

// ---------------- gRPC 服务实现 ----------------
type server struct {
	pb.UnimplementedGreeterServer
	id string // 实例 ID，用于区分多实例
}

// SayHello 实现 Greeter 服务
func (s *server) SayHello(_ context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Reply: fmt.Sprintf("Hello, %s! From instance %s", req.Name, s.id),
	}, nil
}

// ---------------- Consul 封装 ----------------
type consul struct {
	client *api.Client
}

// NewConsul 连接到本地 Consul
func NewConsul(addr string) (*consul, error) {
	c, err := api.NewClient(&api.Config{Address: addr})
	if err != nil {
		return nil, err
	}
	return &consul{c}, nil
}

// RegisterService 注册服务到 Consul（不使用健康检查）
func (c *consul) RegisterService(serviceID, serviceName, ip string, port int) error {
	srv := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: ip,
		Port:    port,
		// Check: nil  // 不使用健康检查，避免 All checks failing
	}
	return c.client.Agent().ServiceRegister(srv)
}

// Deregister 注销服务
func (c *consul) Deregister(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}

// ---------------- 启动单个 gRPC 实例 ----------------
func StartGRPCInstance(id string, port int) {
	ip := "127.0.0.1"
	serviceName := "greeter"
	serviceID := fmt.Sprintf("%s-%s-%d", serviceName, id, port)

	// 1. 连接 Consul
	consulClient, err := NewConsul("127.0.0.1:8500")
	if err != nil {
		log.Fatalf("Consul client failed: %v", err)
	}

	// 2. 注册服务
	err = consulClient.RegisterService(serviceID, serviceName, ip, port)
	if err != nil {
		log.Fatalf("Register service failed: %v", err)
	}
	log.Printf("Service %s registered on port %d", serviceID, port)

	// 3. 启动 gRPC 服务器
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Listen failed: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{id: id})

	// 4. 启动协程监听服务
	go func() {
		log.Printf("gRPC server instance %s running on %s:%d", id, ip, port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Serve failed: %v", err)
		}
	}()

	// 5. 优雅退出处理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Deregistering %s", serviceID)
	consulClient.Deregister(serviceID)
	s.Stop()
}

// ---------------- main ----------------
func main() {
	// 启动 3 个 gRPC 实例
	go StartGRPCInstance("A", 50051)
	go StartGRPCInstance("B", 50052)
	go StartGRPCInstance("C", 50053)

	select {} // 阻塞主 goroutine，保证服务一直运行
}
