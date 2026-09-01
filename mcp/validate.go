package mcp

import "fmt"

// Validate checks MCP bootstrap configuration.
func (c McpConf) Validate() error {
	if err := c.RestConf.Validate(); err != nil { return err }
	if c.Mcp.Version == "" || c.Mcp.SseEndpoint == "" || c.Mcp.MessageEndpoint == "" { return fmt.Errorf("mcp version and endpoints are required") }
	if c.Mcp.SseTimeout <= 0 || c.Mcp.MessageTimeout <= 0 { return fmt.Errorf("mcp timeouts must be positive") }
	return nil
}
