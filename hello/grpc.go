package hello

import (
	"context"

	"google.golang.org/grpc"
	"learn.org/core"
)

type Hello interface {
	SayHello(ctx context.Context, in *HelloRequest) (*HelloResponse, error)
}

type HelloServiceImpl struct {
}

func (HelloServiceImpl) SayHello(ctx context.Context, in *HelloRequest) (*HelloResponse, error) {
	return &HelloResponse{Reply: "hello, " + in.Name}, nil
}

func init() {
	core.S.RegisterService(&grpc.ServiceDesc{
		ServiceName: "hello.HelloService",
		HandlerType: (*Hello)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "SayHello",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := &HelloRequest{}
					if err := dec(in); err != nil {
						return nil, err
					}
					if interceptor == nil {
						return srv.(Hello).SayHello(ctx, in)
					}
					info := &grpc.UnaryServerInfo{
						Server:     srv,
						FullMethod: "/HelloService/SayHello",
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(Hello).SayHello(ctx, req.(*HelloRequest))
					}
					return interceptor(ctx, in, info, handler)
				},
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "hello/hello.proto",
	}, HelloServiceImpl{})
}
