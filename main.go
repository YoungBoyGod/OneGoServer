package main

import (
	"OneGoTask/cmd"
)

// main 程序入口点
// 只负责启动Cobra命令行处理器
func main() {
	// 执行根命令，命令处理逻辑已经移到cmd包中
	cmd.Execute()
}
