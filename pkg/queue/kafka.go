package queue

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/YoungBoyGod/OneGoServer/internal/config"
	pkglog "github.com/YoungBoyGod/OneGoServer/pkg/log"
	"go.uber.org/zap"
)

// === 类型定义 ===

// KafkaMessage Kafka消息结构
type KafkaMessage struct {
	Topic     string            `json:"topic"`
	Key       string            `json:"key,omitempty"`
	Value     interface{}       `json:"value"`
	Headers   map[string]string `json:"headers,omitempty"`
	Partition int32             `json:"partition,omitempty"`
	Offset    int64             `json:"offset,omitempty"`
	Timestamp time.Time         `json:"timestamp,omitempty"`
}

// MessageHandler 消息处理函数类型
type MessageHandler func(ctx context.Context, message *KafkaMessage) error

// KafkaStats Kafka统计信息
type KafkaStats struct {
	ProducerStats ProducerStats `json:"producer_stats"`
	ConsumerStats ConsumerStats `json:"consumer_stats"`
	IsConnected   bool          `json:"is_connected"`
	LastPing      time.Time     `json:"last_ping"`
}

// ProducerStats 生产者统计
type ProducerStats struct {
	MessagesSent      int64 `json:"messages_sent"`
	MessagesSucceeded int64 `json:"messages_succeeded"`
	MessagesFailed    int64 `json:"messages_failed"`
	BytesSent         int64 `json:"bytes_sent"`
}

// ConsumerStats 消费者统计
type ConsumerStats struct {
	MessagesReceived  int64 `json:"messages_received"`
	MessagesProcessed int64 `json:"messages_processed"`
	MessagesFailed    int64 `json:"messages_failed"`
	BytesReceived     int64 `json:"bytes_received"`
}

// KafkaManager Kafka管理器
type KafkaManager struct {
	config   *config.KafkaConfig
	producer sarama.SyncProducer
	consumer sarama.Consumer
	client   sarama.Client
	mu       sync.RWMutex
	isReady  bool
	stats    KafkaStats
	lastPing time.Time
}

// === 全局变量 ===

var (
	// 全局Kafka管理器
	kafkaManager *KafkaManager
	once         sync.Once
)

// === 配置验证 ===

// validateKafkaConfig 验证Kafka配置
func validateKafkaConfig(cfg *config.KafkaConfig) error {
	if cfg == nil {
		return errors.New("kafka config cannot be nil")
	}

	if len(cfg.Brokers) == 0 {
		return errors.New("kafka brokers list cannot be empty")
	}

	for _, broker := range cfg.Brokers {
		if strings.TrimSpace(broker) == "" {
			return errors.New("kafka broker address cannot be empty")
		}
	}

	if cfg.ClientID == "" {
		return errors.New("kafka client id is required")
	}

	if cfg.Version == "" {
		cfg.Version = "2.6.0" // 设置默认版本
	}

	// 验证版本格式
	if _, err := sarama.ParseKafkaVersion(cfg.Version); err != nil {
		return fmt.Errorf("invalid kafka version: %w", err)
	}

	// 设置默认值
	if cfg.ProducerRequiredAcks == 0 {
		cfg.ProducerRequiredAcks = 1 // WaitForLeader
	}

	if cfg.ProducerRetryMax == 0 {
		cfg.ProducerRetryMax = 3
	}

	if cfg.ProducerMaxMessageBytes == 0 {
		cfg.ProducerMaxMessageBytes = 1000000 // 1MB
	}

	if cfg.ConsumerGroupID == "" {
		cfg.ConsumerGroupID = "onegoserver-consumer-group"
	}

	if cfg.ConsumerOffsetInitial == "" {
		cfg.ConsumerOffsetInitial = "newest"
	}

	if cfg.ConsumerSessionTimeout == 0 {
		cfg.ConsumerSessionTimeout = 10000 // 10秒
	}

	if cfg.ConsumerHeartbeatInterval == 0 {
		cfg.ConsumerHeartbeatInterval = 3000 // 3秒
	}

	return nil
}

// === 初始化函数 ===

