package main

import (
	"github.com/lemongoff/hexas/core/load"
	"github.com/lemongoff/hexas/core/logx"
	"github.com/lemongoff/hexas/tools/goctl/cmd"
)

func main() {
	logx.Disable()
	load.Disable()
	cmd.Execute()
}
