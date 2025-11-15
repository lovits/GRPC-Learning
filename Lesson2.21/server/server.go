package main

import (
	"context"
	"fmt"
	"lesson221/proto"
	"log"

	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// 内存存储
type protoServer struct {
	proto.UnimplementedBookstoreServer
	shelves map[int64]*proto.Shelf
	nextID  int64
	mu      sync.Mutex
}

func newServer() *protoServer {
	return &protoServer{
		shelves: make(map[int64]*proto.Shelf),
		nextID:  1,
	}
}

// ListShelves
func (s *protoServer) ListShelves(ctx context.Context, _ *proto.Empty) (*proto.ListShelvesResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	resp := &proto.ListShelvesResponse{}
	for _, sh := range s.shelves {
		resp.Shelves = append(resp.Shelves, sh)
	}
	return resp, nil
}

// GetShelf
func (s *protoServer) GetShelf(ctx context.Context, req *proto.GetShelfRequest) (*proto.Shelf, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sh, ok := s.shelves[req.Id]; ok {
		return sh, nil
	}
	return nil, fmt.Errorf("shelf not found")
}

// CreateShelf
func (s *protoServer) CreateShelf(ctx context.Context, shelf *proto.Shelf) (*proto.Shelf, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	shelf.Id = s.nextID
	s.nextID++
	s.shelves[shelf.Id] = shelf
	return shelf, nil
}

// UpdateShelf
func (s *protoServer) UpdateShelf(ctx context.Context, shelf *proto.Shelf) (*proto.Shelf, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.shelves[shelf.Id]; !ok {
		return nil, fmt.Errorf("shelf not found")
	}
	s.shelves[shelf.Id] = shelf
	return shelf, nil
}

// DeleteShelf
func (s *protoServer) DeleteShelf(ctx context.Context, req *proto.GetShelfRequest) (*proto.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.shelves[req.Id]; !ok {
		return nil, fmt.Errorf("shelf not found")
	}
	delete(s.shelves, req.Id)
	return &proto.Empty{}, nil
}

func main() {
	server := newServer()

	// 启动 gRPC 服务
	go runGrpcServer(server)

	// 启动 HTTP/gRPC-Gateway
	runHttpGateway()
}

func runGrpcServer(srv *protoServer) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterBookstoreServer(grpcServer, srv)

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func runHttpGateway() {
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := proto.RegisterBookstoreHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP gateway: %v", err)
	}

	log.Println("HTTP gateway listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