// InitKafka 初始化Kafka连接
func InitKafka(ctx context.Context, cfg *config.KafkaConfig) error {
	// 先验证配置，无论是否已经初始化过
	if err := validateKafkaConfig(cfg); err != nil {
		return fmt.Errorf("invalid kafka config: %w", err)
	}

	var initErr error

	once.Do(func() {
		// 创建Kafka管理器
		kafkaManager = &KafkaManager{
			config: cfg,
		}

		// 初始化连接
		if err := kafkaManager.connect(); err != nil {
			initErr = fmt.Errorf("failed to connect to kafka: %w", err)
			return
		}

		pkglog.LogInfo("Kafka initialized successfully",
			zap.Strings("brokers", cfg.Brokers),
			zap.String("client_id", cfg.ClientID),
			zap.String("version", cfg.Version))
	})

	return initErr
}

// InitKafkaWithRetry 带重试机制的Kafka初始化
func InitKafkaWithRetry(ctx context.Context, cfg *config.KafkaConfig, maxRetries int, retryInterval time.Duration) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if err := InitKafka(ctx, cfg); err != nil {
			lastErr = err
			pkglog.LogWarn("Kafka connection attempt failed",
				zap.Int("attempt", i+1),
				zap.Int("max_retries", maxRetries),
				zap.Error(err))

			if i < maxRetries-1 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(retryInterval):
					// 重置once以允许重试
					once = sync.Once{}
				}
			}
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to connect to kafka after %d attempts: %w", maxRetries, lastErr)
}

// === Kafka管理器方法 ===

