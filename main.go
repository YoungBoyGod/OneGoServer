package main

import (
	_ "OneGfServer/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"OneGfServer/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
