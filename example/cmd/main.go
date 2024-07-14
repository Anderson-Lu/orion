package main

import (
	"log"

	"github.com/uit/example/service"
	"github.com/uit/pkg/logger"
	"github.com/uit/pkg/xgrpc"

	"github.com/uit/example/proto_go/todo"
)

func main() {

	c := &xgrpc.Config{
		GRPC:          &xgrpc.GRPCConfig{Enable: true, Port: 8081},
		HTTP:          &xgrpc.HTTPConfig{Enable: true, Port: 8080},
		FrameLogger:   &logger.LoggerConfig{Dir: []string{"..", "log", "frame.log"}},
		AccessLogger:  &logger.LoggerConfig{Dir: []string{"..", "log", "access.log"}},
		ServiceLogger: &logger.LoggerConfig{Dir: []string{"..", "log", "service.log"}},
	}

	server, err := xgrpc.New(c)
	if err != nil {
		log.Fatal(err)
	}

	srv := service.NewService(c, server.SvcLogger())

	todo.RegisterUitTodoServer(server.GRPCServer(), srv)

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
