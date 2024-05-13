// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.26.1
// source: v1/log.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationLogProduce = "/api.v1.Log/Produce"

type LogHTTPServer interface {
	Produce(context.Context, *ProduceRequest) (*ProduceResponse, error)
}

func RegisterLogHTTPServer(s *http.Server, srv LogHTTPServer) {
	r := s.Route("/")
	r.POST("/log", _Log_Produce0_HTTP_Handler(srv))
}

func _Log_Produce0_HTTP_Handler(srv LogHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ProduceRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLogProduce)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Produce(ctx, req.(*ProduceRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ProduceResponse)
		return ctx.Result(200, reply)
	}
}

type LogHTTPClient interface {
	Produce(ctx context.Context, req *ProduceRequest, opts ...http.CallOption) (rsp *ProduceResponse, err error)
}

type LogHTTPClientImpl struct {
	cc *http.Client
}

func NewLogHTTPClient(client *http.Client) LogHTTPClient {
	return &LogHTTPClientImpl{client}
}

func (c *LogHTTPClientImpl) Produce(ctx context.Context, in *ProduceRequest, opts ...http.CallOption) (*ProduceResponse, error) {
	var out ProduceResponse
	pattern := "/log"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationLogProduce))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}