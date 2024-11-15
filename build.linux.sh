#!/bin/bash

# 构建和压缩
wails build -platform linux/amd64 -s -ldflags="-s -w"
wails build -platform linux/arm64 -s -ldflags="-s -w"
wails build -platform darwin/amd64 -s -ldflags="-s -w"
wails build -platform darwin/arm64 -s  -ldflags="-s -w"