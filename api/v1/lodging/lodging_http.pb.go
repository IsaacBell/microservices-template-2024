// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.26.1
// source: v1/lodging/lodging.proto

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

const OperationLodgingCreateLodging = "/api.v1.lodging.Lodging/CreateLodging"
const OperationLodgingDeleteLodging = "/api.v1.lodging.Lodging/DeleteLodging"
const OperationLodgingGetLodging = "/api.v1.lodging.Lodging/GetLodging"
const OperationLodgingListLodging = "/api.v1.lodging.Lodging/ListLodging"
const OperationLodgingRealtorStats = "/api.v1.lodging.Lodging/RealtorStats"
const OperationLodgingSearchLodging = "/api.v1.lodging.Lodging/SearchLodging"
const OperationLodgingUpdateLodging = "/api.v1.lodging.Lodging/UpdateLodging"

type LodgingHTTPServer interface {
	CreateLodging(context.Context, *CreateLodgingRequest) (*CreateLodgingReply, error)
	DeleteLodging(context.Context, *DeleteLodgingRequest) (*DeleteLodgingReply, error)
	GetLodging(context.Context, *GetLodgingRequest) (*GetLodgingReply, error)
	ListLodging(context.Context, *ListLodgingRequest) (*ListLodgingReply, error)
	RealtorStats(context.Context, *RealtorStatsRequest) (*RealtorStatsReply, error)
	SearchLodging(context.Context, *SearchLodgingRequest) (*SearchLodgingReply, error)
	UpdateLodging(context.Context, *UpdateLodgingRequest) (*UpdateLodgingReply, error)
}

func RegisterLodgingHTTPServer(s *http.Server, srv LodgingHTTPServer) {
	r := s.Route("/")
	r.POST("/properties", _Lodging_CreateLodging0_HTTP_Handler(srv))
	r.PUT("/properties/{property.id}", _Lodging_UpdateLodging0_HTTP_Handler(srv))
	r.DELETE("/properties/{id}", _Lodging_DeleteLodging0_HTTP_Handler(srv))
	r.GET("/properties/{id}", _Lodging_GetLodging0_HTTP_Handler(srv))
	r.GET("/properties", _Lodging_ListLodging0_HTTP_Handler(srv))
	r.GET("/properties/search", _Lodging_SearchLodging0_HTTP_Handler(srv))
	r.GET("/properties/realtor_stats", _Lodging_RealtorStats0_HTTP_Handler(srv))
}

func _Lodging_CreateLodging0_HTTP_Handler(srv LodgingHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateLodgingRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLodgingCreateLodging)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateLodging(ctx, req.(*CreateLodgingRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateLodgingReply)
		return ctx.Result(200, reply)
	}
}

func _Lodging_UpdateLodging0_HTTP_Handler(srv LodgingHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateLodgingRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLodgingUpdateLodging)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateLodging(ctx, req.(*UpdateLodgingRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateLodgingReply)
		return ctx.Result(200, reply)
	}
}

func _Lodging_DeleteLodging0_HTTP_Handler(srv LodgingHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteLodgingRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLodgingDeleteLodging)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteLodging(ctx, req.(*DeleteLodgingRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteLodgingReply)
		return ctx.Result(200, reply)
	}
}

func _Lodging_GetLodging0_HTTP_Handler(srv LodgingHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetLodgingRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLodgingGetLodging)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetLodging(ctx, req.(*GetLodgingRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetLodgingReply)
		return ctx.Result(200, reply)
	}
}

func _Lodging_ListLodging0_HTTP_Handler(srv LodgingHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListLodgingRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLodgingListLodging)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListLodging(ctx, req.(*ListLodgingRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListLodgingReply)
		return ctx.Result(200, reply)
	}
}

func _Lodging_SearchLodging0_HTTP_Handler(srv LodgingHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SearchLodgingRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLodgingSearchLodging)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SearchLodging(ctx, req.(*SearchLodgingRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SearchLodgingReply)
		return ctx.Result(200, reply)
	}
}

func _Lodging_RealtorStats0_HTTP_Handler(srv LodgingHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RealtorStatsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLodgingRealtorStats)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.RealtorStats(ctx, req.(*RealtorStatsRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RealtorStatsReply)
		return ctx.Result(200, reply)
	}
}

type LodgingHTTPClient interface {
	CreateLodging(ctx context.Context, req *CreateLodgingRequest, opts ...http.CallOption) (rsp *CreateLodgingReply, err error)
	DeleteLodging(ctx context.Context, req *DeleteLodgingRequest, opts ...http.CallOption) (rsp *DeleteLodgingReply, err error)
	GetLodging(ctx context.Context, req *GetLodgingRequest, opts ...http.CallOption) (rsp *GetLodgingReply, err error)
	ListLodging(ctx context.Context, req *ListLodgingRequest, opts ...http.CallOption) (rsp *ListLodgingReply, err error)
	RealtorStats(ctx context.Context, req *RealtorStatsRequest, opts ...http.CallOption) (rsp *RealtorStatsReply, err error)
	SearchLodging(ctx context.Context, req *SearchLodgingRequest, opts ...http.CallOption) (rsp *SearchLodgingReply, err error)
	UpdateLodging(ctx context.Context, req *UpdateLodgingRequest, opts ...http.CallOption) (rsp *UpdateLodgingReply, err error)
}

type LodgingHTTPClientImpl struct {
	cc *http.Client
}

func NewLodgingHTTPClient(client *http.Client) LodgingHTTPClient {
	return &LodgingHTTPClientImpl{client}
}

func (c *LodgingHTTPClientImpl) CreateLodging(ctx context.Context, in *CreateLodgingRequest, opts ...http.CallOption) (*CreateLodgingReply, error) {
	var out CreateLodgingReply
	pattern := "/properties"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationLodgingCreateLodging))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LodgingHTTPClientImpl) DeleteLodging(ctx context.Context, in *DeleteLodgingRequest, opts ...http.CallOption) (*DeleteLodgingReply, error) {
	var out DeleteLodgingReply
	pattern := "/properties/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationLodgingDeleteLodging))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LodgingHTTPClientImpl) GetLodging(ctx context.Context, in *GetLodgingRequest, opts ...http.CallOption) (*GetLodgingReply, error) {
	var out GetLodgingReply
	pattern := "/properties/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationLodgingGetLodging))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LodgingHTTPClientImpl) ListLodging(ctx context.Context, in *ListLodgingRequest, opts ...http.CallOption) (*ListLodgingReply, error) {
	var out ListLodgingReply
	pattern := "/properties"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationLodgingListLodging))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LodgingHTTPClientImpl) RealtorStats(ctx context.Context, in *RealtorStatsRequest, opts ...http.CallOption) (*RealtorStatsReply, error) {
	var out RealtorStatsReply
	pattern := "/properties/realtor_stats"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationLodgingRealtorStats))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LodgingHTTPClientImpl) SearchLodging(ctx context.Context, in *SearchLodgingRequest, opts ...http.CallOption) (*SearchLodgingReply, error) {
	var out SearchLodgingReply
	pattern := "/properties/search"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationLodgingSearchLodging))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LodgingHTTPClientImpl) UpdateLodging(ctx context.Context, in *UpdateLodgingRequest, opts ...http.CallOption) (*UpdateLodgingReply, error) {
	var out UpdateLodgingReply
	pattern := "/properties/{property.id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationLodgingUpdateLodging))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}