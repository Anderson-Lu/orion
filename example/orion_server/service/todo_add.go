package service

import (
	"context"

<<<<<<< HEAD:example/orion_server/service/todo_add.go
	"github.com/orion/example/orion_server/proto_go/proto/todo"
=======
	"github.com/Anderson-Lu/orion/example/orion_server/proto_go/proto/todo"
>>>>>>> dev_0_0_2:example/orion_server/service/todo_add.go
)

func (s *Service) Add(ctx context.Context, in *todo.AddReq) (*todo.AddRsp, error) {
	s.l.Debug("invoke", "in", in)
	return &todo.AddRsp{}, nil
}
