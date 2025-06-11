# Kafka SCRAM 命令行测试工具

本工具基于 Go 语言开发，支持 Kafka SCRAM 认证，方便测试 Kafka 集群的基本功能。

## 功能列表
- 生产消息
- 消费消息
- 列出所有主题
- 创建主题
- 删除主题

## 依赖
- Go 1.20 及以上
- [sarama](https://github.com/IBM/sarama)
- [cobra](https://github.com/spf13/cobra)

## 本地编译

```bash
go mod tidy
go build -o ksc main.go
```

注意：如果 Kafka 运行在容器外，需要使用主机的 IP 地址而不是 localhost。

## 使用方法

所有命令都需指定 Kafka broker 地址、SCRAM 用户名和密码。

### 1. 生产消息

```bash
./ksc produce -b 127.0.0.1:19093 -u kafkaclient -p password -t test -m "你好，Kafka！"
```

### 2. 消费消息

```bash
./ksc consume -b 127.0.0.1:19093 -u kafkaclient -p password -t test
```

### 3. 列出所有主题

```bash
./ksc list -b 127.0.0.1:19093 -u kafkaclient -p password
```

### 4. 创建主题

```bash
./ksc create -b 127.0.0.1:19093 -u kafkaclient -p password -t test
```

### 5. 删除主题

```bash
./ksc delete -b 127.0.0.1:19093 -u kafkaclient -p password -t test
```

### 连接多个 broker

用逗号分隔多个 broker 地址：

```bash
./ksc list -b 127.0.0.1:19093,127.0.0.1:29093 -u kafkaclient -p password
```

## 参数说明
- `-b, --brokers`   Kafka broker 地址，多个用逗号分隔
- `-u, --username`  SCRAM 认证用户名
- `-p, --password`  SCRAM 认证密码
- `-t, --topic`     主题名称（部分命令需要）
- `-m, --message`   要发送的消息（仅生产消息命令需要）

## 注意事项
- 生产环境请配置正确的 TLS 证书，不要使用 `InsecureSkipVerify: true`
- 主题的分区数和副本数可根据需要在代码中调整
- 使用 Docker 时，确保容器可以访问到 Kafka 服务器

---
如有问题欢迎反馈！ 