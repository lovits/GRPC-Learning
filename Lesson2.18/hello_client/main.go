package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"hello_client/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// hello_client

const (
	defaultName = "qimi"
	defaultX    = 10
	defaultY    = 20
)

var (
	addr = flag.String("addr", "127.0.0.1:8972", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
	x    = flag.Int64("x", defaultX, "x is 10")
	y    = flag.Int64("y", defaultY, "y is 20")
)

func main() {
	flag.Parse()
	//有证书的安全连接
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "localhost")
	if err != nil {
		fmt.Printf("credentials failed err:%v\n", err)
	}
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(creds))
	// 连接到server端，此处禁用安全传输
	// conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 执行RPC调用并打印收到的响应数据
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetReply())

	addres, err := c.Add(ctx, &pb.AddRequest{X: *x, Y: *y})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("caculate: %d+%d=%d", *x, *y, addres.GetResult())

}
