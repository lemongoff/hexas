package mcp

import (
	"time"

	"github.com/lemongoff/hexas/rest"
)

// DefaultMcpConf returns Hexas MCP server defaults.
func DefaultMcpConf() McpConf {
	configuration := McpConf{RestConf: rest.DefaultRestConf()}
	configuration.Mcp.Version = "1.0.0"
	configuration.Mcp.SseEndpoint = "/sse"
	configuration.Mcp.MessageEndpoint = "/message"
	configuration.Mcp.SseTimeout = 24 * time.Hour
	configuration.Mcp.MessageTimeout = 30 * time.Second
	return configuration
}
