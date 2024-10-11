package service

import (
	"context"

	"github.com/Anderson-Lu/orion/example/orion_server/proto_go/proto/todo"
)

func (s *Service) List(ctx context.Context, in *todo.ListReq) (*todo.ListRsp, error) {
	return &todo.ListRsp{}, nil
}
