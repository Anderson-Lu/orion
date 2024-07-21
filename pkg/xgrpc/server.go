package xgrpc

import (
	"fmt"

	"net"
	"net/http"

	"github.com/uit/pkg/logger"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/uit/pkg/event"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "runtime/pprof"

	_ "go.uber.org/automaxprocs"
)

type Server struct {
	g *grpc.Server
	m *runtime.ServeMux
	c *Config

	e    *event.EventHub
	lfra *logger.Logger
	lacc *logger.Logger
	lsvc *logger.Logger
}

func New(c *Config) (*Server, error) {

	s := &Server{c: c}
	if err := s.initLogger(); err != nil {
		return nil, err
	}

	ev := event.NewEventHub(
		100,
		event.WithMultiDiapatcher(10),
		event.WithConsumer(uint32(EventTypeLog), s.evLogger),
	)
	s.e = ev
	s.g = grpc.NewServer(grpc.UnaryInterceptor(AccessLoggerInterceptor(s.lacc)))

	reflection.Register(s.g)
	return s, nil
}

func (s *Server) SvcLogger() *logger.Logger {
	return s.lsvc
}

func (s *Server) initLogger() error {

	lg, err := logger.NewLogger(s.c.FrameLogger)
	if err != nil {
		return err
	}
	s.lfra = lg

	sl, err := logger.NewLogger(s.c.AccessLogger)
	if err != nil {
		return err
	}
	s.lacc = sl

	sv, err := logger.NewLogger(s.c.ServiceLogger)
	if err != nil {
		return err
	}
	s.lsvc = sv

	return nil
}

func (s *Server) GRPCServer() *grpc.Server {
	return s.g
}

func (s *Server) MuxServer() *runtime.ServeMux {
	return s.m
}

func (s *Server) Serve() error {

	defer s.lfra.Sync()
	if s.lacc != nil {
		defer s.lacc.Sync()
	}

	eg := errgroup.Group{}

	eg.Go(func() error {
		if s.c.GRPC != nil && s.c.GRPC.Enable {
			return s.serveGRPCServer()
		}
		return nil
	})

	eg.Go(func() error {
		if s.c.HTTP != nil && s.c.HTTP.Enable {
			return s.serveHTTPServer()
		}
		return nil
	})

	return eg.Wait()
}

func (s *Server) Stop() {

}

func (s *Server) serveHTTPServer() error {
	s.m = runtime.NewServeMux()
	fmt.Printf("[Server] http prepare on: %d \n", s.c.HTTP.Port)
	s.e.Publish(&Event{typ: EventTypeLog, data: fmt.Sprintf("[Server] http server stared on: %d", s.c.HTTP.Port)})
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.c.HTTP.Port), s.m); err != nil {
		fmt.Printf("[Server] http started fail, port: %d, err:%+v \n", s.c.HTTP.Port, err.Error())
		return err
	}
	return nil
}

func (s *Server) serveGRPCServer() error {

	fmt.Printf("[Server] gRPC prepare on: %d \n", s.c.GRPC.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.c.GRPC.Port))
	if err != nil {
		fmt.Printf("[Server] gRPC started fail, port: %d, err:%+v \n", s.c.GRPC.Port, err.Error())
		return err
	}
	s.e.Publish(&Event{typ: EventTypeLog, data: fmt.Sprintf("[Server] gRPC server stared on: %d", s.c.GRPC.Port)})
	err = s.g.Serve(lis)
	if err != nil {
		fmt.Printf("[Server] gRPC started fail, port: %d, err:%+v \n", s.c.GRPC.Port, err.Error())
		return err
	}
	return nil
}

func (s *Server) evLogger(msg event.Event) error {
	s.lfra.Info("[Server]", "EvType", msg.Type(), "Ev", msg.Data())
	return nil
}
