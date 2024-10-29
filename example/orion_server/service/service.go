package service

import (
	"github.com/Anderson-Lu/orion/example/orion_server/config"
	"github.com/Anderson-Lu/orion/example/orion_server/proto_go/proto/todo"
	"github.com/Anderson-Lu/orion/orpc/tracing"
	"github.com/Anderson-Lu/orion/pkg/logger"
)

func NewService(c *config.Config) (*Service, error) {

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
	c *config.Config
	l *logger.Logger
	t *tracing.Tracing
	todo.UnimplementedUitTodoServer
}

func (s *Service) SetTracing(t *tracing.Tracing) {
	s.t = t
}
