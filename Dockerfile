FROM golang:latest

# 设置工作目录
WORKDIR /app/

# 将当前目录的内容复制到容器中的/app目录
COPY . .
ENV GOPATH=$GOPATH:/app/ GOPROXY=https://mirrors.aliyun.com/goproxy,https://goproxy.cn,direct
# 构建二进制文件
RUN go build -o godoos ./godo/main.go

# 使用更小的基础镜像用于最终的部署
FROM alpine:latest
WORKDIR /app

# 将构建好的二进制文件复制到新的镜像中
COPY --from=0 /app/godoos /app/
# 复制前端构建结果
COPY --from=0 /app/frontend/dist /app/dist

# 暴露端口
EXPOSE 8210

# 运行命令
CMD ["./godoos"]