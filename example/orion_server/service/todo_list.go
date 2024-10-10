package service

import (
	"context"

<<<<<<< HEAD:example/orion_server/service/todo_list.go
	"github.com/orion/example/orion_server/proto_go/proto/todo"
=======
	"github.com/Anderson-Lu/orion/example/orion_server/proto_go/proto/todo"
>>>>>>> dev_0_0_2:example/orion_server/service/todo_list.go
)

func (s *Service) List(ctx context.Context, in *todo.ListReq) (*todo.ListRsp, error) {
	return &todo.ListRsp{}, nil
}
