package service

import (
	"github.com/uit/example/proto_go/todo"
	"github.com/uit/pkg/logger"
	"github.com/uit/pkg/xgrpc"
)

func NewService(c *xgrpc.Config, lgsvc *logger.Logger) *Service {
	return &Service{
		c: c,
		l: lgsvc,
	}
}

type Service struct {
	c *xgrpc.Config
	l *logger.Logger
	todo.UnimplementedUitTodoServer
}
