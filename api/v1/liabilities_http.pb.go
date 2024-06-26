// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.26.1
// source: v1/liabilities.proto

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

const OperationLiabilitiesGetLiabilities = "/api.v1.Liabilities/GetLiabilities"
const OperationLiabilitiesGetLiability = "/api.v1.Liabilities/GetLiability"

type LiabilitiesHTTPServer interface {
	GetLiabilities(context.Context, *GetLiabilitiesRequest) (*GetLiabilitiesReply, error)
	GetLiability(context.Context, *GetLiabilityRequest) (*GetLiabilityReply, error)
}

func RegisterLiabilitiesHTTPServer(s *http.Server, srv LiabilitiesHTTPServer) {
	r := s.Route("/")
	r.GET("/liabilities/{id}", _Liabilities_GetLiability0_HTTP_Handler(srv))
	r.GET("/liabilities", _Liabilities_GetLiabilities0_HTTP_Handler(srv))
}

func _Liabilities_GetLiability0_HTTP_Handler(srv LiabilitiesHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetLiabilityRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLiabilitiesGetLiability)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetLiability(ctx, req.(*GetLiabilityRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetLiabilityReply)
		return ctx.Result(200, reply)
	}
}

func _Liabilities_GetLiabilities0_HTTP_Handler(srv LiabilitiesHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetLiabilitiesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLiabilitiesGetLiabilities)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetLiabilities(ctx, req.(*GetLiabilitiesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetLiabilitiesReply)
		return ctx.Result(200, reply)
	}
}

type LiabilitiesHTTPClient interface {
	GetLiabilities(ctx context.Context, req *GetLiabilitiesRequest, opts ...http.CallOption) (rsp *GetLiabilitiesReply, err error)
	GetLiability(ctx context.Context, req *GetLiabilityRequest, opts ...http.CallOption) (rsp *GetLiabilityReply, err error)
}

type LiabilitiesHTTPClientImpl struct {
	cc *http.Client
}

func NewLiabilitiesHTTPClient(client *http.Client) LiabilitiesHTTPClient {
	return &LiabilitiesHTTPClientImpl{client}
}

func (c *LiabilitiesHTTPClientImpl) GetLiabilities(ctx context.Context, in *GetLiabilitiesRequest, opts ...http.CallOption) (*GetLiabilitiesReply, error) {
	var out GetLiabilitiesReply
	pattern := "/liabilities"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationLiabilitiesGetLiabilities))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LiabilitiesHTTPClientImpl) GetLiability(ctx context.Context, in *GetLiabilityRequest, opts ...http.CallOption) (*GetLiabilityReply, error) {
	var out GetLiabilityReply
	pattern := "/liabilities/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationLiabilitiesGetLiability))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
