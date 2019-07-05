
#运行环境
FROM alpine:latest

#在容器根目录创建app目录
RUN mkdir /app

#切换工作目录
WORKDIR /app

#将微服务的服务端运行文件拷贝至app下
ADD consignment-service /app/consignment-service

#运行服务端
CMD ["./consignment-service"]

