package main

import (
	"context"
	"github.com/go-learning/grpc/hello/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	// 监听8080段口
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor),
	)
	pb.RegisterGreeterServer(grpcServer, NewServer())
	grpcServer.Serve(lis)
}

type Server struct {
	pb.UnimplementedGreeterServer
}

func (s Server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Welcome " + request.Name}, nil
}

func NewServer() *Server {
	return &Server{}
}

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 在调用方法之前记录日志
	before := time.Now()

	// 调用实际的gRPC方法
	resp, err := handler(ctx, req)

	// 在调用方法之后记录日志
	log.Printf("After Calling Method: %s, time used:%v", info.FullMethod, time.Since(before))

	// 返回结果
	return resp, err
}
