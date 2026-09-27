package proc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultShutdownConf(t *testing.T) {
	c := DefaultShutdownConf()
	assert.Equal(t, time.Second, c.WrapUpTime)
	assert.Equal(t, 5500*time.Millisecond, c.WaitTime)
}
