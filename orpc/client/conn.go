package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"google.golang.org/grpc"
)

type OrionConns interface {
	Invoke(ctx context.Context, connIndex int, method string, req, rsp interface{}, options ...OrionClientInvokeOption) error
	Size() int
}

type RpcConns struct {
	host  string
	size  int
	conns []grpc.ClientConnInterface
}

func (r *RpcConns) buildOptions(opts ...OrionClientInvokeOption) []grpc.CallOption {
	c := []grpc.CallOption{}
	for _, opt := range opts {
		switch opt.Type() {
		case OptionTypeGrpcCallOption:
			for _, v := range opt.Params() {
				c = append(c, v.(grpc.CallOption))
			}
		}
	}
	return c
}

func (r *RpcConns) Invoke(ctx context.Context, connIndex int, method string, req, rsp interface{}, options ...OrionClientInvokeOption) error {
	opts := r.buildOptions(options...)
	return r.conns[connIndex%len(r.conns)].Invoke(ctx, method, req, rsp, opts...)
}

func (r *RpcConns) Size() int {
	return r.size
}

func newRpcConns(host string, size int, opts ...grpc.DialOption) (OrionConns, error) {
	r := &RpcConns{}
	r.host = host
	r.size = size
	for i := 0; i < int(r.size); i++ {
		c, err := grpc.NewClient(r.host, opts...)
		if err != nil {
			return nil, err
		}
		r.conns = append(r.conns, c)
	}
	return r, nil
}

type HttpConns struct {
	conns []*http.Client
	size  int
}

// note that HttpConns only support method 'post' and content-type 'application/json'
func (h *HttpConns) Invoke(ctx context.Context, connIndex int, method string, req, rsp interface{}, options ...OrionClientInvokeOption) error {

	reqBody, _ := json.Marshal(req)

	r, err := http.NewRequest(http.MethodPost, method, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	h.setHeaders(r)

	httpRsp, err := h.conns[connIndex%len(h.conns)].Do(r)
	if err != nil {
		return err
	}

	if httpRsp.StatusCode != http.StatusOK {
		return errors.New("error http code:" + httpRsp.Status)
	}

	defer httpRsp.Body.Close()
	bs, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, &rsp)

	return err
}

func (h *HttpConns) Size() int {
	return len(h.conns)
}

func (h *HttpConns) setHeaders(r *http.Request, opts ...OrionClientInvokeOption) {
	for _, opt := range opts {
		ho, ok := opt.(*HeaderOption)
		if !ok {
			continue
		}
		for k, v := range ho.headers {
			r.Header.Set(k, v)
		}
	}
}

func newHttpConns(size int) (OrionConns, error) {
	r := &HttpConns{}
	r.size = size
	for i := 0; i < int(r.size); i++ {
		c := &http.Client{}
		r.conns = append(r.conns, c)
	}
	return r, nil
}
