package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"github.com/herry-go/ksl/scram"
	"time"

	"github.com/IBM/sarama"
	"github.com/spf13/cobra"
)

var (
	brokers     []string
	username    string
	password    string
	topic       string
	message     string
	scramType   string // 新增：SCRAM认证类型
	enableTLS   bool
)

func getKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	
	// SASL 配置
	config.Net.SASL.Enable = true
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.Net.SASL.Handshake = true

	// 根据选择的认证类型设置
	switch scramType {
	case "SHA256":
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &scram.XDGSCRAMClient{HashGeneratorFcn: scram.SHA256}
		}
	default: // 默认使用 SHA512
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &scram.XDGSCRAMClient{HashGeneratorFcn: scram.SHA512}
		}
	}

	// TLS 配置
	config.Net.TLS.Enable = enableTLS
	if enableTLS {
		config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	// 连接配置
	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second
	
	// 重试配置
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 100 * time.Millisecond
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	
	// 消费者配置
	config.Consumer.Retry.Backoff = 100 * time.Millisecond
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	
	// 版本配置
	config.Version = sarama.V2_0_1_0

	return config
}

func produceMessage(cmd *cobra.Command, args []string) error {
	config := getKafkaConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("创建生产者失败: %v", err)
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}

	fmt.Printf("消息已发送到分区 %d，偏移量 %d\n", partition, offset)
	return nil
}

func consumeMessages(cmd *cobra.Command, args []string) error {
	config := getKafkaConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return fmt.Errorf("创建消费者失败: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("创建分区消费者失败: %v", err)
	}
	defer partitionConsumer.Close()

	fmt.Printf("开始消费主题 %s 的消息...\n", topic)
	for msg := range partitionConsumer.Messages() {
		fmt.Printf("分区: %d, 偏移量: %d, 消息: %s\n", msg.Partition, msg.Offset, string(msg.Value))
	}
	return nil
}

func listTopics(cmd *cobra.Command, args []string) error {
	config := getKafkaConfig()
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return fmt.Errorf("创建管理员客户端失败: %v", err)
	}
	defer admin.Close()

	topics, err := admin.ListTopics()
	if err != nil {
		return fmt.Errorf("获取主题列表失败: %v", err)
	}

	fmt.Println("主题列表:")
	for topic, detail := range topics {
		fmt.Printf("- %s (分区数: %d)\n", topic, detail.NumPartitions)
	}
	return nil
}

func createTopic(cmd *cobra.Command, args []string) error {
	config := getKafkaConfig()
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return fmt.Errorf("创建管理员客户端失败: %v", err)
	}
	defer admin.Close()

	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	if err != nil {
		return fmt.Errorf("创建主题失败: %v", err)
	}

	fmt.Printf("主题 %s 创建成功\n", topic)
	return nil
}

func deleteTopic(cmd *cobra.Command, args []string) error {
	config := getKafkaConfig()
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return fmt.Errorf("创建管理员客户端失败: %v", err)
	}
	defer admin.Close()

	err = admin.DeleteTopic(topic)
	if err != nil {
		return fmt.Errorf("删除主题失败: %v", err)
	}

	fmt.Printf("主题 %s 删除成功\n", topic)
	return nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "kcl",
		Short: "Kafka 命令行测试工具",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(brokers) == 0 {
				log.Fatal("请指定 Kafka 代理地址")
			}
			if username == "" || password == "" {
				log.Fatal("请指定 SCRAM 认证的用户名和密码")
			}
			if scramType != "SHA256" && scramType != "SHA512" {
				log.Fatal("SCRAM 认证类型必须是 SHA256 或 SHA512")
			}
		},
	}

	rootCmd.PersistentFlags().StringSliceVarP(&brokers, "brokers", "b", []string{}, "Kafka 代理地址列表 (例如: localhost:9092)")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "SCRAM 认证用户名")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "SCRAM 认证密码")
	rootCmd.PersistentFlags().StringVarP(&topic, "topic", "t", "", "主题名称")
	rootCmd.PersistentFlags().StringVarP(&scramType, "scram", "s", "SHA512", "SCRAM 认证类型 (SHA256 或 SHA512)")
	rootCmd.PersistentFlags().BoolVarP(&enableTLS, "tls", "", true, "是否启用 TLS")

	// 生产消息命令
	produceCmd := &cobra.Command{
		Use:   "produce",
		Short: "生产消息",
		RunE:  produceMessage,
	}
	produceCmd.Flags().StringVarP(&message, "message", "m", "", "要发送的消息")
	produceCmd.MarkFlagRequired("message")
	rootCmd.AddCommand(produceCmd)

	// 消费消息命令
	consumeCmd := &cobra.Command{
		Use:   "consume",
		Short: "消费消息",
		RunE:  consumeMessages,
	}
	rootCmd.AddCommand(consumeCmd)

	// 列出主题命令
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "列出所有主题",
		RunE:  listTopics,
	}
	rootCmd.AddCommand(listCmd)

	// 创建主题命令
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "创建主题",
		RunE:  createTopic,
	}
	rootCmd.AddCommand(createCmd)

	// 删除主题命令
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "删除主题",
		RunE:  deleteTopic,
	}
	rootCmd.AddCommand(deleteCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
