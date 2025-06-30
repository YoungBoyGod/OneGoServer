package queue

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	pkglog "github.com/YoungBoyGod/OneGoServer/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// === 测试配置和帮助函数 ===

// getTestKafkaConfig 获取测试用的Kafka配置
func getTestKafkaConfig() *config.KafkaConfig {
	// 支持多种环境的Kafka配置
	kafkaHost := os.Getenv("KAFKA_HOST")
	if kafkaHost == "" {
		kafkaHost = "localhost" // 默认使用localhost
	}

	kafkaPort := os.Getenv("KAFKA_PORT")
	if kafkaPort == "" {
		kafkaPort = "9092" // 默认端口
	}

	// 支持多个broker地址，用于容错
	brokers := []string{
		fmt.Sprintf("%s:%s", kafkaHost, kafkaPort),
	}

	// 如果是localhost，也尝试kafka主机名（用于Docker环境）
	if kafkaHost == "localhost" {
		brokers = append(brokers, fmt.Sprintf("kafka:%s", kafkaPort))
	}

	return &config.KafkaConfig{
		Brokers:                   brokers,
		ClientID:                  "onegoserver-test-client",
		Version:                   "2.6.0",
		Username:                  "",
		Password:                  "",
		EnableSASL:                false,
		SASLMechanism:             "PLAIN",
		EnableTLS:                 false,
		ProducerReturnSuccesses:   true,
		ProducerReturnErrors:      true,
		ProducerRequiredAcks:      1,
		ProducerRetryMax:          3,
		ProducerMaxMessageBytes:   1000000,
		ConsumerGroupID:           "test-consumer-group",
		ConsumerOffsetInitial:     "newest",
		ConsumerSessionTimeout:    10000,
		ConsumerHeartbeatInterval: 3000,
	}
}

// getTestLoggingConfig 获取测试用的日志配置
func getTestLoggingConfig() *config.LoggingConfig {
	return &config.LoggingConfig{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		FilePath:   "logs/test.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}
}

// setupTest 测试初始化
func setupTest(t *testing.T) {
	// 初始化日志系统
	cfg := getTestLoggingConfig()
	err := pkglog.InitLoggerEnhanced(cfg)
	require.NoError(t, err)

	// 重置全局状态
	kafkaManager = nil
	once = sync.Once{}
}

// === 配置验证测试 ===

func TestValidateKafkaConfig(t *testing.T) {
	setupTest(t)

	t.Run("有效配置", func(t *testing.T) {
		cfg := getTestKafkaConfig()
		err := validateKafkaConfig(cfg)
		assert.NoError(t, err)
	})

	t.Run("空配置", func(t *testing.T) {
		err := validateKafkaConfig(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kafka config cannot be nil")
	})

	t.Run("空Brokers列表", func(t *testing.T) {
		cfg := getTestKafkaConfig()
		cfg.Brokers = []string{}
		err := validateKafkaConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "brokers list cannot be empty")
	})

	t.Run("空Broker地址", func(t *testing.T) {
		cfg := getTestKafkaConfig()
		cfg.Brokers = []string{"localhost:9092", "", "127.0.0.1:9092"}
		err := validateKafkaConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "broker address cannot be empty")
	})

	t.Run("空ClientID", func(t *testing.T) {
		cfg := getTestKafkaConfig()
		cfg.ClientID = ""
		err := validateKafkaConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client id is required")
	})

	t.Run("无效版本", func(t *testing.T) {
		cfg := getTestKafkaConfig()
		cfg.Version = "invalid.version"
		err := validateKafkaConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid kafka version")
	})

	t.Run("默认值设置", func(t *testing.T) {
		cfg := &config.KafkaConfig{
			Brokers:  []string{"localhost:9092"},
			ClientID: "test-client",
		}

		err := validateKafkaConfig(cfg)
		assert.NoError(t, err)

		// 验证默认值
		assert.Equal(t, "2.6.0", cfg.Version)
		assert.Equal(t, 1, cfg.ProducerRequiredAcks)
		assert.Equal(t, 3, cfg.ProducerRetryMax)
		assert.Equal(t, 1000000, cfg.ProducerMaxMessageBytes)
		assert.Equal(t, "onegoserver-consumer-group", cfg.ConsumerGroupID)
		assert.Equal(t, "newest", cfg.ConsumerOffsetInitial)
		assert.Equal(t, 10000, cfg.ConsumerSessionTimeout)
		assert.Equal(t, 3000, cfg.ConsumerHeartbeatInterval)
	})
}

// === 初始化测试 ===

