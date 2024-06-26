// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: v1/b2b/lead.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Companies_GetCompany_FullMethodName    = "/api.v1.Companies/GetCompany"
	Companies_CreateCompany_FullMethodName = "/api.v1.Companies/CreateCompany"
	Companies_UpdateCompany_FullMethodName = "/api.v1.Companies/UpdateCompany"
	Companies_DeleteCompany_FullMethodName = "/api.v1.Companies/DeleteCompany"
	Companies_ListCompanys_FullMethodName  = "/api.v1.Companies/ListCompanys"
)

// CompaniesClient is the client API for Companies service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CompaniesClient interface {
	GetCompany(ctx context.Context, in *GetCompanyRequest, opts ...grpc.CallOption) (*GetCompanyReply, error)
	CreateCompany(ctx context.Context, in *CreateCompanyRequest, opts ...grpc.CallOption) (*CreateCompanyReply, error)
	UpdateCompany(ctx context.Context, in *UpdateCompanyRequest, opts ...grpc.CallOption) (*UpdateCompanyReply, error)
	DeleteCompany(ctx context.Context, in *DeleteCompanyRequest, opts ...grpc.CallOption) (*DeleteCompanyReply, error)
	ListCompanys(ctx context.Context, in *ListCompaniesRequest, opts ...grpc.CallOption) (*ListCompaniesReply, error)
}

type companiesClient struct {
	cc grpc.ClientConnInterface
}

func NewCompaniesClient(cc grpc.ClientConnInterface) CompaniesClient {
	return &companiesClient{cc}
}

