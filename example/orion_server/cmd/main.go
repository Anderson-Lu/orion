package main

import (
	"log"

<<<<<<< HEAD
	"github.com/orion/orpc"
	_ "github.com/orion/orpc/build"

	"github.com/orion/example/orion_server/config"
	"github.com/orion/example/orion_server/proto_go/proto/todo"
	"github.com/orion/example/orion_server/service"
=======
	"github.com/Anderson-Lu/orion/orpc"
	_ "github.com/Anderson-Lu/orion/orpc/build"

	"github.com/Anderson-Lu/orion/example/orion_server/config"
	"github.com/Anderson-Lu/orion/example/orion_server/proto_go/proto/todo"
	"github.com/Anderson-Lu/orion/example/orion_server/service"
>>>>>>> dev_0_0_2
)

func main() {

	c := &config.Config{}

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