func TestKafkaManager(t *testing.T) {
	setupTest(t)

	t.Run("Kafka管理器初始化状态", func(t *testing.T) {
		// 检查未初始化状态
		assert.False(t, IsReady())
		assert.Nil(t, GetKafkaClient())
	})

	t.Run("配置验证", func(t *testing.T) {
		cfg := getTestKafkaConfig()
		err := validateKafkaConfig(cfg)
		assert.NoError(t, err)
	})
}

// === 无连接操作测试 ===

func TestKafkaOperationsWithoutConnection(t *testing.T) {
	setupTest(t)

	ctx := context.Background()

	t.Run("Ping操作无连接", func(t *testing.T) {
		err := Ping(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kafka client not initialized")
	})

	t.Run("SendMessage操作无连接", func(t *testing.T) {
		message := &KafkaMessage{
			Topic: "test-topic",
			Value: "test message",
		}
		err := SendMessage(ctx, message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kafka client not initialized")
	})

	t.Run("SendMessages操作无连接", func(t *testing.T) {
		messages := []*KafkaMessage{
			{Topic: "test-topic", Value: "message1"},
			{Topic: "test-topic", Value: "message2"},
		}
		err := SendMessages(ctx, messages)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kafka client not initialized")
	})

	t.Run("ConsumeMessages操作无连接", func(t *testing.T) {
		handler := func(ctx context.Context, message *KafkaMessage) error {
			return nil
		}

		err := ConsumeMessages(ctx, "test-topic", 0, 0, handler)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kafka client not initialized")
	})

	t.Run("GetTopics操作无连接", func(t *testing.T) {
		_, err := GetTopics(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kafka client not initialized")
	})

	t.Run("GetPartitions操作无连接", func(t *testing.T) {
		_, err := GetPartitions(ctx, "test-topic")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kafka client not initialized")
	})

	t.Run("GetKafkaStats操作无连接", func(t *testing.T) {
		_, err := GetKafkaStats(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kafka client not initialized")
	})
}

// === 消息类型测试 ===

func TestKafkaMessage(t *testing.T) {
	t.Run("消息结构验证", func(t *testing.T) {
		message := &KafkaMessage{
			Topic:     "test-topic",
			Key:       "test-key",
			Value:     "test-value",
			Headers:   map[string]string{"Content-Type": "application/json"},
			Partition: 1,
			Offset:    100,
			Timestamp: time.Now(),
		}

		assert.Equal(t, "test-topic", message.Topic)
		assert.Equal(t, "test-key", message.Key)
		assert.Equal(t, "test-value", message.Value)
		assert.Equal(t, "application/json", message.Headers["Content-Type"])
		assert.Equal(t, int32(1), message.Partition)
		assert.Equal(t, int64(100), message.Offset)
	})

	t.Run("复杂对象消息", func(t *testing.T) {
		complexValue := map[string]interface{}{
			"id":     12345,
			"name":   "测试对象",
			"active": true,
			"tags":   []string{"tag1", "tag2"},
			"metadata": map[string]interface{}{
				"created_at": time.Now().Format(time.RFC3339),
				"version":    "1.0.0",
			},
		}

		message := &KafkaMessage{
			Topic: "complex-topic",
			Value: complexValue,
		}

		assert.Equal(t, "complex-topic", message.Topic)
		assert.NotNil(t, message.Value)
	})
}

// === 错误处理测试 ===

func TestKafkaErrorHandling(t *testing.T) {
	setupTest(t)

	t.Run("配置错误处理", func(t *testing.T) {
		invalidCfg := &config.KafkaConfig{
			Brokers:  []string{},
			ClientID: "",
		}

		err := validateKafkaConfig(invalidCfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "brokers list cannot be empty")
	})

	t.Run("多重错误验证", func(t *testing.T) {
		invalidCfg := &config.KafkaConfig{
			Brokers:  []string{"", "localhost:9092"},
			ClientID: "",
			Version:  "invalid.version",
		}

		err := validateKafkaConfig(invalidCfg)
		assert.Error(t, err)
		// 验证会返回第一个遇到的错误
		assert.Contains(t, err.Error(), "broker address cannot be empty")
	})
}

// === 重试机制测试 ===

func TestKafkaRetryMechanism(t *testing.T) {
	setupTest(t)

	t.Run("重试配置验证", func(t *testing.T) {
		// 重置状态确保独立测试
		kafkaManager = nil
		once = sync.Once{}

		cfg := getTestKafkaConfig()
		cfg.Brokers = []string{"invalid-host:9092"} // 使用无效主机

		ctx := context.Background()

		err := InitKafkaWithRetry(ctx, cfg, 2, 100*time.Millisecond)
		assert.Error(t, err)
		// 检查错误消息包含重试信息
		assert.True(t,
			strings.Contains(err.Error(), "failed to connect to kafka after 2 attempts") ||
				strings.Contains(err.Error(), "context deadline exceeded") ||
				strings.Contains(err.Error(), "failed to connect to kafka"))
	})

	t.Run("配置错误重试", func(t *testing.T) {
		// 重置状态确保独立测试
		kafkaManager = nil
		once = sync.Once{}

		invalidCfg := &config.KafkaConfig{
			Brokers: []string{}, // 空的brokers列表
		}

		ctx := context.Background()
		err := InitKafkaWithRetry(ctx, invalidCfg, 2, 100*time.Millisecond)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid kafka config")
	})
}

// === 并发安全测试 ===

func TestKafkaManagerConcurrency(t *testing.T) {
	setupTest(t)

	t.Run("并发IsReady检查", func(t *testing.T) {
		kafkaManager = &KafkaManager{
			isReady: true,
		}

		// 启动多个goroutine并发访问
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					IsReady() // 并发读取
				}
				done <- true
			}()
		}

		// 等待所有goroutine完成
		for i := 0; i < 10; i++ {
			<-done
		}

		assert.True(t, IsReady()) // 确保状态正确
	})

	t.Run("并发获取客户端", func(t *testing.T) {
		kafkaManager = &KafkaManager{
			client:  nil,
			isReady: true,
		}

		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					GetKafkaClient() // 并发获取客户端
				}
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		assert.Nil(t, GetKafkaClient()) // 客户端为nil
	})
}

