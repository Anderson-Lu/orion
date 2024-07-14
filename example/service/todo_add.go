package service

import (
	"context"

	"github.com/uit/example/proto_go/todo"
)

func (s *Service) Add(ctx context.Context, in *todo.AddReq) (*todo.AddRsp, error) {
	s.l.Debug("invoke", "in", in)
	return &todo.AddRsp{}, nil
}
