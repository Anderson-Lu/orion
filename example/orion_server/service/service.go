package service

import (
<<<<<<< HEAD
	"github.com/orion/example/orion_server/config"
	"github.com/orion/example/orion_server/proto_go/proto/todo"
	"github.com/orion/pkg/logger"
=======
	"github.com/Anderson-Lu/orion/example/orion_server/config"
	"github.com/Anderson-Lu/orion/example/orion_server/proto_go/proto/todo"
	"github.com/Anderson-Lu/orion/pkg/logger"
>>>>>>> dev_0_0_2
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
	todo.UnimplementedUitTodoServer
}
