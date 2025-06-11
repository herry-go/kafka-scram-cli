# 使用多阶段构建
# 构建阶段
FROM registry.cn-hangzhou.aliyuncs.com/redbean/golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

## 复制 go.mod 和 go.sum
#COPY go.mod go.sum ./
#
## 下载依赖
#RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o ksc main.go || exit 1 ;

# 运行阶段
FROM registry.cn-hangzhou.aliyuncs.com/redbean/alpine:latest

# 安装基础工具
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/ksc .

RUN ls /app
# 设置入口点
ENTRYPOINT ["/app/ksc"]

# 默认命令
CMD ["--help"]