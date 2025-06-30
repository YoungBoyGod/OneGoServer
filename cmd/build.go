package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	// 构建平台
	targetOS     string
	targetArch   string
	outputDir    string
	crossCompile bool
	skipClean    bool
)

// buildCmd build子命令
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "构建OneGoServer二进制文件",
	Long: `构建OneGoServer二进制文件，支持多平台交叉编译。

该命令提供以下功能：
  • 本地平台构建
  • 交叉编译 (Windows/Linux/MacOS)
  • 版本信息注入
  • 输出目录自定义
  • 构建优化选项

构建过程会自动注入版本信息，包括：
  • Git版本标签
  • 构建时间
  • Git提交哈希
  • Git分支信息
  • Go版本信息
  • 构建者信息`,
	Example: `  # 构建本地平台版本
  onego build

  # 构建Windows版本
  onego build --os windows --arch amd64

  # 构建所有平台版本
  onego build --cross

  # 自定义输出目录
  onego build --output ./release

  # 跳过清理，保留之前的构建文件
  onego build --skip-clean`,
	Run: runBuild,
}

func runBuild(cmd *cobra.Command, args []string) {
	if verbose {
		fmt.Printf("🚀 开始构建 OneGoServer...\n")
	}

	// 获取构建信息
	buildInfo := getBuildInfo()

	if verbose {
		printBuildInfo(buildInfo)
	}

	// 清理旧文件
	if !skipClean {
		if err := cleanBuildDir(); err != nil {
			fmt.Printf("❌ 清理构建目录失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("❌ 创建输出目录失败: %v\n", err)
		os.Exit(1)
	}

	if crossCompile {
		// 交叉编译所有平台
		buildAllPlatforms(buildInfo)
	} else {
		// 构建指定平台
		if targetOS == "" {
			targetOS = runtime.GOOS
		}
		if targetArch == "" {
			targetArch = runtime.GOARCH
		}

		buildPlatform(targetOS, targetArch, buildInfo)
	}

	fmt.Printf("✅ 构建完成!\n")

	// 列出构建文件
	if verbose {
		listBuildFiles()
	}
}

// BuildInfo 构建信息
type BuildInfo struct {
	Version   string
	BuildTime string
	GitCommit string
	GitBranch string
	GoVersion string
	BuildBy   string
}

func getBuildInfo() BuildInfo {
	buildTime := time.Now().UTC().Format(time.RFC3339)

	// 获取Git信息
	gitCommit := getGitCommit()
	gitBranch := getGitBranch()

	// 获取Go版本
	goVersion := runtime.Version()

	// 获取构建者信息
	buildBy := getBuildBy()

	return BuildInfo{
		Version:   Version,
		BuildTime: buildTime,
		GitCommit: gitCommit,
		GitBranch: gitBranch,
		GoVersion: goVersion,
		BuildBy:   buildBy,
	}
}

func getGitCommit() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return string(output[:min(len(output)-1, 8)]) // 取前8位
}

func getGitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return string(output[:len(output)-1]) // 去掉换行符
}

