// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.26.1
// source: v1/finance.proto

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

const OperationFinanceGetSenateLobbying = "/api.v1.Finance/GetSenateLobbying"
const OperationFinanceGetStockQuote = "/api.v1.Finance/GetStockQuote"
const OperationFinanceGetUSASpending = "/api.v1.Finance/GetUSASpending"

type FinanceHTTPServer interface {
	GetSenateLobbying(context.Context, *GetSenateLobbyingRequest) (*GetSenateLobbyingReply, error)
	GetStockQuote(context.Context, *GetStockQuoteRequest) (*GetStockQuoteReply, error)
	GetUSASpending(context.Context, *GetUSASpendingRequest) (*GetUSASpendingReply, error)
}

func RegisterFinanceHTTPServer(s *http.Server, srv FinanceHTTPServer) {
	r := s.Route("/")
	r.GET("/quote/{symbol}", _Finance_GetStockQuote0_HTTP_Handler(srv))
	r.GET("/stock/usa-spending/{symbol}", _Finance_GetUSASpending0_HTTP_Handler(srv))
	r.GET("/stock/lobbying/{symbol}", _Finance_GetSenateLobbying0_HTTP_Handler(srv))
}

func _Finance_GetStockQuote0_HTTP_Handler(srv FinanceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetStockQuoteRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationFinanceGetStockQuote)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetStockQuote(ctx, req.(*GetStockQuoteRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetStockQuoteReply)
		return ctx.Result(200, reply)
	}
}

func _Finance_GetUSASpending0_HTTP_Handler(srv FinanceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetUSASpendingRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationFinanceGetUSASpending)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetUSASpending(ctx, req.(*GetUSASpendingRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetUSASpendingReply)
		return ctx.Result(200, reply)
	}
}

func _Finance_GetSenateLobbying0_HTTP_Handler(srv FinanceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetSenateLobbyingRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationFinanceGetSenateLobbying)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetSenateLobbying(ctx, req.(*GetSenateLobbyingRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetSenateLobbyingReply)
		return ctx.Result(200, reply)
	}
}

type FinanceHTTPClient interface {
	GetSenateLobbying(ctx context.Context, req *GetSenateLobbyingRequest, opts ...http.CallOption) (rsp *GetSenateLobbyingReply, err error)
	GetStockQuote(ctx context.Context, req *GetStockQuoteRequest, opts ...http.CallOption) (rsp *GetStockQuoteReply, err error)
	GetUSASpending(ctx context.Context, req *GetUSASpendingRequest, opts ...http.CallOption) (rsp *GetUSASpendingReply, err error)
}

type FinanceHTTPClientImpl struct {
	cc *http.Client
}

func NewFinanceHTTPClient(client *http.Client) FinanceHTTPClient {
	return &FinanceHTTPClientImpl{client}
}

func (c *FinanceHTTPClientImpl) GetSenateLobbying(ctx context.Context, in *GetSenateLobbyingRequest, opts ...http.CallOption) (*GetSenateLobbyingReply, error) {
	var out GetSenateLobbyingReply
	pattern := "/stock/lobbying/{symbol}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationFinanceGetSenateLobbying))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *FinanceHTTPClientImpl) GetStockQuote(ctx context.Context, in *GetStockQuoteRequest, opts ...http.CallOption) (*GetStockQuoteReply, error) {
	var out GetStockQuoteReply
	pattern := "/quote/{symbol}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationFinanceGetStockQuote))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *FinanceHTTPClientImpl) GetUSASpending(ctx context.Context, in *GetUSASpendingRequest, opts ...http.CallOption) (*GetUSASpendingReply, error) {
	var out GetUSASpendingReply
	pattern := "/stock/usa-spending/{symbol}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationFinanceGetUSASpending))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
