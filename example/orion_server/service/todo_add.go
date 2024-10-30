package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/Anderson-Lu/orion/example/orion_server/proto_go/proto/todo"
	"go.opentelemetry.io/otel/trace"
)

func (s *Service) Add(ctx context.Context, in *todo.AddReq) (*todo.AddRsp, error) {

	var span trace.Span

	if s.t != nil {
		_, span = s.t.SpanClient(ctx, "服务端处理逻辑")
	}
	// time.Sleep(time.Second * 10)
	// now := time.Now().Unix()
	// gapSec := (now - c)
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	time.Sleep(time.Millisecond * time.Duration(r.Intn(500)))

	// needMockFail := (gapSec > 10 && gapSec < 30) || (gapSec > 50 && gapSec < 70) || (gapSec > 100 && gapSec < 120)
	// if needMockFail {
	// return &todo.AddRsp{}, errors.New("mock error")
	// }

	span.End()
	return &todo.AddRsp{Msg: "ok"}, nil
}
