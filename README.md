# grpc-learning-examples
grpc 样例
1、普通请求
2、客户端流
3、服务端流
4、双向流

开发步骤 (需要安装protoc，安装grpc插件)
1、编写proto文件 (*.proto)
2、生成.pb.go文件, 命令 proto -I. --go_out=<.pb.go的输出路径> --plugin=grpc:. <proto文件路径,eg: hello/hello.proto>
3、编写服务并注册服务 (grpc.go)

grpc服务启动步骤
1、go mod tidy
2、go mod vendor
3、go run main.go


grpcui调试
1、安装好grpcui，并设置好命令环境
2、grpcui --plaintext localhost:8081 （8081为grpc服务端口，命令执行后可以在浏览器调试）
