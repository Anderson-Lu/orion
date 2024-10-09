package service

import (
	"context"

	"github.com/orion/example/uit_grpc_server/proto_go/proto/todo"
)

func (s *Service) List(ctx context.Context, in *todo.ListReq) (*todo.ListRsp, error) {
	return &todo.ListRsp{}, nil
}
