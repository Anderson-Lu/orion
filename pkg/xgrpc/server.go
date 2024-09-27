package xgrpc

import (
	"fmt"
	"os"

	"net"
	"net/http"

	"github.com/uit/pkg/logger"
	"github.com/uit/pkg/xgrpc/build"
	"github.com/uit/pkg/xgrpc/interceptors"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "runtime/pprof"

	_ "go.uber.org/automaxprocs"
)

var (
	defaultFrameLogger  = &logger.LoggerConfig{Path: []string{"..", "log", "frame.log"}, LogLevel: "info"}
	defaultAccessLogger = &logger.LoggerConfig{Path: []string{"..", "log", "access.log"}, LogLevel: "info"}
	defaultPanicLogger  = &logger.LoggerConfig{Path: []string{"..", "log", "panic.log"}, LogLevel: "error"}
)

type Server struct {
	g *grpc.Server
	m *runtime.ServeMux
	p *http.ServeMux
	c *Config

	panicLogger *logger.Logger
	accLogger   *logger.Logger
	frameLogger *logger.Logger

	cmdMode bool
}

func New(c *Config, opts ...ServerOption) (*Server, error) {
	s := &Server{c: c}
	if err := s.initLogger(); err != nil {
		return nil, err
	}
	s.initGrpcServer()
	s.initOptions()
	return s, nil
}

func (s *Server) initOptions(opts ...ServerOption) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *Server) initGrpcServer() {
	s.g = grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.ChainInterceptors(
			interceptors.ContextWrapperInterceptor(s.frameLogger),
			interceptors.AccessInterceptor(s.accLogger),
			interceptors.PanicInterceptor(s.panicLogger),
		)),
	)
	reflection.Register(s.g)
}

func (s *Server) initLogger() error {

	if s.c.FrameLogger == nil {
		s.c.FrameLogger = defaultFrameLogger
	}
	if s.c.AccessLogger == nil {
		s.c.AccessLogger = defaultAccessLogger
	}
	if s.c.PanicLogger == nil {
		s.c.AccessLogger = defaultPanicLogger
	}

	lg, err := logger.NewLogger(s.c.FrameLogger)
	if err != nil {
		return err
	}
	s.frameLogger = lg

	sl, err := logger.NewLogger(s.c.AccessLogger)
	if err != nil {
		return err
	}
	s.accLogger = sl

	pl, err := logger.NewLogger(s.c.PanicLogger)
	if err != nil {
		return err
	}
	s.panicLogger = pl

	return nil
}

func (s *Server) serveFlags() bool {
	for _, args := range os.Args {
		switch args {
		case "-v", "--verbose":
			build.PrintVerbose()
			return true
		}
	}
	return false
}

func (s *Server) ListenAndServe() error {

	if s.cmdMode {
		if ok := s.serveFlags(); ok {
			return nil
		}
	}

	defer s.frameLogger.Sync()
	defer s.frameLogger.Sync()
	defer s.panicLogger.Sync()

	eg := errgroup.Group{}
	eg.Go(func() error {
		if s.c.GRPC != nil && s.c.GRPC.Enable {
			return s.runGRPCServer()
		}
		return nil
	})

	eg.Go(func() error {
		if s.c.HTTP != nil && s.c.HTTP.Enable {
			return s.runGRPCGateway()
		}
		return nil
	})

	eg.Go(func() error {
		if s.c.PromtheusConfig == nil || !s.c.PromtheusConfig.Enable || s.c.PromtheusConfig.Port == 0 {
			return nil
		}
		return s.runPromtheusMetrics()
	})
	return eg.Wait()
}

func (s *Server) Stop() {
	s.frameLogger.Info("[Server] server stopped", "port", s.c.GRPC.Port)
}

func (s *Server) runPromtheusMetrics() error {
	s.p = http.NewServeMux()
	s.p.Handle("/metrics", promhttp.Handler())
	s.frameLogger.Info("[Server] promtheus metrics server started succ", "port", s.c.PromtheusConfig.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.c.PromtheusConfig.Port), s.p); err != nil {
		s.frameLogger.Info("[Server] promtheus metrics server started fail", "port", s.c.PromtheusConfig.Port, "err", err.Error())
		return err
	}
	return nil
}

func (s *Server) runGRPCGateway() error {
	s.m = runtime.NewServeMux()
	s.frameLogger.Info("[Server] http server started succ", "port", s.c.HTTP.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.c.HTTP.Port), s.m); err != nil {
		s.frameLogger.Info("[Server] http server started fail", "port", s.c.HTTP.Port, "err", err.Error())
		return err
	}
	return nil
}

func (s *Server) runGRPCServer() error {
	s.frameLogger.Info("[Server] gRPC server started succ", "port", s.c.GRPC.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.c.GRPC.Port))
	if err != nil {
		s.frameLogger.Info("[Server] gRPC server started fail", "port", s.c.GRPC.Port, "err", err.Error())
		return err
	}
	err = s.g.Serve(lis)
	if err != nil {
		s.frameLogger.Info("[Server] gRPC server started fail", "port", s.c.GRPC.Port, "err", err.Error())
		return err
	}
	return nil
}