func (c *companiesClient) GetCompany(ctx context.Context, in *GetCompanyRequest, opts ...grpc.CallOption) (*GetCompanyReply, error) {
	out := new(GetCompanyReply)
	err := c.cc.Invoke(ctx, Companies_GetCompany_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *companiesClient) CreateCompany(ctx context.Context, in *CreateCompanyRequest, opts ...grpc.CallOption) (*CreateCompanyReply, error) {
	out := new(CreateCompanyReply)
	err := c.cc.Invoke(ctx, Companies_CreateCompany_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *companiesClient) UpdateCompany(ctx context.Context, in *UpdateCompanyRequest, opts ...grpc.CallOption) (*UpdateCompanyReply, error) {
	out := new(UpdateCompanyReply)
	err := c.cc.Invoke(ctx, Companies_UpdateCompany_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *companiesClient) DeleteCompany(ctx context.Context, in *DeleteCompanyRequest, opts ...grpc.CallOption) (*DeleteCompanyReply, error) {
	out := new(DeleteCompanyReply)
	err := c.cc.Invoke(ctx, Companies_DeleteCompany_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *companiesClient) ListCompanys(ctx context.Context, in *ListCompaniesRequest, opts ...grpc.CallOption) (*ListCompaniesReply, error) {
	out := new(ListCompaniesReply)
	err := c.cc.Invoke(ctx, Companies_ListCompanys_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompaniesServer is the server API for Companies service.
// All implementations must embed UnimplementedCompaniesServer
// for forward compatibility
type CompaniesServer interface {
	GetCompany(context.Context, *GetCompanyRequest) (*GetCompanyReply, error)
	CreateCompany(context.Context, *CreateCompanyRequest) (*CreateCompanyReply, error)
	UpdateCompany(context.Context, *UpdateCompanyRequest) (*UpdateCompanyReply, error)
	DeleteCompany(context.Context, *DeleteCompanyRequest) (*DeleteCompanyReply, error)
	ListCompanys(context.Context, *ListCompaniesRequest) (*ListCompaniesReply, error)
	mustEmbedUnimplementedCompaniesServer()
}

// UnimplementedCompaniesServer must be embedded to have forward compatible implementations.
type UnimplementedCompaniesServer struct {
}

func (UnimplementedCompaniesServer) GetCompany(context.Context, *GetCompanyRequest) (*GetCompanyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCompany not implemented")
}
func (UnimplementedCompaniesServer) CreateCompany(context.Context, *CreateCompanyRequest) (*CreateCompanyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCompany not implemented")
}
func (UnimplementedCompaniesServer) UpdateCompany(context.Context, *UpdateCompanyRequest) (*UpdateCompanyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCompany not implemented")
}
func (UnimplementedCompaniesServer) DeleteCompany(context.Context, *DeleteCompanyRequest) (*DeleteCompanyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCompany not implemented")
}
func (UnimplementedCompaniesServer) ListCompanys(context.Context, *ListCompaniesRequest) (*ListCompaniesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCompanys not implemented")
}
func (UnimplementedCompaniesServer) mustEmbedUnimplementedCompaniesServer() {}

// UnsafeCompaniesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CompaniesServer will
// result in compilation errors.
type UnsafeCompaniesServer interface {
	mustEmbedUnimplementedCompaniesServer()
}

func RegisterCompaniesServer(s grpc.ServiceRegistrar, srv CompaniesServer) {
	s.RegisterService(&Companies_ServiceDesc, srv)
}

func _Companies_GetCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCompanyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompaniesServer).GetCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Companies_GetCompany_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompaniesServer).GetCompany(ctx, req.(*GetCompanyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Companies_CreateCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCompanyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompaniesServer).CreateCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Companies_CreateCompany_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompaniesServer).CreateCompany(ctx, req.(*CreateCompanyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Companies_UpdateCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCompanyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompaniesServer).UpdateCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Companies_UpdateCompany_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompaniesServer).UpdateCompany(ctx, req.(*UpdateCompanyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Companies_DeleteCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCompanyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompaniesServer).DeleteCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Companies_DeleteCompany_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompaniesServer).DeleteCompany(ctx, req.(*DeleteCompanyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Companies_ListCompanys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCompaniesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompaniesServer).ListCompanys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Companies_ListCompanys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompaniesServer).ListCompanys(ctx, req.(*ListCompaniesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Companies_ServiceDesc is the grpc.ServiceDesc for Companies service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Companies_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.Companies",
	HandlerType: (*CompaniesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCompany",
			Handler:    _Companies_GetCompany_Handler,
		},
		{
			MethodName: "CreateCompany",
			Handler:    _Companies_CreateCompany_Handler,
		},
		{
			MethodName: "UpdateCompany",
			Handler:    _Companies_UpdateCompany_Handler,
		},
		{
			MethodName: "DeleteCompany",
			Handler:    _Companies_DeleteCompany_Handler,
		},
		{
			MethodName: "ListCompanys",
			Handler:    _Companies_ListCompanys_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/b2b/lead.proto",
}

const (
	Leads_CreateLead_FullMethodName = "/api.v1.Leads/CreateLead"
	Leads_UpdateLead_FullMethodName = "/api.v1.Leads/UpdateLead"
	Leads_DeleteLead_FullMethodName = "/api.v1.Leads/DeleteLead"
	Leads_GetLead_FullMethodName    = "/api.v1.Leads/GetLead"
	Leads_ListLeads_FullMethodName  = "/api.v1.Leads/ListLeads"
)

// LeadsClient is the client API for Leads service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LeadsClient interface {
	CreateLead(ctx context.Context, in *CreateLeadRequest, opts ...grpc.CallOption) (*CreateLeadReply, error)
	UpdateLead(ctx context.Context, in *UpdateLeadRequest, opts ...grpc.CallOption) (*UpdateLeadReply, error)
	DeleteLead(ctx context.Context, in *DeleteLeadRequest, opts ...grpc.CallOption) (*DeleteLeadReply, error)
	GetLead(ctx context.Context, in *GetLeadRequest, opts ...grpc.CallOption) (*GetLeadReply, error)
	ListLeads(ctx context.Context, in *ListLeadsRequest, opts ...grpc.CallOption) (*ListLeadsReply, error)
}

type leadsClient struct {
	cc grpc.ClientConnInterface
}

func NewLeadsClient(cc grpc.ClientConnInterface) LeadsClient {
	return &leadsClient{cc}
}

func (c *leadsClient) CreateLead(ctx context.Context, in *CreateLeadRequest, opts ...grpc.CallOption) (*CreateLeadReply, error) {
	out := new(CreateLeadReply)
	err := c.cc.Invoke(ctx, Leads_CreateLead_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leadsClient) UpdateLead(ctx context.Context, in *UpdateLeadRequest, opts ...grpc.CallOption) (*UpdateLeadReply, error) {
	out := new(UpdateLeadReply)
	err := c.cc.Invoke(ctx, Leads_UpdateLead_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leadsClient) DeleteLead(ctx context.Context, in *DeleteLeadRequest, opts ...grpc.CallOption) (*DeleteLeadReply, error) {
	out := new(DeleteLeadReply)
	err := c.cc.Invoke(ctx, Leads_DeleteLead_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leadsClient) GetLead(ctx context.Context, in *GetLeadRequest, opts ...grpc.CallOption) (*GetLeadReply, error) {
	out := new(GetLeadReply)
	err := c.cc.Invoke(ctx, Leads_GetLead_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leadsClient) ListLeads(ctx context.Context, in *ListLeadsRequest, opts ...grpc.CallOption) (*ListLeadsReply, error) {
	out := new(ListLeadsReply)
	err := c.cc.Invoke(ctx, Leads_ListLeads_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LeadsServer is the server API for Leads service.
// All implementations must embed UnimplementedLeadsServer
// for forward compatibility
type LeadsServer interface {
	CreateLead(context.Context, *CreateLeadRequest) (*CreateLeadReply, error)
	UpdateLead(context.Context, *UpdateLeadRequest) (*UpdateLeadReply, error)
	DeleteLead(context.Context, *DeleteLeadRequest) (*DeleteLeadReply, error)
	GetLead(context.Context, *GetLeadRequest) (*GetLeadReply, error)
	ListLeads(context.Context, *ListLeadsRequest) (*ListLeadsReply, error)
	mustEmbedUnimplementedLeadsServer()
}

// UnimplementedLeadsServer must be embedded to have forward compatible implementations.
type UnimplementedLeadsServer struct {
}

func (UnimplementedLeadsServer) CreateLead(context.Context, *CreateLeadRequest) (*CreateLeadReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLead not implemented")
}
func (UnimplementedLeadsServer) UpdateLead(context.Context, *UpdateLeadRequest) (*UpdateLeadReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLead not implemented")
}
func (UnimplementedLeadsServer) DeleteLead(context.Context, *DeleteLeadRequest) (*DeleteLeadReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLead not implemented")
}
func (UnimplementedLeadsServer) GetLead(context.Context, *GetLeadRequest) (*GetLeadReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLead not implemented")
}
func (UnimplementedLeadsServer) ListLeads(context.Context, *ListLeadsRequest) (*ListLeadsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLeads not implemented")
}
func (UnimplementedLeadsServer) mustEmbedUnimplementedLeadsServer() {}

// UnsafeLeadsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LeadsServer will
// result in compilation errors.
type UnsafeLeadsServer interface {
	mustEmbedUnimplementedLeadsServer()
}

func RegisterLeadsServer(s grpc.ServiceRegistrar, srv LeadsServer) {
	s.RegisterService(&Leads_ServiceDesc, srv)
}

func _Leads_CreateLead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLeadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeadsServer).CreateLead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leads_CreateLead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeadsServer).CreateLead(ctx, req.(*CreateLeadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Leads_UpdateLead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLeadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeadsServer).UpdateLead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leads_UpdateLead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeadsServer).UpdateLead(ctx, req.(*UpdateLeadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Leads_DeleteLead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteLeadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeadsServer).DeleteLead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leads_DeleteLead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeadsServer).DeleteLead(ctx, req.(*DeleteLeadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Leads_GetLead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLeadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeadsServer).GetLead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leads_GetLead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeadsServer).GetLead(ctx, req.(*GetLeadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Leads_ListLeads_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLeadsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeadsServer).ListLeads(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leads_ListLeads_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeadsServer).ListLeads(ctx, req.(*ListLeadsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Leads_ServiceDesc is the grpc.ServiceDesc for Leads service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Leads_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.Leads",
	HandlerType: (*LeadsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateLead",
			Handler:    _Leads_CreateLead_Handler,
		},
		{
			MethodName: "UpdateLead",
			Handler:    _Leads_UpdateLead_Handler,
		},
		{
			MethodName: "DeleteLead",
			Handler:    _Leads_DeleteLead_Handler,
		},
		{
			MethodName: "GetLead",
			Handler:    _Leads_GetLead_Handler,
		},
		{
			MethodName: "ListLeads",
			Handler:    _Leads_ListLeads_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/b2b/lead.proto",
}
