#!/bin/bash

BINARY_NAME=ksc
CURRENT_DIR=$(pwd)
BIN_DIR="/usr/local/bin"

# 查找当前目录下的二进制文件
BINARY_PATH=$(find "$CURRENT_DIR" -type f -name "$BINARY_NAME-*" | head -n1)

if [ -z "$BINARY_PATH" ]; then
    echo "错误：未找到二进制文件。请确保当前目录包含 '$BINARY_NAME-*' 可执行文件。"
    exit 1
fi

# 安装
sudo cp "$BINARY_PATH" "$BIN_DIR/$BINARY_NAME"
sudo chmod +x "$BIN_DIR/$BINARY_NAME"

echo "✅ $BINARY_NAME 已成功安装到 $BIN_DIR/"
echo "💡 使用方式: $BINARY_NAME --help"
