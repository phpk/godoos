#!/bin/bash

# 获取用户输入的平台和架构，默认为当前系统平台和架构
GOOS=${1:-$(go env GOOS)}
GOARCH=${2:-$(go env GOARCH)}

# 设置二进制文件的输出名称
BINARY_NAME="godo"

# 编译项目
echo "Building $BINARY_NAME for $GOOS/$GOARCH..."
go build  -ldflags="-s -w" -o $BINARY_NAME

echo "Build completed."
