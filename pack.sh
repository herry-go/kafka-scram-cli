#!/bin/bash

# 定义版本和项目名
VERSION="v1.0.0"
PROJECT_NAME=ksc
OUTPUT_DIR=dist

# 创建输出目录
mkdir -p ${OUTPUT_DIR}

# 支持的平台架构
PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64")

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM##*/}

    echo "Building for $GOOS/$GOARCH"

    # 构建命令
    CGO_ENABLED=0 go build -o ${OUTPUT_DIR}/${PROJECT_NAME}-${GOOS}-${GOARCH} \
        -ldflags "-s -w -X 'main.version=${VERSION}'" \
        ./ksc/main.go

    # 打包
    cp install.sh ${OUTPUT_DIR}/install.sh
    (cd ${OUTPUT_DIR} && tar -czf ${PROJECT_NAME}-${GOOS}-${GOARCH}-${VERSION}.tar.gz ${PROJECT_NAME}-${GOOS}-${GOARCH} install.sh)
done

echo "Build complete. Files in: ${OUTPUT_DIR}/"
