#!/bin/bash

# 定义一个函数来执行压缩操作
compress_to_zip() {
    local target_file="./build/bin/godoos$2"
    local zip_name="./build/bin/$1.zip"
    
    # 检查ZIP文件是否已经存在，如果存在则跳过压缩
    if [ ! -f "$zip_name" ]; then
        echo "Compressing $target_file to $zip_name"
        zip -j "$zip_name" "$target_file"
    else
        echo "$zip_name already exists. Skipping compression."
    fi
}

# 构建和压缩
wails build -platform linux/amd64 -s
compress_to_zip "godoos-linux-amd64" ""
wails build -platform linux/arm64 -s
compress_to_zip "godoos-linux-arm64" ""
wails build -platform darwin/amd64 -s
compress_to_zip "godoos-darwin-amd64" ".app"
wails build -platform darwin/arm64 -s
compress_to_zip "godoos-darwin-arm64" ".app"