// connect 建立Kafka连接
func (km *KafkaManager) connect() error {
	// 打印调试信息，查看实际的broker地址
	pkglog.LogInfo("Connecting to Kafka brokers",
		zap.Strings("brokers", km.config.Brokers),
		zap.String("client_id", km.config.ClientID))

	// 创建Sarama配置
	saramaConfig := sarama.NewConfig()

	// 设置版本
	version, err := sarama.ParseKafkaVersion(km.config.Version)
	if err != nil {
		return fmt.Errorf("failed to parse kafka version: %w", err)
	}
	saramaConfig.Version = version

	// 设置客户端ID
	saramaConfig.ClientID = km.config.ClientID

	// 设置SASL认证
	if km.config.EnableSASL {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.User = km.config.Username
		saramaConfig.Net.SASL.Password = km.config.Password

		switch strings.ToUpper(km.config.SASLMechanism) {
		case "PLAIN":
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		case "SCRAM-SHA-256":
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		case "SCRAM-SHA-512":
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		default:
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		}
	}

	// 设置TLS
	if km.config.EnableTLS {
		saramaConfig.Net.TLS.Enable = true
		saramaConfig.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	// 生产者配置
	saramaConfig.Producer.Return.Successes = km.config.ProducerReturnSuccesses
	saramaConfig.Producer.Return.Errors = km.config.ProducerReturnErrors
	saramaConfig.Producer.RequiredAcks = sarama.RequiredAcks(km.config.ProducerRequiredAcks)
	saramaConfig.Producer.Retry.Max = km.config.ProducerRetryMax
	saramaConfig.Producer.MaxMessageBytes = km.config.ProducerMaxMessageBytes

	// 消费者配置
	switch strings.ToLower(km.config.ConsumerOffsetInitial) {
	case "oldest":
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	case "newest":
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	default:
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	saramaConfig.Consumer.Group.Session.Timeout = time.Duration(km.config.ConsumerSessionTimeout) * time.Millisecond
	saramaConfig.Consumer.Group.Heartbeat.Interval = time.Duration(km.config.ConsumerHeartbeatInterval) * time.Millisecond

	// 创建客户端
	client, err := sarama.NewClient(km.config.Brokers, saramaConfig)
	if err != nil {
		return fmt.Errorf("failed to create kafka client: %w", err)
	}

	// 创建生产者
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to create kafka producer: %w", err)
	}

	// 创建消费者
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		producer.Close()
		client.Close()
		return fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	// 设置连接
	km.mu.Lock()
	km.client = client
	km.producer = producer
	km.consumer = consumer
	km.isReady = true
	km.lastPing = time.Now()
	km.mu.Unlock()

	// 记录操作日志
	pkglog.LogKafkaOperation("connect", "kafka_manager", time.Since(km.lastPing),
		zap.Strings("brokers", km.config.Brokers),
		zap.String("client_id", km.config.ClientID))

	return nil
}

// === 公共接口函数 ===

// IsReady 检查Kafka是否就绪
func IsReady() bool {
	if kafkaManager == nil {
		return false
	}

	kafkaManager.mu.RLock()
	defer kafkaManager.mu.RUnlock()
	return kafkaManager.isReady
}

// GetKafkaClient 获取Kafka客户端
func GetKafkaClient() sarama.Client {
	if kafkaManager == nil {
		return nil
	}

	kafkaManager.mu.RLock()
	defer kafkaManager.mu.RUnlock()
	return kafkaManager.client
}

// Ping 检查Kafka连接
func Ping(ctx context.Context) error {
	if !IsReady() {
		return errors.New("kafka client not initialized")
	}

	kafkaManager.mu.RLock()
	client := kafkaManager.client
	kafkaManager.mu.RUnlock()

	if client.Closed() {
		return errors.New("kafka client is closed")
	}

	// 获取broker列表作为健康检查
	brokers := client.Brokers()
	if len(brokers) == 0 {
		return errors.New("no active kafka brokers")
	}

	kafkaManager.mu.Lock()
	kafkaManager.lastPing = time.Now()
	kafkaManager.mu.Unlock()

	return nil
}

// === 生产者操作 ===

// SendMessage 发送消息
func SendMessage(ctx context.Context, message *KafkaMessage) error {
	if !IsReady() {
		return errors.New("kafka client not initialized")
	}

	kafkaManager.mu.RLock()
	producer := kafkaManager.producer
	kafkaManager.mu.RUnlock()

	if producer == nil {
		return errors.New("kafka producer not initialized")
	}

	// 序列化消息值
	var valueBytes []byte
	var err error

	switch v := message.Value.(type) {
	case string:
		valueBytes = []byte(v)
	case []byte:
		valueBytes = v
	default:
		valueBytes, err = json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal message value: %w", err)
		}
	}

	// 构建Sarama消息
	producerMessage := &sarama.ProducerMessage{
		Topic: message.Topic,
		Value: sarama.ByteEncoder(valueBytes),
	}

	// 设置Key（如果有）
	if message.Key != "" {
		producerMessage.Key = sarama.StringEncoder(message.Key)
	}

	// 设置Headers（如果有）
	if len(message.Headers) > 0 {
		headers := make([]sarama.RecordHeader, 0, len(message.Headers))
		for k, v := range message.Headers {
			headers = append(headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: []byte(v),
			})
		}
		producerMessage.Headers = headers
	}

	// 发送消息
	startTime := time.Now()
	partition, offset, err := producer.SendMessage(producerMessage)
	duration := time.Since(startTime)

	if err != nil {
		// 更新失败统计
		kafkaManager.mu.Lock()
		kafkaManager.stats.ProducerStats.MessagesFailed++
		kafkaManager.mu.Unlock()

		pkglog.LogKafkaOperation("send_message_failed", message.Topic, duration,
			zap.String("key", message.Key),
			zap.Error(err))

		return fmt.Errorf("failed to send message: %w", err)
	}

	// 更新成功统计
	kafkaManager.mu.Lock()
	kafkaManager.stats.ProducerStats.MessagesSent++
	kafkaManager.stats.ProducerStats.MessagesSucceeded++
	kafkaManager.stats.ProducerStats.BytesSent += int64(len(valueBytes))
	kafkaManager.mu.Unlock()

	// 设置返回信息
	message.Partition = partition
	message.Offset = offset
	message.Timestamp = time.Now()

	pkglog.LogKafkaOperation("send_message", message.Topic, duration,
		zap.String("key", message.Key),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset))

	return nil
}

// SendMessages 批量发送消息
func SendMessages(ctx context.Context, messages []*KafkaMessage) error {
	if len(messages) == 0 {
		return nil
	}

	for _, message := range messages {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := SendMessage(ctx, message); err != nil {
				return err
			}
		}
	}

	return nil
}

// === 消费者操作 ===

