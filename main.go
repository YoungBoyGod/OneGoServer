package main

import (
	_ "OneGfServer/internal/packed"

	_ "OneGfServer/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"OneGfServer/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
