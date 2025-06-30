package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	// 版本信息
	Version   = "v0.1.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GitBranch = "unknown"
	GoVersion = runtime.Version()
	BuildBy   = "unknown"
)

// versionCmd version子命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示OneGoServer版本信息",
	Long: `显示OneGoServer的详细版本信息，包括：
  • 版本号
  • 构建时间
  • Git提交信息
  • Go版本信息
  • 构建环境信息`,
	Example: `  # 显示简单版本信息
  onego version

  # 显示详细版本信息  
  onego version --verbose

  # 只显示版本号
  onego --version`,
	Run: showVersion,
}

func showVersion(cmd *cobra.Command, args []string) {
	if verbose {
		// 详细版本信息
		fmt.Printf("OneGoServer 详细版本信息:\n")
		fmt.Printf("┌─────────────────────────────────────────────────────────────┐\n")
		fmt.Printf("│                    OneGoServer                              │\n")
		fmt.Printf("│                分布式任务管理平台                              │\n")
		fmt.Printf("├─────────────────────────────────────────────────────────────┤\n")
		fmt.Printf("│ 版本号        : %-40s │\n", Version)
		fmt.Printf("│ 构建时间      : %-40s │\n", getBuildTime())
		fmt.Printf("│ Git提交       : %-40s │\n", getShortCommit())
		fmt.Printf("│ Git分支       : %-40s │\n", GitBranch)
		fmt.Printf("│ Go版本        : %-40s │\n", GoVersion)
		fmt.Printf("│ 操作系统      : %-40s │\n", runtime.GOOS)
		fmt.Printf("│ 系统架构      : %-40s │\n", runtime.GOARCH)
		fmt.Printf("│ 构建者        : %-40s │\n", BuildBy)
		fmt.Printf("└─────────────────────────────────────────────────────────────┘\n")
		fmt.Printf("\n")
		fmt.Printf("🏗️  架构: DDD (领域驱动设计) + Clean Architecture\n")
		fmt.Printf("🚀 技术栈: Gin + PostgreSQL + Redis + Kafka + JWT\n")
		fmt.Printf("📚 项目地址: https://github.com/YoungBoyGod/OneGoServer\n")
	} else {
		// 简单版本信息
		fmt.Printf("OneGoServer %s\n", Version)
		fmt.Printf("构建时间: %s\n", getBuildTime())
		fmt.Printf("Go版本: %s\n", GoVersion)
	}
}

// getBuildTime 获取格式化的构建时间
func getBuildTime() string {
	if BuildTime == "unknown" {
		return "未知"
	}

	// 尝试解析构建时间
	if t, err := time.Parse(time.RFC3339, BuildTime); err == nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return BuildTime
}

// getShortCommit 获取短Git提交哈希
func getShortCommit() string {
	if GitCommit == "unknown" {
		return "未知"
	}

	if len(GitCommit) > 8 {
		return GitCommit[:8]
	}
	return GitCommit
}

func init() {
	// 添加version子命令到根命令
	rootCmd.AddCommand(versionCmd)
}
