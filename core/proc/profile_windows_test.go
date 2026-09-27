package proc

import (
	"testing"

	"github.com/lemongoff/hexas/core/logx/logtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfile(t *testing.T) {
	c := logtest.NewCollector(t)
	first := StartProfile()
	second := StartProfile()
	require.NotNil(t, first)
	require.NotNil(t, second)
	assert.IsType(t, nilStopper{}, first)
	assert.IsType(t, nilStopper{}, second)
	assert.NotPanics(t, func() {
		first.Stop()
		first.Stop()
		second.Stop()
		second.Stop()
	})
	assert.NotContains(t, c.String(), ".pprof")
}
