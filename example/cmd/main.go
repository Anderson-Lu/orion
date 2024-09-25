package main

import (
	"log"

	"github.com/uit/pkg/logger"
	"github.com/uit/pkg/xgrpc"
	_ "github.com/uit/pkg/xgrpc/build"
	"github.com/uit/pkg/xgrpc/options"

	"github.com/uit/example/proto_go/todo"
	"github.com/uit/example/service"
)

func main() {

	c := &xgrpc.Config{
		GRPC:          &xgrpc.GRPCConfig{Enable: true, Port: 8081},
		HTTP:          &xgrpc.HTTPConfig{Enable: true, Port: 8080},
		FrameLogger:   &logger.LoggerConfig{Path: []string{"..", "log", "frame.log"}, LogLevel: "info"},
		AccessLogger:  &logger.LoggerConfig{Path: []string{"..", "log", "access.log"}},
		ServiceLogger: &logger.LoggerConfig{Path: []string{"..", "log", "service.log"}},
		PanicLogger:   &logger.LoggerConfig{Path: []string{"..", "log", "panic.log"}},
	}

	handler, _ := service.NewService(c)
	server, err := xgrpc.New(c,
		options.WithHandler(handler, &todo.UitTodo_ServiceDesc),
		options.WithFlags(),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
