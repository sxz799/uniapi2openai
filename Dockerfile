# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.21-alpine as builder

# 设置工作目录
WORKDIR /home

# 将应用的代码复制到容器中
COPY . .


# 编译应用程序
RUN go build -o app .

FROM alpine:latest

WORKDIR /home


COPY --from=0 /home/app ./


EXPOSE 8080

# 运行应用程序
CMD ["./app"]