# docker build --no-cache -t godoos/godoos:latest .
# docker run -it --rm -p 56780:56780 -p 8185:80 godoos/godoos:latest
# 使用 golang:alpine 作为构建阶段的基础镜像
FROM golang:alpine AS builder

# 在容器内部设置环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
    GODOTOPTYPE=docker

# 设置后续指令的工作目录
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 将代码编译成二进制可执行文件
RUN go build -o godoos ./godo/main.go

# 创建一个小镜像
FROM nginx:alpine

# 设置工作目录
WORKDIR /

# 更改镜像源
RUN echo 'http://mirrors.aliyun.com/alpine/v3.19/main' > /etc/apk/repositories && \
    echo 'http://mirrors.aliyun.com/alpine/v3.19/community' >> /etc/apk/repositories

# 复制前端构建结果到 nginx 的默认文档目录
COPY --from=builder /build/frontend/dist /usr/share/nginx/html

# 替换默认的 nginx 配置文件
COPY docker/nginx.conf /etc/nginx/conf.d/

# 暴露 nginx 默认端口
EXPOSE 80

# 从builder镜像中把 /build/godoos 拷贝到当前目录
COPY --from=builder /build/godoos /godoos

# 添加执行权限
RUN chmod +x /godoos

# 暴露端口
EXPOSE 56780

# 需要运行的命令
USER root

# 启动 Nginx 和 godoos 服务
CMD ["sh", "-c", "nginx -g 'daemon off;' & /godoos"]