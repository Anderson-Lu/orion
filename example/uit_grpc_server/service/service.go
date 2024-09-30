package service

import (
	"github.com/uit/example/uit_grpc_server/proto_go/proto/todo"
	"github.com/uit/pkg/logger"
	"github.com/uit/pkg/uit"
)

func NewService(c *uit.Config) (*Service, error) {

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
	c *uit.Config
	l *logger.Logger
	todo.UnimplementedUitTodoServer
}
