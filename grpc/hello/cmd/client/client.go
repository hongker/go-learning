package main

import (
	"context"
	"fmt"
	"github.com/go-learning/grpc/hello/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	// 与服务器建立连接
	cc, err := grpc.Dial(":8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
	)
	if err != nil {
		log.Fatalf("failed to connect server: %v", err)
	}

	client := pb.NewGreeterClient(cc)
	// 配置参数
	req := &pb.HelloRequest{Name: "alice"}

	// 发送请求
	resp, err := client.SayHello(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to say hello: %v", err)
	}

	log.Println("success receive response: ", resp.Message)
}
