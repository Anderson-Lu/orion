package service

import (
	"context"

	"github.com/orion/example/orion_server/proto_go/proto/todo"
)

func (s *Service) Remove(ctx context.Context, in *todo.RemoveReq) (*todo.RemoveRsp, error) {
	return &todo.RemoveRsp{}, nil
}