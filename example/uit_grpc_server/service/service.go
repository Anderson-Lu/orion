package service

import (
	"github.com/orion/example/uit_grpc_server/proto_go/proto/todo"
	"github.com/orion/orpc"
	"github.com/orion/pkg/logger"
)

func NewService(c *orpc.Config) (*Service, error) {

	lg, err := logger.NewLogger(c.ServiceLogger)
	if err != nil {
		return nil, err
	}

	return &Service{
		c: c,
		l: lg,
	}, nil
}

type Service struct {
	c *orpc.Config
	l *logger.Logger
	todo.UnimplementedUitTodoServer
}
