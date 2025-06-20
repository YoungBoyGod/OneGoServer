package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// ClientCredential 客户端凭证结构
type ClientCredential struct {
	Type        string    `json:"type"`
	Namespace   string    `json:"namespace"`
	ClientName  string    `json:"client_name"`
	SecretKey   string    `json:"secret_key,omitempty"`
	Token       string    `json:"token,omitempty"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	Certificate string    `json:"certificate,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
}

var credentialCmd = &cobra.Command{
	Use:   "credential",
	Short: "客户端凭证管理",
	Long:  `管理客户端注册凭证，支持PSK、JWT、BasicAuth等认证方式`,
}

var generateCredentialCmd = &cobra.Command{
	Use:   "generate",
	Short: "生成客户端凭证",
	Long:  `生成不同类型的客户端注册凭证`,
	Run: func(cmd *cobra.Command, args []string) {
		credType, _ := cmd.Flags().GetString("type")
		namespace, _ := cmd.Flags().GetString("namespace")
		clientName, _ := cmd.Flags().GetString("client-name")
		output, _ := cmd.Flags().GetString("output")
		validDays, _ := cmd.Flags().GetInt("valid-days")

		if namespace == "" {
			namespace = "default"
		}
		if clientName == "" {
			clientName = "OneGoClient"
		}

		var credential *ClientCredential
		var err error

		switch credType {
		case "psk":
			credential, err = generatePSKCredential(namespace, clientName, validDays)
		case "basic":
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			credential, err = generateBasicCredential(namespace, clientName, username, password, validDays)
		case "jwt":
			credential, err = generateJWTCredential(namespace, clientName, validDays*24) // 转换为小时
		default:
			fmt.Printf("错误：不支持的凭证类型: %s\n", credType)
			fmt.Println("支持的类型: psk, basic, jwt")
			os.Exit(1)
		}

		if err != nil {
			fmt.Printf("生成凭证失败: %v\n", err)
			os.Exit(1)
		}

		// 保存到文件
		if output != "" {
			if err := saveCredentialToFile(credential, output); err != nil {
				fmt.Printf("保存凭证文件失败: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ 凭证已保存到: %s\n", output)
		}

		// 显示凭证信息
		printCredentialInfo(credential)
	},
}

func generatePSKCredential(namespace, clientName string, validDays int) (*ClientCredential, error) {
	secretKey := "your-psk-secret-key-here" // 在实际环境中应该生成随机密钥

	credential := &ClientCredential{
		Type:       "psk",
		Namespace:  namespace,
		ClientName: clientName,
		SecretKey:  secretKey,
		CreatedAt:  time.Now(),
		CreatedBy:  "admin",
	}

	if validDays > 0 {
		credential.ExpiresAt = time.Now().AddDate(0, 0, validDays)
	}

	return credential, nil
}

func generateBasicCredential(namespace, clientName, username, password string, validDays int) (*ClientCredential, error) {
	if username == "" {
		username = "onegoclient"
	}
	if password == "" {
		password = "onegoclient123"
	}

	credential := &ClientCredential{
		Type:       "basic",
		Namespace:  namespace,
		ClientName: clientName,
		Username:   username,
		Password:   password,
		CreatedAt:  time.Now(),
		CreatedBy:  "admin",
	}

	if validDays > 0 {
		credential.ExpiresAt = time.Now().AddDate(0, 0, validDays)
	}

	return credential, nil
}

func generateJWTCredential(namespace, clientName string, validHours int) (*ClientCredential, error) {
	// 简化的JWT Token（实际环境应使用真正的JWT库）
	token := fmt.Sprintf("jwt-token-%s-%s-%d", namespace, clientName, time.Now().Unix())

	credential := &ClientCredential{
		Type:       "jwt",
		Namespace:  namespace,
		ClientName: clientName,
		Token:      token,
		CreatedAt:  time.Now(),
		CreatedBy:  "admin",
	}

	if validHours > 0 {
		credential.ExpiresAt = time.Now().Add(time.Duration(validHours) * time.Hour)
	}

	return credential, nil
}

func saveCredentialToFile(credential *ClientCredential, filename string) error {
	data, err := json.MarshalIndent(credential, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0600)
}

func printCredentialInfo(credential *ClientCredential) {
	fmt.Println("\n📋 生成的客户端凭证信息:")
	fmt.Printf("类型: %s\n", credential.Type)
	fmt.Printf("命名空间: %s\n", credential.Namespace)
	fmt.Printf("客户端名称: %s\n", credential.ClientName)
	fmt.Printf("创建时间: %s\n", credential.CreatedAt.Format("2006-01-02 15:04:05"))

	if !credential.ExpiresAt.IsZero() {
		fmt.Printf("过期时间: %s\n", credential.ExpiresAt.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Println("过期时间: 永不过期")
	}

	switch credential.Type {
	case "psk":
		fmt.Printf("密钥: %s\n", credential.SecretKey)
	case "basic":
		fmt.Printf("用户名: %s\n", credential.Username)
		fmt.Printf("密码: %s\n", credential.Password)
	case "jwt":
		fmt.Printf("Token: %s\n", credential.Token)
	}

	fmt.Println("\n🔐 使用说明:")
	fmt.Println("1. 将生成的凭证文件放置在客户端的 config/ 目录下")
	fmt.Printf("2. 文件名建议为: client_credential.json\n")
	fmt.Println("3. 客户端将自动检测并使用该凭证进行认证")
}

func init() {
	rootCmd.AddCommand(credentialCmd)
	credentialCmd.AddCommand(generateCredentialCmd)

	// 凭证生成参数
	generateCredentialCmd.Flags().StringP("type", "t", "psk", "凭证类型 (psk, basic, jwt)")
	generateCredentialCmd.Flags().StringP("namespace", "n", "default", "命名空间")
	generateCredentialCmd.Flags().StringP("client-name", "", "OneGoClient", "客户端名称")
	generateCredentialCmd.Flags().StringP("output", "o", "client_credential.json", "输出文件路径")
	generateCredentialCmd.Flags().IntP("valid-days", "d", 0, "有效天数 (0表示永不过期)")

	// BasicAuth特有参数
	generateCredentialCmd.Flags().String("username", "", "用户名 (仅用于basic类型)")
	generateCredentialCmd.Flags().String("password", "", "密码 (仅用于basic类型)")
}
