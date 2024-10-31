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
	var span2 trace.Span

	if s.t != nil {
		ctx, span = s.t.ServerSpan(ctx, "server handle")
	}

	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	time.Sleep(time.Millisecond * time.Duration(r.Intn(500)))

	_, span2 = s.t.ServerSpan(ctx, "query db")
	time.Sleep(time.Millisecond * 10)
	span2.End()

	span.End()
	return &todo.AddRsp{Msg: "ok"}, nil
}
