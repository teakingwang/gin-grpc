package main

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teakingwang/gin-grpc/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	go grpcServer()
	r := gin.Default()
	r.GET("/say", Say)
	r.Run(":18081")
}

func grpcServer() {
	lis, err := net.Listen("tcp", ":18082")
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	pb.RegisterHelloServiceServer(s, &helloService{})

	s.Serve(lis)
}

type helloService struct{}

func (h helloService) Hello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	return &pb.HelloResp{Result: "hello " + req.Name}, nil
}

func Say(ctx *gin.Context) {
	//conn, err := grpc.Dial("localhost:10088", grpc.WithInsecure()) // 旧的方式，已废弃
	conn, err := grpc.Dial("localhost:10088", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	c := pb.NewHelloServiceClient(conn)

	name := ctx.Query("name")
	r, err := c.Hello(context.Background(), &pb.HelloReq{Name: name})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, r)
}
