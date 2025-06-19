package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"OneGoTask/pkg/config"

	"gopkg.in/yaml.v3"
)

// Client HTTP客户端结构体
type Client struct {
	config     *config.Config
	httpClient *http.Client
}

// New 创建新的客户端实例
func New(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Client.Timeout) * time.Second,
		},
	}
}

// Ping 测试服务器连接
func (c *Client) Ping() error {
	return c.makeRequest("GET", "/ping", nil, "Ping服务器")
}

// Health 获取服务器健康状态
func (c *Client) Health() error {
	return c.makeRequest("GET", "/health", nil, "检查服务器健康状态")
}

// Config 获取服务器配置
func (c *Client) Config() error {
	return c.makeRequest("GET", "/config", nil, "获取服务器配置")
}

// Status 获取服务器状态
func (c *Client) Status() error {
	return c.makeRequest("GET", "/api/v1/status", nil, "获取服务器状态")
}

// Version 获取服务器版本
func (c *Client) Version() error {
	return c.makeRequest("GET", "/api/v1/version", nil, "获取服务器版本")
}

// Restart 重启服务器
func (c *Client) Restart() error {
	return c.makeRequest("POST", "/restart", nil, "重启服务器")
}

// makeRequest 发起HTTP请求
func (c *Client) makeRequest(method, path string, body interface{}, description string) error {
	url := c.config.Client.ServerURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("序列化请求数据失败: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// 重试机制
	var lastErr error
	for i := 0; i <= c.config.Client.RetryCount; i++ {
		if i > 0 {
			fmt.Printf("⏳ 重试第 %d 次...\n", i)
			time.Sleep(time.Duration(c.config.Client.RetryDelay) * time.Second)
		}

		req, err := http.NewRequest(method, url, reqBody)
		if err != nil {
			lastErr = fmt.Errorf("创建请求失败: %v", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "OneGoTask-Client/1.0")

		fmt.Printf("🌐 %s %s\n", description, url)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("请求失败: %v", err)
			continue
		}

		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("读取响应失败: %v", err)
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return c.formatOutput(respBody)
		} else {
			lastErr = fmt.Errorf("服务器返回错误: %d %s", resp.StatusCode, string(respBody))
			continue
		}
	}

	return fmt.Errorf("请求失败，已重试 %d 次: %v", c.config.Client.RetryCount, lastErr)
}

// formatOutput 格式化输出
func (c *Client) formatOutput(data []byte) error {
	switch c.config.Client.OutputFormat {
	case "json":
		return c.outputJSON(data)
	case "yaml":
		return c.outputYAML(data)
	case "table":
		return c.outputTable(data)
	default:
		return c.outputJSON(data) // 默认JSON格式
	}
}

// outputJSON JSON格式输出
func (c *Client) outputJSON(data []byte) error {
	var formatted bytes.Buffer
	err := json.Indent(&formatted, data, "", "  ")
	if err != nil {
		return fmt.Errorf("格式化JSON失败: %v", err)
	}

	fmt.Printf("📄 响应结果 (JSON):\n")
	fmt.Println(formatted.String())
	return nil
}

// outputYAML YAML格式输出
func (c *Client) outputYAML(data []byte) error {
	var jsonData interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return fmt.Errorf("解析JSON失败: %v", err)
	}

	yamlData, err := yaml.Marshal(jsonData)
	if err != nil {
		return fmt.Errorf("转换为YAML失败: %v", err)
	}

	fmt.Printf("📄 响应结果 (YAML):\n")
	fmt.Print(string(yamlData))
	return nil
}

// outputTable 表格格式输出
func (c *Client) outputTable(data []byte) error {
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return fmt.Errorf("解析JSON失败: %v", err)
	}

	fmt.Printf("📊 响应结果 (表格):\n")
	fmt.Println("┌─────────────────────┬─────────────────────────────────────┐")
	fmt.Println("│        字段         │                 值                  │")
	fmt.Println("├─────────────────────┼─────────────────────────────────────┤")

	c.printTableRow(jsonData, "")

	fmt.Println("└─────────────────────┴─────────────────────────────────────┘")
	return nil
}

// printTableRow 打印表格行
func (c *Client) printTableRow(data interface{}, prefix string) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			fullKey := key
			if prefix != "" {
				fullKey = prefix + "." + key
			}

			switch val := value.(type) {
			case map[string]interface{}:
				c.printTableRow(val, fullKey)
			default:
				fmt.Printf("│ %-19s │ %-35v │\n", fullKey, val)
			}
		}
	default:
		fmt.Printf("│ %-19s │ %-35v │\n", prefix, v)
	}
}
