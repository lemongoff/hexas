package mcp

import (
	"testing"
	"time"

	"github.com/lemongoff/hexas/internal/configtest"
	"github.com/stretchr/testify/assert"
)

func TestMcpConfDefaults(t *testing.T) {
	// Test default values are set correctly
	jsonConfig := `name: test-service
port: 8080
mcp:
  name: test-mcp-server
  version: 1.0.0
`

	c := DefaultMcpConf()
	err := configtest.LoadYAML([]byte(jsonConfig), &c)
	assert.NoError(t, err)

	// Check default values
	assert.Equal(t, "test-mcp-server", c.Mcp.Name)
	assert.Equal(t, "1.0.0", c.Mcp.Version)
	assert.Equal(t, "/sse", c.Mcp.SseEndpoint)
	assert.Equal(t, "/message", c.Mcp.MessageEndpoint)
	assert.Equal(t, 30*time.Second, c.Mcp.MessageTimeout)
}
