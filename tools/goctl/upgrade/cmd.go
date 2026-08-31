package upgrade

import "github.com/JellyGoFF/FF-Hexas/tools/goctl/internal/cobrax"

// Cmd describes an upgrade command.
var Cmd = cobrax.NewCommand("upgrade", cobrax.WithRunE(upgrade))
