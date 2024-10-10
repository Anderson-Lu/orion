package service

import (
	"context"

<<<<<<< HEAD:example/orion_server/service/todo_modify.go
	"github.com/orion/example/orion_server/proto_go/proto/todo"
=======
	"github.com/Anderson-Lu/orion/example/orion_server/proto_go/proto/todo"
>>>>>>> dev_0_0_2:example/orion_server/service/todo_modify.go
)

func (s *Service) Modify(ctx context.Context, in *todo.ModifyReq) (*todo.ModifyRsp, error) {
	return &todo.ModifyRsp{}, nil
}
