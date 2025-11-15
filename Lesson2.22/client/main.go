package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"lesson222/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	// 1. 连接 gRPC 服务端
	conn, err := grpc.Dial(
		"127.0.0.1:8972", // gRPC 服务端口
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// 2. 创建客户端
	client := proto.NewBookstoreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// -------------------------------
	// 3. 创建书架
	fmt.Println("== 创建书架 ==")
	createResp, err := client.CreateShelf(ctx, &proto.CreateShelfRequest{
		Shelf: &proto.Shelf{
			Theme: "Science",
			Size:  50,
		},
	})
	if err != nil {
		log.Fatalf("CreateShelf failed: %v", err)
	}
	fmt.Printf("创建成功：ID=%d, Theme=%s, Size=%d\n", createResp.Id, createResp.Theme, createResp.Size)

	// -------------------------------
	// 4. 查询指定书架
	fmt.Println("== 查询书架 ==")
	getResp, err := client.GetShelf(ctx, &proto.GetShelfRequest{
		Shelf: createResp.Id,
	})
	if err != nil {
		log.Fatalf("GetShelf failed: %v", err)
	}
	fmt.Printf("查询结果：ID=%d, Theme=%s, Size=%d\n", getResp.Id, getResp.Theme, getResp.Size)

	// -------------------------------
	// 5. 列出所有书架
	fmt.Println("== 列出所有书架 ==")
	listResp, err := client.ListShelves(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("ListShelves failed: %v", err)
	}
	for _, s := range listResp.Shelves {
		fmt.Printf("ID=%d, Theme=%s, Size=%d\n", s.Id, s.Theme, s.Size)
	}

	// -------------------------------
	// 6. 删除书架
	fmt.Println("== 删除书架 ==")
	_, err = client.DeleteShelf(ctx, &proto.DeleteShelfRequest{
		Shelf: createResp.Id,
	})
	if err != nil {
		log.Fatalf("DeleteShelf failed: %v", err)
	}
	fmt.Println("删除成功")
}
