package orpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"net"
	"net/http"

	"github.com/Anderson-Lu/orion/orpc/build"
	"github.com/Anderson-Lu/orion/orpc/interceptors"
	"github.com/Anderson-Lu/orion/orpc/registry"
	"github.com/Anderson-Lu/orion/orpc/registry/consul"
	"github.com/Anderson-Lu/orion/orpc/tracing"

	"github.com/Anderson-Lu/orion/pkg/logger"
	"github.com/Anderson-Lu/orion/pkg/utils"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	_ "runtime/pprof"

	_ "github.com/Anderson-Lu/orion/orpc/codec"

	_ "go.uber.org/automaxprocs"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	defaultFrameLogger  = &logger.LoggerConfig{Path: "../log/frame.log", LogLevel: "info"}
	defaultAccessLogger = &logger.LoggerConfig{Path: "../log/access.log", LogLevel: "info"}
	defaultPanicLogger  = &logger.LoggerConfig{Path: "../log/panic.log", LogLevel: "error"}
	defaultGRPCOptions  = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
)

type Server struct {
	gServer    *grpc.Server
	gatewayMux *runtime.ServeMux
	httpMux    *http.ServeMux
	promMux    *http.ServeMux
	c          *Config

	panicLogger *logger.Logger
	accLogger   *logger.Logger
	frameLogger *logger.Logger

	cmdMode      bool
	grpcOpts     []grpc.DialOption
	gatewayFunc  func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	grpcHandlers map[interface{}][]*grpc.ServiceDesc

	rsy  registry.IRegistry
	sig  chan os.Signal
	gsig chan any

	trace *tracing.Tracing
}

func New(opts ...ServerOption) (*Server, error) {

	s := &Server{}
	s.sig = make(chan os.Signal, 1)
	s.gsig = make(chan any, 1)

	if err := s.initOptions(opts...); err != nil {
		return nil, err
	}

	if err := s.initLogger(); err != nil {
		return nil, err
	}

	if err := s.initTracing(); err != nil {
		return nil, err
	}

	s.grpcOpts = defaultGRPCOptions
	return s, nil
}

func (s *Server) Tracing() *tracing.Tracing {
	return s.trace
}

func (s *Server) initTracing() error {
	if s.c.Tracing == nil {
		return nil
	}

	bs := &tracing.Resources{}
	bs.Env(s.c.Tracing.Env)
	bs.Namespace(s.c.Tracing.Namespace)
	bs.IP(utils.IP().GetLocalIP())
	bs.ServiceName(s.c.Tracing.ServiceName)
	bs.InstanceId(s.c.Tracing.InstanceId)

	tr, err := tracing.NewTracing(s.c.Tracing.ServiceName, tracing.WithOpenTelemetryAddress(s.c.Tracing.Address), tracing.WithResource(bs))
	if err != nil {
		return err
	}
	s.trace = tr
	s.trace.Start()

	return nil
}

func (s *Server) initOptions(opts ...ServerOption) error {
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) initGrpcServer() error {
	s.gServer = grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.ChainInterceptors(
			interceptors.ContextWrapperInterceptor(s.frameLogger),
			interceptors.AccessInterceptor(s.accLogger),
			interceptors.PanicInterceptor(s.panicLogger),
			interceptors.RateLimitorInterceptor(s.c.RateLimit, s.frameLogger),
		)),
	)
	reflection.Register(s.gServer)
	for handler, sds := range s.grpcHandlers {
		for _, sd := range sds {
			s.gServer.RegisterService(sd, handler)
		}
	}

	switch rsyi := s.rsy.(type) {
	case *consul.OrionConsulRegistry:
		ip := utils.IP().GetLocalIP()
		port := s.c.Server.Port
		if err := rsyi.AddNode(context.Background(), ip, port); err != nil {
			return err
		}
		rsyi.RegisterHealthHandler(s.gServer)
	}

	return nil
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
		case "-v", "--version":
			build.PrintVerbose()
			return true
		}
	}
	return false
}

func (s *Server) ListenAndServe() error {

	signal.Notify(s.sig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT)

	if s.cmdMode {
		if ok := s.serveFlags(); ok {
			return nil
		}
	}

	defer s.frameLogger.Sync()
	defer s.frameLogger.Sync()
	defer s.panicLogger.Sync()

	go s.startServers()

	select {
	case stopSig := <-s.sig:
		s.stop(fmt.Sprintf("accept system signal stopped, signal:[%s]", stopSig.String()))
	case stopInfo := <-s.gsig:
		s.stop(stopInfo)
	}
	return nil
}

func (s *Server) startServers() {
	fn := func() error {
		eg := errgroup.Group{}
		eg.Go(func() error {
			if s.c.Server != nil {
				return s.start()
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
	if err := fn(); err != nil {
		s.gsig <- err
	}
}

func (s *Server) stop(info any) {
	s.frameLogger.Info("[Server] server stopped", "port", s.c.Server.Port, "context", info)

	if s.rsy != nil {
		s.rsy.RemoveNode(context.Background())
	}
	if s.trace != nil {
		s.trace.Shutdown(context.Background())
	}
}

func (s *Server) runPromtheusMetrics() error {
	s.promMux = http.NewServeMux()
	s.promMux.Handle("/metrics", promhttp.Handler())
	s.frameLogger.Info("[Server] promtheus metrics server started succ", "port", s.c.PromtheusConfig.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.c.PromtheusConfig.Port), s.promMux); err != nil {
		s.frameLogger.Info("[Server] promtheus metrics server started fail", "port", s.c.PromtheusConfig.Port, "err", err.Error())
		return err
	}
	return nil
}

func (s *Server) start() error {

	if err := s.initGrpcServer(); err != nil {
		s.frameLogger.Info("[Server] gRPC server init fail", "port", s.c.Server.Port, "err", err.Error())
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.c.Server.Port))
	if err != nil {
		s.frameLogger.Info("[Server] gRPC server started fail", "port", s.c.Server.Port, "err", err.Error())
		log.Fatal(err)
	}

	// only grpc server
	if !s.c.Server.EnableGRPCGateway {
		s.frameLogger.Info("[Server] gRPC server started succ", "port", s.c.Server.Port)
		err = s.gServer.Serve(lis)
		if err != nil {
			s.frameLogger.Info("[Server] gRPC server started fail", "port", s.c.Server.Port, "err", err.Error())
			log.Fatal(err)
		}
		return nil
	}

	s.gatewayMux = runtime.NewServeMux()

	if s.gatewayFunc == nil {
		return nil
	}
	if err := s.gatewayFunc(context.Background(), s.gatewayMux, fmt.Sprintf(":%d", s.c.Server.Port), defaultGRPCOptions); err != nil {
		s.frameLogger.Info("[Server] gRPC server started fail", "port", s.c.Server.Port, "err", err.Error())
		log.Fatal(err)
	}

	s.httpMux = http.NewServeMux()
	s.httpMux.Handle("/", s.gatewayMux)

	// both grpc server and http server in one port
	integrateServer := h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			s.gServer.ServeHTTP(w, r)
		} else {
			s.httpMux.ServeHTTP(w, r)
		}
	}), &http2.Server{})

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.c.Server.Port),
		Handler: integrateServer,
	}
	if err := gwServer.Serve(lis); err != nil {
		s.frameLogger.Info("[Server] gRPC gateway server started fail", "port", s.c.Server.Port, "err", err.Error())
		log.Fatal(err)
	}
	return nil
}
