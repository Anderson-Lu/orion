package service

import (
	"context"
	"time"

	"github.com/Anderson-Lu/orion/example/orion_server/proto_go/proto/todo"
)

var c int64 = time.Now().Unix()

func (s *Service) Add(ctx context.Context, in *todo.AddReq) (*todo.AddRsp, error) {
	s.l.Debug("invoke", "in", in)
	// now := time.Now().Unix()
	// gapSec := (now - c)

	// needMockFail := (gapSec > 10 && gapSec < 30) || (gapSec > 50 && gapSec < 70) || (gapSec > 100 && gapSec < 120)
	// if needMockFail {
	// return &todo.AddRsp{}, errors.New("mock error")
	// }
	return &todo.AddRsp{Msg: "ok"}, nil
}
