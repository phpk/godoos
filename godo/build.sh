#!/bin/bash

# 定义要构建的平台
PLATFORMS=("linux/amd64" "windows/amd64" "darwin/amd64" "linux/arm64" "windows/arm64" "darwin/arm64")

# 定义版本号
SCRIPT_VERSION="1.0.0"
# 获取当前脚本的绝对路径
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
for PLATFORM in "${PLATFORMS[@]}"; do
    # 分割平台字符串
    OS=$(echo $PLATFORM | cut -d '/' -f 1)
    ARCH=$(echo $PLATFORM | cut -d '/' -f 2)

    # 设置后缀
    case $OS in
        "windows") SUFFIX=".exe" ;;
        *) SUFFIX="" ;;
    esac
    OUT_PATH="../dist/"
    if [ ! -d "$OUT_PATH" ]; then
        mkdir "$OUT_PATH"
    fi
    # 输出文件名
    OUTPUT_FILE="${OUT_PATH}godoos_web_${OS}_${ARCH}${SUFFIX}"

    # 设置GOOS和GOARCH环境变量
    export GOOS=$OS
    export GOARCH=$ARCH
    export GODOTOPTYPE="web"

    # 执行编译命令，并处理可能的错误
    go build  -ldflags="-s -w" -o "$OUTPUT_FILE" ./main.go || { echo "编译 $OS/$ARCH 失败，请检查错误并尝试解决。"; continue; }

    echo "编译 $OS/$ARCH 成功，生成文件: $OUTPUT_FILE"
done