// === 统计信息测试 ===

func TestKafkaStats(t *testing.T) {
	t.Run("统计信息结构验证", func(t *testing.T) {
		stats := &KafkaStats{
			ProducerStats: ProducerStats{
				MessagesSent:      100,
				MessagesSucceeded: 95,
				MessagesFailed:    5,
				BytesSent:         10240,
			},
			ConsumerStats: ConsumerStats{
				MessagesReceived:  200,
				MessagesProcessed: 190,
				MessagesFailed:    10,
				BytesReceived:     20480,
			},
			IsConnected: true,
			LastPing:    time.Now(),
		}

		assert.Equal(t, int64(100), stats.ProducerStats.MessagesSent)
		assert.Equal(t, int64(95), stats.ProducerStats.MessagesSucceeded)
		assert.Equal(t, int64(5), stats.ProducerStats.MessagesFailed)
		assert.Equal(t, int64(10240), stats.ProducerStats.BytesSent)

		assert.Equal(t, int64(200), stats.ConsumerStats.MessagesReceived)
		assert.Equal(t, int64(190), stats.ConsumerStats.MessagesProcessed)
		assert.Equal(t, int64(10), stats.ConsumerStats.MessagesFailed)
		assert.Equal(t, int64(20480), stats.ConsumerStats.BytesReceived)

		assert.True(t, stats.IsConnected)
	})
}

// === 关闭测试 ===

func TestCloseKafka(t *testing.T) {
	setupTest(t)

	t.Run("关闭不存在的连接", func(t *testing.T) {
		kafkaManager = nil
		err := CloseKafka()
		assert.NoError(t, err) // 应该优雅处理
	})

	t.Run("关闭存在的管理器", func(t *testing.T) {
		kafkaManager = &KafkaManager{
			client:   nil,
			producer: nil,
			consumer: nil,
			isReady:  true,
		}

		err := CloseKafka()
		assert.NoError(t, err)
		assert.False(t, kafkaManager.isReady)
	})
}

// === Context测试 ===

