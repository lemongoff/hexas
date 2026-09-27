package profiling

import (
	"context"
	"math"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/lemongoff/hexas/core/syncx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	t.Run("profiling", func(t *testing.T) {
		c := DefaultConfig()
		c.Name = "test"
		p := newProfiler(c)
		assert.NotNil(t, p)
		assert.NoError(t, p.Start())
		assert.NoError(t, p.Stop())
	})

	t.Run("invalid config", func(t *testing.T) {
		mp := newMockProfiler()
		setTestProfiler(t, mp)
		Start(Config{})
		assert.False(t, mp.started.True())
	})

	t.Run("test start profiler", func(t *testing.T) {
		mp := newMockProfiler()
		setTestProfiler(t, mp)

		c := Config{
			Name:              "test",
			ServerAddr:        "localhost:4040",
			CheckInterval:     time.Millisecond,
			ProfilingDuration: time.Millisecond * 10,
			CpuThreshold:      0,
		}
		stop := runTestProfiler(t, c)
		waitProfilerSignal(t, mp.startCalled, "profiler start")
		waitProfilerSignal(t, mp.stopCalled, "profiler stop")
		stop()

		assert.True(t, mp.started.True())
		assert.True(t, mp.stopped.True())
	})

	t.Run("cpu below threshold", func(t *testing.T) {
		mp := newMockProfiler()
		setTestProfiler(t, mp)

		c := Config{
			Name:              "test",
			ServerAddr:        "localhost:4040",
			CheckInterval:     time.Millisecond,
			ProfilingDuration: time.Millisecond * 10,
			CpuThreshold:      math.MaxInt64,
		}
		stop := runTestProfiler(t, c)
		// This is an observation window for the absence of a start, not a
		// deadline by which an asynchronous operation must have completed.
		timer := time.NewTimer(50 * time.Millisecond)
		defer timer.Stop()
		select {
		case <-mp.startCalled:
			t.Fatal("profiler started below the CPU threshold")
		case <-timer.C:
		}
		stop()

		assert.False(t, mp.started.True())
	})

	for _, startFails := range []bool{true, false} {
		name := "stop error"
		if startFails {
			name = "start error"
		}
		t.Run(name, func(t *testing.T) {
			mp := newMockProfiler()
			mp.stopErr = assert.AnError
			if startFails {
				mp.startErr = assert.AnError
			}
			setTestProfiler(t, mp)
			stop := runTestProfiler(t, Config{
				Name: "test", CheckInterval: time.Millisecond,
				ProfilingDuration: 10 * time.Millisecond,
			})
			waitProfilerSignal(t, mp.startCalled, "profiler start attempt")
			if !startFails {
				waitProfilerSignal(t, mp.stopCalled, "profiler stop attempt")
			}
			stop()
			assert.Equal(t, !startFails, mp.started.True())
			assert.False(t, mp.stopped.True())
		})
	}
}

func TestStartProcess(t *testing.T) {
	const childEnv = "HEXAS_TEST_PROFILING_START_PROCESS"
	if os.Getenv(childEnv) == "1" {
		// Start has a process-lifetime worker and sync.Once. Keep both in a
		// child process so they cannot leak into other tests or repeated runs.
		mp := newMockProfiler()
		newProfiler = func(Config) profiler { return mp }
		Start(Config{
			ServerAddr: "localhost:4040", CheckInterval: time.Millisecond,
			ProfilingDuration: 10 * time.Millisecond,
		})
		waitProfilerSignal(t, mp.startCalled, "public Start")
		waitProfilerSignal(t, mp.stopCalled, "public Start profiling cycle")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestStartProcess$", "-test.count=1")
	cmd.Env = append(os.Environ(), childEnv+"=1")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", output)
}

func setTestProfiler(t *testing.T, mp *mockProfiler) {
	t.Helper()
	previous := newProfiler
	newProfiler = func(Config) profiler { return mp }
	t.Cleanup(func() { newProfiler = previous })
}

func runTestProfiler(t *testing.T, c Config) func() {
	t.Helper()
	done := make(chan struct{})
	exited := make(chan struct{})
	go func() {
		defer close(exited)
		startPyroscope(c, done)
	}()
	stop := sync.OnceFunc(func() {
		close(done)
		waitProfilerSignal(t, exited, "worker exit")
	})
	// Registered after setTestProfiler, so the worker exits before restoration.
	t.Cleanup(stop)
	return stop
}

func waitProfilerSignal(t *testing.T, signal <-chan struct{}, name string) {
	t.Helper()
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	select {
	case <-signal:
	case <-timer.C:
		t.Fatalf("timed out waiting for %s", name)
	}
}

func TestGenPyroscopeConf(t *testing.T) {
	c := Config{
		Name:         "",
		ServerAddr:   "localhost:4040",
		AuthUser:     "user",
		AuthPassword: "password",
		ProfileType: ProfileType{
			Logger:     true,
			CPU:        true,
			Goroutines: true,
			Memory:     true,
			Mutex:      true,
			Block:      true,
		},
	}

	pyroscopeConf := genPyroscopeConf(c)
	assert.Equal(t, c.ServerAddr, pyroscopeConf.ServerAddress)
	assert.Equal(t, c.AuthUser, pyroscopeConf.BasicAuthUser)
	assert.Equal(t, c.AuthPassword, pyroscopeConf.BasicAuthPassword)
	assert.Equal(t, c.Name, pyroscopeConf.ApplicationName)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileCPU)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileGoroutines)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileAllocObjects)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileAllocSpace)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileInuseObjects)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileInuseSpace)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileMutexCount)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileMutexDuration)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileBlockCount)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileBlockDuration)

	setFraction(c)
	resetFraction(c)

	newPyroscopeProfiler(c)
}

func TestNewPyroscopeProfiler(t *testing.T) {
	p := newPyroscopeProfiler(Config{})

	assert.Error(t, p.Start())
	assert.NoError(t, p.Stop())
}

type mockProfiler struct {
	started     syncx.AtomicBool
	stopped     syncx.AtomicBool
	startErr    error
	stopErr     error
	startCalled chan struct{}
	stopCalled  chan struct{}
	startOnce   sync.Once
	stopOnce    sync.Once
}

func newMockProfiler() *mockProfiler {
	return &mockProfiler{startCalled: make(chan struct{}), stopCalled: make(chan struct{})}
}

func (m *mockProfiler) Start() error {
	if m.startErr == nil {
		m.started.Set(true)
	}
	m.startOnce.Do(func() { close(m.startCalled) })
	return m.startErr
}

func (m *mockProfiler) Stop() error {
	if m.stopErr == nil {
		m.stopped.Set(true)
	}
	m.stopOnce.Do(func() { close(m.stopCalled) })
	return m.stopErr
}
