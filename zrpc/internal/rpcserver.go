package internal

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/lemongoff/hexas/core/proc"
	"github.com/lemongoff/hexas/internal/health"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const probeNamePrefix = "zrpc"

type (
	// ServerOption defines the method to customize a rpcServerOptions.
	ServerOption func(options *rpcServerOptions)

	rpcServerOptions struct {
		health bool
	}

	rpcServer struct {
		*baseRpcServer
		name          string
		healthManager health.Probe
		mu            sync.Mutex
		server        *grpc.Server
		listener      net.Listener
		onReady       func() error
		onDrain       func() error
		onStop        func() error
		lifecycleMu   sync.Mutex
		drainOnce     sync.Once
		stopOnce      sync.Once
		drainErr      error
		stopErr       error
		stopped       bool
		ready         bool
	}
)

// NewRpcServer returns a Server.
func NewRpcServer(addr string, opts ...ServerOption) Server {
	var options rpcServerOptions
	for _, opt := range opts {
		opt(&options)
	}

	return &rpcServer{
		baseRpcServer: newBaseRpcServer(addr, &options),
		healthManager: health.NewHealthManager(fmt.Sprintf("%s-%s", probeNamePrefix, addr)),
	}
}

func (s *rpcServer) SetName(name string) {
	s.name = name
}

func (s *rpcServer) SetLifecycle(onReady func() error, onDrain func() error, onStop func() error) {
	s.mu.Lock()
	s.onReady, s.onDrain, s.onStop = onReady, onDrain, onStop
	s.mu.Unlock()
}

func (s *rpcServer) Start(register RegisterFn) error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	unaryInterceptorOption := grpc.ChainUnaryInterceptor(s.unaryInterceptors...)
	streamInterceptorOption := grpc.ChainStreamInterceptor(s.streamInterceptors...)

	options := append(s.options, unaryInterceptorOption, streamInterceptorOption)
	server := grpc.NewServer(options...)
	if register != nil {
		register(server)
	}

	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		_ = lis.Close()
		return nil
	}
	s.server = server
	s.listener = lis
	onReady := s.onReady
	s.mu.Unlock()

	// Register health before serving, but do not report ready until discovery
	// publication has completed.
	if s.health != nil {
		grpc_health_v1.RegisterHealthServer(server, s.health)
	}
	health.AddProbe(s.healthManager)
	serveDone := make(chan error, 1)
	go func() { serveDone <- server.Serve(lis) }()

	s.lifecycleMu.Lock()
	s.mu.Lock()
	stopped := s.stopped
	s.mu.Unlock()
	if stopped {
		s.lifecycleMu.Unlock()
		return normalizeServeError(<-serveDone)
	}
	if onReady != nil {
		if err := onReady(); err != nil {
			s.stopOnce.Do(func() {
				s.stopErr = err
				s.mu.Lock()
				s.stopped = true
				onStop := s.onStop
				s.mu.Unlock()
				server.Stop()
				_ = lis.Close()
				if onStop != nil {
					s.stopErr = errors.Join(s.stopErr, onStop())
				}
			})
			s.lifecycleMu.Unlock()
			<-serveDone
			return s.stopErr
		}
	}
	if s.health != nil {
		s.health.Resume()
	}
	s.healthManager.MarkReady()
	s.mu.Lock()
	s.ready = true
	s.mu.Unlock()
	s.lifecycleMu.Unlock()

	// we need to make sure all others are wrapped up,
	// so we do graceful stop at shutdown phase instead of wrap up phase
	proc.AddShutdownListener(func() { _ = s.Stop() })

	err = normalizeServeError(<-serveDone)
	if err != nil {
		_ = s.Stop()
	}
	return err
}

func (s *rpcServer) BeginDrain() error {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	return s.beginDrain()
}

func (s *rpcServer) beginDrain() error {
	s.drainOnce.Do(func() {
		if s.health != nil {
			s.health.Shutdown()
		}
		s.mu.Lock()
		ready, onDrain := s.ready, s.onDrain
		s.mu.Unlock()
		if ready && onDrain != nil {
			s.drainErr = onDrain()
		}
	})
	return s.drainErr
}

func (s *rpcServer) Stop() error {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()
	s.stopOnce.Do(func() {
		s.stopErr = s.beginDrain()
		s.mu.Lock()
		s.stopped = true
		server, listener, onStop := s.server, s.listener, s.onStop
		s.mu.Unlock()
		if server != nil {
			server.GracefulStop()
		}
		if listener != nil {
			_ = listener.Close()
		}
		if onStop != nil {
			if err := onStop(); s.stopErr == nil {
				s.stopErr = err
			}
		}
	})
	return s.stopErr
}

func normalizeServeError(err error) error {
	if errors.Is(err, grpc.ErrServerStopped) {
		return nil
	}
	return err
}

// WithRpcHealth returns a func that sets rpc health switch to a Server.
func WithRpcHealth(health bool) ServerOption {
	return func(options *rpcServerOptions) {
		options.health = health
	}
}