func TestContextHandling(t *testing.T) {
	setupTest(t)

	t.Run("Context超时处理", func(t *testing.T) {
		// 创建一个已经取消的context
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消

		err := Ping(ctx)
		assert.Error(t, err)
		// 由于Kafka客户端没有初始化，会先返回初始化错误
		assert.Contains(t, err.Error(), "kafka client not initialized")
	})

	t.Run("Context带值传递", func(t *testing.T) {
		// 测试context可以携带值
		ctx := context.WithValue(context.Background(), "operation_id", "test_op_kafka_123")

		// 验证值存在
		value := ctx.Value("operation_id")
		assert.Equal(t, "test_op_kafka_123", value)

		// 尝试Kafka操作（预期失败因为没有连接）
		message := &KafkaMessage{
			Topic: "test-topic",
			Value: "test message",
		}
		err := SendMessage(ctx, message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kafka client not initialized")
	})

	t.Run("Context超时设置", func(t *testing.T) {
		// 创建带超时的context
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		message := &KafkaMessage{
			Topic: "test-topic",
			Value: "test message",
		}
		err := SendMessage(ctx, message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kafka client not initialized")
	})
}

// === 配置边界情况测试 ===

func TestKafkaConfigEdgeCases(t *testing.T) {
	setupTest(t)

	t.Run("最小有效配置", func(t *testing.T) {
		cfg := &config.KafkaConfig{
			Brokers:  []string{"localhost:9092"},
			ClientID: "minimal-client",
		}
		err := validateKafkaConfig(cfg)
		assert.NoError(t, err)
	})

	t.Run("多Broker配置", func(t *testing.T) {
		cfg := &config.KafkaConfig{
			Brokers:  []string{"broker1:9092", "broker2:9092", "broker3:9092"},
			ClientID: "multi-broker-client",
		}
		err := validateKafkaConfig(cfg)
		assert.NoError(t, err)
	})

	t.Run("SASL配置", func(t *testing.T) {
		cfg := getTestKafkaConfig()
		cfg.EnableSASL = true
		cfg.Username = "test-user"
		cfg.Password = "test-password"
		cfg.SASLMechanism = "PLAIN"
		err := validateKafkaConfig(cfg)
		assert.NoError(t, err)
	})

	t.Run("TLS配置", func(t *testing.T) {
		cfg := getTestKafkaConfig()
		cfg.EnableTLS = true
		err := validateKafkaConfig(cfg)
		assert.NoError(t, err)
	})
}

// === 基准测试 ===

// BenchmarkKafkaConfigValidation Kafka配置验证性能基准测试
func BenchmarkKafkaConfigValidation(b *testing.B) {
	cfg := getTestKafkaConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateKafkaConfig(cfg)
	}
}

// BenchmarkKafkaMessageCreation 消息创建性能基准测试
func BenchmarkKafkaMessageCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		message := &KafkaMessage{
			Topic: fmt.Sprintf("benchmark-topic-%d", i),
			Key:   fmt.Sprintf("key-%d", i),
			Value: fmt.Sprintf("This is benchmark message number %d", i),
			Headers: map[string]string{
				"Content-Type": "text/plain",
				"Source":       "benchmark",
			},
		}
		_ = message // 防止编译器优化
	}
}

// === 真实Kafka集成测试 ===

