build:
	#一定要注意Makefile中的缩进，否则make build 可能报错Nothing ti be done for build
	#protoc 命令前面是一个Tab，不是四个或八个空格
	#protoc -I. --go_out=plugins=grpc:$(GOPATH)/src/facedamon/shippy/consignment-service proto/consignment/consignment.proto
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. ./proto/consignment/consignment.proto

	#告知go编译器生成二进制文件的目标环境: amd64 CPU的Linux系统
	GOOS=linux GOARCH=amd64 go build

	#docker 编译服务
	docker build -t consignment-service .

run:
	# 在 Docker alpine 容器的 50001 端口上运行 consignment-service 服务
	# 可添加 -d 参数将微服务放到后台运行
	docker run -p 50051:50051 \
	 -e MICRO_SERVER_ADDRESS=:50051 \
	 -e MICRO_REGISTRY=mdns \
	 consignment-service