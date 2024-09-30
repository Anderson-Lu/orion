package service

import (
	"github.com/uit/example/uit_grpc_server/proto_go/proto/todo"
	"github.com/uit/modules/logger"
	"github.com/uit/urpc"
)

func NewService(c *urpc.Config) (*Service, error) {

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
	c *urpc.Config
	l *logger.Logger
	todo.UnimplementedUitTodoServer
}
