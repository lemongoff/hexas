package internal

import (
	"sync"
	"testing"
	"time"

	"github.com/lemongoff/hexas/core/proc"
	"github.com/lemongoff/hexas/internal/mock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestRpcServer(t *testing.T) {
	server := NewRpcServer("localhost:54321", WithRpcHealth(true))
	server.SetName("mock")
	var wg, wgDone sync.WaitGroup
	var grpcServer *grpc.Server
	var lock sync.Mutex
	wg.Add(1)
	wgDone.Add(1)
	go func() {
		err := server.Start(func(server *grpc.Server) {
			lock.Lock()
			mock.RegisterDepositServiceServer(server, new(mock.DepositServer))
			grpcServer = server
			lock.Unlock()
			wg.Done()
		})
		assert.Nil(t, err)
		wgDone.Done()
	}()

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	lock.Lock()
	grpcServer.GracefulStop()
	lock.Unlock()

	proc.Shutdown()
	wgDone.Wait()
}

func TestRpcServer_WithBadAddress(t *testing.T) {
	server := NewRpcServer("localhost:111111", WithRpcHealth(true))
	server.SetName("mock")
	err := server.Start(func(server *grpc.Server) {
		mock.RegisterDepositServiceServer(server, new(mock.DepositServer))
	})
	assert.NotNil(t, err)

	proc.WrapUp()
}

func TestRpcServerLifecycle(t *testing.T) {
	server := NewRpcServer("localhost:0")
	ready := make(chan struct{})
	drained := make(chan struct{})
	stopped := make(chan struct{})
	server.SetLifecycle(
		func() error { close(ready); return nil },
		func() error { close(drained); return nil },
		func() error { close(stopped); return nil },
	)
	done := make(chan error, 1)
	go func() { done <- server.Start(nil) }()
	<-ready
	assert.NoError(t, server.BeginDrain())
	<-drained
	assert.NoError(t, server.Stop())
	<-stopped
	assert.NoError(t, <-done)
}

func TestRpcServerSerializesReadyAndStop(t *testing.T) {
	server := NewRpcServer("localhost:0")
	readyEntered := make(chan struct{})
	releaseReady := make(chan struct{})
	drained := make(chan struct{})
	stopped := make(chan struct{})
	server.SetLifecycle(
		func() error {
			close(readyEntered)
			<-releaseReady
			return nil
		},
		func() error { close(drained); return nil },
		func() error { close(stopped); return nil },
	)
	startDone := make(chan error, 1)
	go func() { startDone <- server.Start(nil) }()
	<-readyEntered
	stopDone := make(chan error, 1)
	go func() { stopDone <- server.Stop() }()
	select {
	case <-stopped:
		t.Fatal("server stopped before ready publication completed")
	case <-time.After(20 * time.Millisecond):
	}
	close(releaseReady)
	assert.NoError(t, <-stopDone)
	<-drained
	<-stopped
	assert.NoError(t, <-startDone)
}
