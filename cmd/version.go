package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// versionCmd 代表version命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long: `显示应用程序的详细版本信息，包括:

• 应用版本号
• Git提交哈希
• 构建时间
• Go版本信息
• 编译平台信息`,
	Run: showVersion,
}

func showVersion(cmd *cobra.Command, args []string) {
	short, _ := cmd.Flags().GetBool("short")

	if short {
		fmt.Printf("%s\n", version)
		return
	}

	fmt.Printf("📦 %s\n", rootCmd.Short)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("版本:     %s\n", version)
	fmt.Printf("提交:     %s\n", commit)
	fmt.Printf("构建时间:  %s\n", date)
	fmt.Printf("Go版本:   %s\n", runtime.Version())
	fmt.Printf("平台:     %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("编译器:   %s\n", runtime.Compiler)
}

func init() {
	// 添加version特定的标志
	versionCmd.Flags().BoolP("short", "s", false, "只显示版本号")

	// 将version命令添加到根命令
	rootCmd.AddCommand(versionCmd)
}
