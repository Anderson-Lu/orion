package main

import (
	"log"

	"github.com/uit/modules/logger"
	"github.com/uit/urpc"
	_ "github.com/uit/urpc/build"

	"github.com/uit/example/uit_grpc_server/proto_go/proto/todo"
	"github.com/uit/example/uit_grpc_server/service"
)

func main() {

	c := &urpc.Config{
		Server:          &urpc.ServerConfig{Port: 8080, EnableGRPCGateway: true},
		PromtheusConfig: &urpc.PromtheusConfig{Enable: true, Port: 9092},
		FrameLogger:     &logger.LoggerConfig{Path: "../log/frame.log", LogLevel: "info"},
		AccessLogger:    &logger.LoggerConfig{Path: "../log/access.log"},
		ServiceLogger:   &logger.LoggerConfig{Path: "../log/service.log"},
		PanicLogger:     &logger.LoggerConfig{Path: "../log/panic.log"},
	}

	handler, _ := service.NewService(c)
	server, err := urpc.New(
		urpc.WithConfigFile("../config/config.toml"),
		urpc.WithGRPCHandler(handler, &todo.UitTodo_ServiceDesc),
		urpc.WithGrpcGatewayEndpointFunc(todo.RegisterUitTodoHandlerFromEndpoint),
		urpc.WithFlags(),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
