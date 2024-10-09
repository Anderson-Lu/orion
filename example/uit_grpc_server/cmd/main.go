package main

import (
	"log"

	"github.com/orion/orpc"
	_ "github.com/orion/orpc/build"
	"github.com/orion/pkg/logger"

	"github.com/orion/example/uit_grpc_server/proto_go/proto/todo"
	"github.com/orion/example/uit_grpc_server/service"
)

func main() {

	c := &orpc.Config{
		Server:          &orpc.ServerConfig{Port: 8080, EnableGRPCGateway: true},
		PromtheusConfig: &orpc.PromtheusConfig{Enable: true, Port: 9092},
		FrameLogger:     &logger.LoggerConfig{Path: "../log/frame.log", LogLevel: "info"},
		AccessLogger:    &logger.LoggerConfig{Path: "../log/access.log"},
		ServiceLogger:   &logger.LoggerConfig{Path: "../log/service.log"},
		PanicLogger:     &logger.LoggerConfig{Path: "../log/panic.log"},
	}

	handler, _ := service.NewService(c)
	server, err := orpc.New(
		orpc.WithConfigFile("../config/config.toml"),
		orpc.WithGRPCHandler(handler, &todo.UitTodo_ServiceDesc),
		orpc.WithGrpcGatewayEndpointFunc(todo.RegisterUitTodoHandlerFromEndpoint),
		orpc.WithFlags(),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