// TestKafkaIntegration 真实Kafka连接和操作测试
func TestKafkaIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试：使用 -short 标志")
	}

	setupTest(t)

	// 重置状态确保独立测试
	kafkaManager = nil
	once = sync.Once{}

	// 尝试连接真实Kafka
	kafkaCfg := getTestKafkaConfig()
	ctx := context.Background()

	err := InitKafka(ctx, kafkaCfg)
	if err != nil {
		t.Skipf("跳过Kafka集成测试：无法连接Kafka服务器 (%v)", err)
		return
	}

	// 确保测试结束后清理
	defer CloseKafka()

	t.Run("真实连接验证", func(t *testing.T) {
		assert.True(t, IsReady())
		assert.NotNil(t, GetKafkaClient())

		err := Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("发送单条消息", func(t *testing.T) {
		message := &KafkaMessage{
			Topic: "test-integration-topic",
			Key:   "integration-test-key",
			Value: "Hello Kafka Integration Test!",
			Headers: map[string]string{
				"Test-Type": "integration",
				"Source":    "onegoserver",
				"Timestamp": time.Now().Format(time.RFC3339),
			},
		}

		err := SendMessage(ctx, message)
		assert.NoError(t, err)
		assert.NotZero(t, message.Partition)
		assert.NotZero(t, message.Offset)
		assert.False(t, message.Timestamp.IsZero())
	})

	t.Run("批量发送消息", func(t *testing.T) {
		messages := make([]*KafkaMessage, 5)
		for i := 0; i < 5; i++ {
			messages[i] = &KafkaMessage{
				Topic: "test-batch-topic",
				Key:   fmt.Sprintf("batch-key-%d", i),
				Value: fmt.Sprintf("Batch message %d", i),
				Headers: map[string]string{
					"Batch-Index": fmt.Sprintf("%d", i),
					"Test-Type":   "batch",
				},
			}
		}

		err := SendMessages(ctx, messages)
		assert.NoError(t, err)

		// 验证所有消息都有分区和偏移量
		for i, msg := range messages {
			assert.NotZero(t, msg.Partition, "Message %d should have partition", i)
			assert.NotZero(t, msg.Offset, "Message %d should have offset", i)
		}
	})

	t.Run("获取主题列表", func(t *testing.T) {
		topics, err := GetTopics(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, topics)
		t.Logf("可用主题: %v", topics)
	})

	t.Run("获取分区信息", func(t *testing.T) {
		// 首先发送一条消息确保主题存在
		testTopic := "test-partition-topic"
		message := &KafkaMessage{
			Topic: testTopic,
			Value: "创建主题的测试消息",
		}
		err := SendMessage(ctx, message)
		assert.NoError(t, err)

		// 获取分区信息
		partitions, err := GetPartitions(ctx, testTopic)
		assert.NoError(t, err)
		assert.NotEmpty(t, partitions)
		t.Logf("主题 %s 的分区: %v", testTopic, partitions)
	})

	t.Run("获取统计信息", func(t *testing.T) {
		stats, err := GetKafkaStats(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.True(t, stats.IsConnected)
		assert.True(t, stats.ProducerStats.MessagesSent > 0)
		t.Logf("Kafka统计信息: %+v", stats)
	})
}

// TestKafkaRetryWithRealConnection 测试真实连接的重试机制
func TestKafkaRetryWithRealConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试：使用 -short 标志")
	}

	setupTest(t)

	t.Run("无效主机重试测试", func(t *testing.T) {
		// 重置状态确保独立测试
		kafkaManager = nil
		once = sync.Once{}

		invalidCfg := &config.KafkaConfig{
			Brokers:                   []string{"invalid-kafka-host-12345:9092"},
			ClientID:                  "retry-test-client",
			Version:                   "2.6.0",
			ProducerReturnSuccesses:   true,
			ProducerReturnErrors:      true,
			ProducerRequiredAcks:      1,
			ProducerRetryMax:          1,
			ProducerMaxMessageBytes:   1000000,
			ConsumerGroupID:           "test-group",
			ConsumerOffsetInitial:     "newest",
			ConsumerSessionTimeout:    5000, // 短超时
			ConsumerHeartbeatInterval: 1000,
		}

		ctx := context.Background()
		startTime := time.Now()

		err := InitKafkaWithRetry(ctx, invalidCfg, 3, 100*time.Millisecond)

		duration := time.Since(startTime)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to kafka after 3 attempts")

		// 验证重试确实花费了时间（至少2次重试间隔）
		expectedMinDuration := 200 * time.Millisecond
		assert.True(t, duration >= expectedMinDuration,
			"重试应该花费至少 %v，实际花费 %v", expectedMinDuration, duration)

		t.Logf("重试测试完成，耗时: %v", duration)
	})
}

// === 性能基准测试 ===

// BenchmarkKafkaOperations Kafka操作性能基准测试
func BenchmarkKafkaOperations(b *testing.B) {
	if testing.Short() {
		b.Skip("跳过性能测试：使用 -short 标志")
	}

	// 初始化日志系统
	cfg := getTestLoggingConfig()
	err := pkglog.InitLoggerEnhanced(cfg)
	if err != nil {
		b.Skipf("无法初始化日志系统: %v", err)
		return
	}

	// 重置状态确保独立测试
	kafkaManager = nil
	once = sync.Once{}

	kafkaCfg := getTestKafkaConfig()
	ctx := context.Background()

	err = InitKafka(ctx, kafkaCfg)
	if err != nil {
		b.Skipf("跳过性能测试：无法连接Kafka服务器 (%v)", err)
		return
	}

	defer CloseKafka()

	b.Run("SendMessage操作", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			message := &KafkaMessage{
				Topic: "benchmark-topic",
				Key:   fmt.Sprintf("bench-key-%d", i),
				Value: fmt.Sprintf("Benchmark message %d", i),
			}
			SendMessage(ctx, message)
		}
	})

	b.Run("消息创建", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			message := &KafkaMessage{
				Topic: "benchmark-topic",
				Key:   fmt.Sprintf("bench-key-%d", i),
				Value: map[string]interface{}{
					"id":        i,
					"message":   fmt.Sprintf("Benchmark object message %d", i),
					"timestamp": time.Now(),
				},
				Headers: map[string]string{
					"Content-Type": "application/json",
					"Source":       "benchmark",
				},
			}
			_ = message
		}
	})

	b.Run("并发发送消息", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			counter := 0
			for pb.Next() {
				message := &KafkaMessage{
					Topic: "benchmark-parallel-topic",
					Key:   fmt.Sprintf("parallel-key-%d", counter),
					Value: fmt.Sprintf("Parallel message %d", counter),
				}
				SendMessage(ctx, message)
				counter++
			}
		})
	})
}