// ConsumeMessages 消费消息（简单消费者）
func ConsumeMessages(ctx context.Context, topic string, partition int32, offset int64, handler MessageHandler) error {
	if !IsReady() {
		return errors.New("kafka client not initialized")
	}

	kafkaManager.mu.RLock()
	consumer := kafkaManager.consumer
	kafkaManager.mu.RUnlock()

	if consumer == nil {
		return errors.New("kafka consumer not initialized")
	}

	// 创建分区消费者
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return fmt.Errorf("failed to create partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	pkglog.LogInfo("Started consuming messages",
		zap.String("topic", topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case message := <-partitionConsumer.Messages():
			if message == nil {
				continue
			}

			startTime := time.Now()

			// 构建消息
			kafkaMsg := &KafkaMessage{
				Topic:     message.Topic,
				Key:       string(message.Key),
				Value:     message.Value,
				Partition: message.Partition,
				Offset:    message.Offset,
				Timestamp: message.Timestamp,
			}

			// 设置Headers
			if len(message.Headers) > 0 {
				kafkaMsg.Headers = make(map[string]string)
				for _, header := range message.Headers {
					kafkaMsg.Headers[string(header.Key)] = string(header.Value)
				}
			}

			// 处理消息
			err := handler(ctx, kafkaMsg)
			duration := time.Since(startTime)

			// 更新统计
			kafkaManager.mu.Lock()
			kafkaManager.stats.ConsumerStats.MessagesReceived++
			kafkaManager.stats.ConsumerStats.BytesReceived += int64(len(message.Value))

			if err != nil {
				kafkaManager.stats.ConsumerStats.MessagesFailed++
				pkglog.LogKafkaOperation("process_message_failed", topic, duration,
					zap.String("key", kafkaMsg.Key),
					zap.Int32("partition", partition),
					zap.Int64("offset", message.Offset),
					zap.Error(err))
			} else {
				kafkaManager.stats.ConsumerStats.MessagesProcessed++
				pkglog.LogKafkaOperation("process_message", topic, duration,
					zap.String("key", kafkaMsg.Key),
					zap.Int32("partition", partition),
					zap.Int64("offset", message.Offset))
			}
			kafkaManager.mu.Unlock()

		case err := <-partitionConsumer.Errors():
			if err != nil {
				pkglog.LogError("Consumer error", zap.Error(err))
				return err
			}
		}
	}
}

// === 主题管理 ===

// GetTopics 获取所有主题
func GetTopics(ctx context.Context) ([]string, error) {
	if !IsReady() {
		return nil, errors.New("kafka client not initialized")
	}

	kafkaManager.mu.RLock()
	client := kafkaManager.client
	kafkaManager.mu.RUnlock()

	topics, err := client.Topics()
	if err != nil {
		return nil, fmt.Errorf("failed to get topics: %w", err)
	}

	return topics, nil
}

// GetPartitions 获取主题的分区信息
func GetPartitions(ctx context.Context, topic string) ([]int32, error) {
	if !IsReady() {
		return nil, errors.New("kafka client not initialized")
	}

	kafkaManager.mu.RLock()
	client := kafkaManager.client
	kafkaManager.mu.RUnlock()

	partitions, err := client.Partitions(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to get partitions for topic %s: %w", topic, err)
	}

	return partitions, nil
}

// === 统计和监控 ===

// GetKafkaStats 获取Kafka统计信息
func GetKafkaStats(ctx context.Context) (*KafkaStats, error) {
	if !IsReady() {
		return nil, errors.New("kafka client not initialized")
	}

	kafkaManager.mu.RLock()
	stats := kafkaManager.stats
	stats.IsConnected = kafkaManager.isReady && !kafkaManager.client.Closed()
	stats.LastPing = kafkaManager.lastPing
	kafkaManager.mu.RUnlock()

	return &stats, nil
}

// === 资源清理 ===

// CloseKafka 关闭Kafka连接
func CloseKafka() error {
	if kafkaManager == nil {
		return nil
	}

	kafkaManager.mu.Lock()
	defer kafkaManager.mu.Unlock()

	var errs []error

	if kafkaManager.producer != nil {
		if err := kafkaManager.producer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close producer: %w", err))
		}
		kafkaManager.producer = nil
	}

	if kafkaManager.consumer != nil {
		if err := kafkaManager.consumer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close consumer: %w", err))
		}
		kafkaManager.consumer = nil
	}

	if kafkaManager.client != nil {
		if err := kafkaManager.client.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close client: %w", err))
		}
		kafkaManager.client = nil
	}

	kafkaManager.isReady = false

	if len(errs) > 0 {
		return fmt.Errorf("errors closing kafka: %v", errs)
	}

	pkglog.LogInfo("Kafka connection closed successfully")
	return nil
}
