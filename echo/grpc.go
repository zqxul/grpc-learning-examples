package echo

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"learn.org/core"

	"context"
)

// 服务接口
type EchoService interface {
	UnaryEcho(context.Context, *EchoRequest) (*EchoResponse, error)
	ServerStreamingEcho(*EchoRequest, grpc.ServerStream) error
	ClientStreamingEcho(grpc.ServerStream) error
	BidirectionalStreamingEcho(grpc.ServerStream) error
}

// 服务实现
type EchoServiceImpl struct {
}

// 普通输出实现
func (EchoServiceImpl) UnaryEcho(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	data, _ := json.Marshal(req)
	log.Println("receive Unary Echo from client:" + string(data))
	return &EchoResponse{
		Message: "pong",
	}, nil
}

// 服务端流输出实现
func (EchoServiceImpl) ServerStreamingEcho(req *EchoRequest, stream grpc.ServerStream) error {
	data, _ := json.Marshal(req)
	log.Println("server stream receive Unary Echo from client:" + string(data))
	// 开启一个goroutine, 创建一个定时器，每5秒往客户端发送数据
	ticker := time.NewTicker(time.Second * 5)

	startTime := time.Now()
	ctx, cancel := context.WithCancel(stream.Context())
	go func() {
		for {
			// 30秒后，终止推送，跳出循环，数据推送结束
			if time.Now().After(startTime.Add(time.Second * 30)) {
				log.Println("ServerStream is done")
				cancel()
				break
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("ServerStream closed")
			return ctx.Err()
		case <-ticker.C:
			res := &EchoResponse{Message: time.Now().String()}
			stream.SendMsg(res)
			data, _ := json.Marshal(res)
			log.Println("server stream send message to client:" + string(data))
		}
	}
}

// 客户端流输出实现
func (EchoServiceImpl) ClientStreamingEcho(stream grpc.ServerStream) error {

	// 创建数据传输通道
	c := make(chan *EchoRequest)

	// 开启goroutine接收客户端数据
	go func(c chan<- *EchoRequest) {
		for {
			in := &EchoRequest{}
			err := stream.RecvMsg(in)
			if err == io.EOF {
				break
			}
			c <- in //接收成功后发到通道中
		}
	}(c)

	// 循环等待数据
	for {
		select {
		case <-stream.Context().Done():
			log.Println("ClientStream closed")
			break
		case in := <-c:
			data, _ := json.Marshal(in)
			log.Println("receive client stream Echo from client:" + string(data))
		}
	}

}

// 双向流输出实现
func (EchoServiceImpl) BidirectionalStreamingEcho(stream grpc.ServerStream) error {

	startTime := time.Now()
	ctx, cancel := context.WithCancel(stream.Context())
	go func() {
		for {
			// 30秒后，终止推送，跳出循环，数据推送结束
			if time.Now().After(startTime.Add(time.Second * 30)) {
				log.Println("ServerStream is done")
				cancel()
				break
			}
		}
	}()

	// 开启一个goroutine, 创建一个定时器，每3秒往客户端发送数据
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				res := &EchoResponse{Message: time.Now().String()}
				stream.SendMsg(res)
				data, _ := json.Marshal(res)
				log.Println("send server stream Echo to client:" + string(data))
			}
		}
	}()

	// 创建数据传输通道
	c := make(chan *EchoRequest)

	// 开启goroutine接收客户端数据
	go func(c chan<- *EchoRequest) {
		for {
			in := &EchoRequest{}
			err := stream.RecvMsg(in)
			if err == io.EOF {
				break
			}
			c <- in //接收成功后发到通道中
		}
	}(c)

	// 循环等待数据
	for {
		select {
		case <-ctx.Done(): // 流关闭后退出循环
			log.Println("BidirectionalStream closed")
			return ctx.Err()
		case in := <-c: // 从通道中接收客户端数据，并打印
			data, _ := json.Marshal(in)
			log.Println("receive client stream Echo from client:" + string(data))
		}
	}

}

var (
	// 服务方法描述
	sd = grpc.ServiceDesc{
		ServiceName: "echo.EchoService",
		HandlerType: (*EchoService)(nil),
		// 普通方法
		Methods: []grpc.MethodDesc{
			{
				MethodName: "UnaryEcho",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := &EchoRequest{}
					if err := dec(in); err != nil {
						grpclog.Printf("decode request err:%v", err)
						return nil, err
					}
					if interceptor == nil {
						return srv.(EchoService).UnaryEcho(ctx, in)
					}
					info := &grpc.UnaryServerInfo{
						Server:     srv,
						FullMethod: "/echo.EchoService/UnaryEcho",
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(EchoService).UnaryEcho(ctx, req.(*EchoRequest))
					}
					return interceptor(ctx, in, info, handler)
				},
			},
		},
		// 流方法
		Streams: []grpc.StreamDesc{
			{
				StreamName: "ServerStreamingEcho",
				Handler: func(srv interface{}, stream grpc.ServerStream) error {
					in := &EchoRequest{}
					stream.RecvMsg(in)
					log.Println(json.Marshal(in))
					return srv.(EchoService).ServerStreamingEcho(in, stream)
				},
				ClientStreams: false,
				ServerStreams: true,
			},
			{
				StreamName: "ClientStreamingEcho",
				Handler: func(srv interface{}, stream grpc.ServerStream) error {
					return srv.(EchoService).ClientStreamingEcho(stream)
				},
				ClientStreams: true,
				ServerStreams: false,
			},
			{
				StreamName: "BidirectionalStreamingEcho",
				Handler: func(srv interface{}, stream grpc.ServerStream) error {
					return srv.(EchoService).BidirectionalStreamingEcho(stream)
				},
				ClientStreams: true,
				ServerStreams: true,
			},
		},
		Metadata: "echo/echo.proto",
	}
)

func init() {
	// 注册服务
	core.S.RegisterService(&sd, EchoServiceImpl{})
}
