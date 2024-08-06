FROM golang:alpine AS builder

# 在容器内部设置环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 设置后续指令的工作目录
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 将代码编译成二进制可执行文件
RUN go build -o godoos ./godo/main.go

# 创建一个小镜像
FROM debian:stretch-slim

# 设置工作目录
WORKDIR /

# 复制前端构建结果
COPY ./frontend/dist /dist

# 从builder镜像中把 /build/godoos 拷贝到当前目录
COPY --from=builder /build/godoos /godoos

# 添加执行权限
RUN chmod +x /godoos

# 暴露端口
EXPOSE 56780

# 需要运行的命令
ENTRYPOINT ["/godoos"]