func getBuildBy() string {
	username := os.Getenv("USERNAME")
	if username == "" {
		username = os.Getenv("USER")
	}
	if username == "" {
		username = "unknown"
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return fmt.Sprintf("%s@%s", username, hostname)
}

func printBuildInfo(info BuildInfo) {
	fmt.Printf("📋 构建信息:\n")
	fmt.Printf("  版本号:   %s\n", info.Version)
	fmt.Printf("  构建时间: %s\n", info.BuildTime)
	fmt.Printf("  Git提交:  %s\n", info.GitCommit)
	fmt.Printf("  Git分支:  %s\n", info.GitBranch)
	fmt.Printf("  Go版本:   %s\n", info.GoVersion)
	fmt.Printf("  构建者:   %s\n", info.BuildBy)
	fmt.Printf("  输出目录: %s\n", outputDir)
	fmt.Printf("\n")
}

func cleanBuildDir() error {
	if verbose {
		fmt.Printf("🧹 清理构建目录: %s\n", outputDir)
	}

	return os.RemoveAll(outputDir)
}

func buildPlatform(goos, goarch string, info BuildInfo) {
	binaryName := "onego"
	if goos == "windows" {
		binaryName += ".exe"
	}

	if crossCompile || (goos != runtime.GOOS || goarch != runtime.GOARCH) {
		binaryName = fmt.Sprintf("onego-%s-%s", goos, goarch)
		if goos == "windows" {
			binaryName += ".exe"
		}
	}

	outputPath := filepath.Join(outputDir, binaryName)

	if verbose {
		fmt.Printf("🔨 构建 %s/%s: %s\n", goos, goarch, outputPath)
	}

	// 构建LDFLAGS
	ldflags := fmt.Sprintf("-ldflags=-X 'github.com/YoungBoyGod/OneGoServer/cmd.Version=%s' "+
		"-X 'github.com/YoungBoyGod/OneGoServer/cmd.BuildTime=%s' "+
		"-X 'github.com/YoungBoyGod/OneGoServer/cmd.GitCommit=%s' "+
		"-X 'github.com/YoungBoyGod/OneGoServer/cmd.GitBranch=%s' "+
		"-X 'github.com/YoungBoyGod/OneGoServer/cmd.GoVersion=%s' "+
		"-X 'github.com/YoungBoyGod/OneGoServer/cmd.BuildBy=%s' "+
		"-w -s",
		info.Version, info.BuildTime, info.GitCommit,
		info.GitBranch, info.GoVersion, info.BuildBy)

	// 构建命令
	buildCmd := exec.Command("go", "build", ldflags, "-o", outputPath, "./main.go")
	buildCmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", goos),
		fmt.Sprintf("GOARCH=%s", goarch),
		"CGO_ENABLED=0",
	)

	if verbose {
		fmt.Printf("执行命令: %s\n", buildCmd.String())
	}

	output, err := buildCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ 构建失败 (%s/%s): %v\n", goos, goarch, err)
		if len(output) > 0 {
			fmt.Printf("输出: %s\n", string(output))
		}
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("✅ 构建成功: %s\n", outputPath)
	}
}

func buildAllPlatforms(info BuildInfo) {
	platforms := []struct {
		OS   string
		Arch string
	}{
		{"windows", "amd64"},
		{"linux", "amd64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
	}

	fmt.Printf("🌐 交叉编译所有平台...\n")

	for _, platform := range platforms {
		buildPlatform(platform.OS, platform.Arch, info)
	}
}

func listBuildFiles() {
	fmt.Printf("\n📋 构建文件列表:\n")

	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Dir(path) == outputDir {
			size := info.Size()
			fmt.Printf("  %s (%s)\n", info.Name(), formatFileSize(size))
		}

		return nil
	})

	if err != nil {
		fmt.Printf("⚠️ 列出文件失败: %v\n", err)
	}
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	// 添加build子命令到根命令
	rootCmd.AddCommand(buildCmd)

	// build命令的标志
	buildCmd.Flags().StringVar(&targetOS, "os", "", "目标操作系统 (windows/linux/darwin)")
	buildCmd.Flags().StringVar(&targetArch, "arch", "", "目标架构 (amd64/arm64)")
	buildCmd.Flags().StringVar(&outputDir, "output", "bin", "输出目录")
	buildCmd.Flags().BoolVar(&crossCompile, "cross", false, "交叉编译所有平台")
	buildCmd.Flags().BoolVar(&skipClean, "skip-clean", false, "跳过清理，保留之前的构建文件")

	// 设置标志说明
	buildCmd.Flags().Lookup("os").Usage = "指定目标操作系统"
	buildCmd.Flags().Lookup("arch").Usage = "指定目标架构"
	buildCmd.Flags().Lookup("output").Usage = "指定输出目录"
	buildCmd.Flags().Lookup("cross").Usage = "交叉编译Windows/Linux/MacOS所有平台"
	buildCmd.Flags().Lookup("skip-clean").Usage = "跳过清理步骤"
}
