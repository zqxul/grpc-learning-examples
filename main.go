package main

import (
	"net"

	_ "github.com/fullstorydev/grpcui"
	_ "github.com/fullstorydev/grpcurl"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"learn.org/core"
	_ "learn.org/echo"
	_ "learn.org/hello"
)

var (
	Addr = ":8081"
)

func init() {}

func main() {
	var err error
	listen, err := net.Listen("tcp", Addr)
	if err != nil {
		grpclog.Printf("listen err:%v", err)
	}
	reflection.Register(core.S)
	if err = core.S.Serve(listen); err != nil {
		grpclog.Printf("Serve err:%v", err)
	}
	grpclog.Println("Listen on " + Addr)
}
