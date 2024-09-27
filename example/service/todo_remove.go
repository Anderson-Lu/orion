package service

import (
	"context"

	"github.com/uit/example/proto_go/proto/todo"
)

func (s *Service) Remove(ctx context.Context, in *todo.RemoveReq) (*todo.RemoveRsp, error) {
	return &todo.RemoveRsp{}, nil
}
