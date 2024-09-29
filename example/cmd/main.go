package main

import (
	"log"

	"github.com/uit/pkg/logger"
	"github.com/uit/pkg/uit"
	_ "github.com/uit/pkg/uit/build"

	"github.com/uit/example/proto_go/proto/todo"
	"github.com/uit/example/service"
)

func main() {

	c := &uit.Config{
		Server:          &uit.ServerConfig{Port: 8080, EnableGRPCGateway: true},
		PromtheusConfig: &uit.PromtheusConfig{Enable: true, Port: 9092},
		FrameLogger:     &logger.LoggerConfig{Path: []string{"..", "log", "frame.log"}, LogLevel: "info"},
		AccessLogger:    &logger.LoggerConfig{Path: []string{"..", "log", "access.log"}},
		ServiceLogger:   &logger.LoggerConfig{Path: []string{"..", "log", "service.log"}},
		PanicLogger:     &logger.LoggerConfig{Path: []string{"..", "log", "panic.log"}},
	}

	handler, _ := service.NewService(c)
	server, err := uit.New(
		uit.WithConfigFile("../config/config.toml"),
		uit.WithGRPCHandler(handler, &todo.UitTodo_ServiceDesc),
		uit.WithGrpcGatewayEndpointFunc(todo.RegisterUitTodoHandlerFromEndpoint),
		uit.WithFlags(),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
