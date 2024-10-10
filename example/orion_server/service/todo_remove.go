package service

import (
	"context"

<<<<<<< HEAD:example/orion_server/service/todo_remove.go
	"github.com/orion/example/orion_server/proto_go/proto/todo"
=======
	"github.com/Anderson-Lu/orion/example/orion_server/proto_go/proto/todo"
>>>>>>> dev_0_0_2:example/orion_server/service/todo_remove.go
)

func (s *Service) Remove(ctx context.Context, in *todo.RemoveReq) (*todo.RemoveRsp, error) {
	return &todo.RemoveRsp{}, nil
}
