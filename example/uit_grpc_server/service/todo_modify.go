package service

import (
	"context"

	"github.com/uit/example/uit_grpc_server/proto_go/proto/todo"
)

func (s *Service) Modify(ctx context.Context, in *todo.ModifyReq) (*todo.ModifyRsp, error) {
	return &todo.ModifyRsp{}, nil
}
