package service

import (
	"context"

	"github.com/orion/example/orion_server/proto_go/proto/todo"
)

func (s *Service) Add(ctx context.Context, in *todo.AddReq) (*todo.AddRsp, error) {
	s.l.Debug("invoke", "in", in)
	return &todo.AddRsp{}, nil
}
