package main

// server.go

import (
	"fmt"
	"log"
	"net"

	pb "hello"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":2019"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Req: %v", req)
	return &pb.HelloResponse{Reply: "Hello " + req.Greeting, Number: []int32{1, 2, 3}}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})
	serviceInfos := s.GetServiceInfo()
	for k, v := range serviceInfos {
		fmt.Printf("s.name:[%v]\n", k)
		for _, m := range v.Methods {
			fmt.Printf("    m.name:[%v]\n", m.Name)
		}
	}

	s.Serve(lis)
}
