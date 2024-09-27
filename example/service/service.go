package service

import (
	"github.com/uit/example/proto_go/proto/todo"
	"github.com/uit/pkg/logger"
	"github.com/uit/pkg/xgrpc"
)

func NewService(c *xgrpc.Config) (*Service, error) {

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
	c *xgrpc.Config
	l *logger.Logger
	todo.UnimplementedUitTodoServer
}
