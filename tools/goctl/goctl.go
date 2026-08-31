package main

import (
	"github.com/JellyGoFF/FF-Hexas/core/load"
	"github.com/JellyGoFF/FF-Hexas/core/logx"
	"github.com/JellyGoFF/FF-Hexas/tools/goctl/cmd"
)

func main() {
	logx.Disable()
	load.Disable()
	cmd.Execute()
}
