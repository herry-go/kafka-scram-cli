#!/bin/bash

# 默认配置
USE_DOCKER=true
DOCKER_IMAGE="ksc:latest"
DOCKER_NETWORK="kafka-cluster-network"
KAFKA_BROKER="kafka-broker-1:19094"
KAFKA_USER="kafkaclient"
KAFKA_PASSWORD="password"
SCRAM_TYPE="SHA512"
ENABLE_TLS=false


# 显示帮助信息
show_help() {
    echo "Kafka 命令行工具使用说明"
    echo "用法: $0 [选项] 命令 [参数]"
    echo ""
    echo "选项:"
    echo "  -d, --docker        使用 Docker 运行"
    echo "  -b, --broker        指定 Kafka broker 地址 (默认: $KAFKA_BROKER)"
    echo "  -u, --user          指定用户名 (默认: $KAFKA_USER)"
    echo "  -p, --password      指定密码 (默认: $KAFKA_PASSWORD)"
    echo "  -n, --network       指定 Docker 网络 (默认: $DOCKER_NETWORK)"
    echo "  -i, --image         指定 Docker 镜像 (默认: $DOCKER_IMAGE)"
    echo "  -s, --scram         SCRAM 认证类型 (SHA256 或 SHA512) (默认: $SCRAM_TYPE)"
    echo "  --no-tls           禁用 TLS (默认启用)"
    echo "  -h, --help          显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 produce -t test -m 'Hello'"
    echo "  $0 -d produce -t test -m 'Hello'"
    echo "  $0 -d -b localhost:9092 -s SHA512 --no-tls produce -t test -m 'Hello'"
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--docker)
            USE_DOCKER=true
            shift
            ;;
        -b|--broker)
            KAFKA_BROKER="$2"
            shift 2
            ;;
        -u|--user)
            KAFKA_USER="$2"
            shift 2
            ;;
        -p|--password)
            KAFKA_PASSWORD="$2"
            shift 2
            ;;
        -n|--network)
            DOCKER_NETWORK="$2"
            shift 2
            ;;
        -i|--image)
            DOCKER_IMAGE="$2"
            shift 2
            ;;
        -s|--scram)
            SCRAM_TYPE="$2"
            if [[ "$SCRAM_TYPE" != "SHA256" && "$SCRAM_TYPE" != "SHA512" ]]; then
                echo "错误: SCRAM 认证类型必须是 SHA256 或 SHA512"
                exit 1
            fi
            shift 2
            ;;
        --no-tls)
            ENABLE_TLS=false
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            break
            ;;
    esac
done

# 构建基本命令参数
BASE_ARGS="-b $KAFKA_BROKER -u $KAFKA_USER -p $KAFKA_PASSWORD -s $SCRAM_TYPE"
if [ "$ENABLE_TLS" = false ]; then
    BASE_ARGS="$BASE_ARGS --tls=false"
fi

# 如果使用 Docker
if [ "$USE_DOCKER" = true ]; then
    # 检查 Docker 是否运行
    if ! docker info > /dev/null 2>&1; then
        echo "错误: Docker 未运行"
        exit 1
    fi

    # 检查镜像是否存在
    if ! docker image inspect $DOCKER_IMAGE > /dev/null 2>&1; then
        echo "错误: Docker 镜像 $DOCKER_IMAGE 不存在"
        exit 1
    fi

    # 构建 Docker 运行命令
    CMD="docker run -it --rm --network $DOCKER_NETWORK $DOCKER_IMAGE $BASE_ARGS $@"
else
    # 检查本地二进制文件是否存在
    if [ ! -f "./ksc" ]; then
        echo "错误: 本地 ksc 二进制文件不存在"
        exit 1
    fi

    # 构建本地运行命令
    CMD="./ksc $BASE_ARGS $@"
fi

# 执行命令
echo "执行命令: $CMD"
eval $CMD