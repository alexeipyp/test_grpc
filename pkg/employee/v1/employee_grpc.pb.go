// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: employee.proto

package employee_v1

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
	Employee_PopulateWithAbsenceStatus_FullMethodName = "/employee_v1.Employee/PopulateWithAbsenceStatus"
)

// EmployeeClient is the client API for Employee service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmployeeClient interface {
	PopulateWithAbsenceStatus(ctx context.Context, in *PopulateRequest, opts ...grpc.CallOption) (*PopulateResponse, error)
}

type employeeClient struct {
	cc grpc.ClientConnInterface
}

func NewEmployeeClient(cc grpc.ClientConnInterface) EmployeeClient {
	return &employeeClient{cc}
}

func (c *employeeClient) PopulateWithAbsenceStatus(ctx context.Context, in *PopulateRequest, opts ...grpc.CallOption) (*PopulateResponse, error) {
	out := new(PopulateResponse)
	err := c.cc.Invoke(ctx, Employee_PopulateWithAbsenceStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmployeeServer is the server API for Employee service.
// All implementations must embed UnimplementedEmployeeServer
// for forward compatibility
type EmployeeServer interface {
	PopulateWithAbsenceStatus(context.Context, *PopulateRequest) (*PopulateResponse, error)
	mustEmbedUnimplementedEmployeeServer()
}

// UnimplementedEmployeeServer must be embedded to have forward compatible implementations.
type UnimplementedEmployeeServer struct {
}

func (UnimplementedEmployeeServer) PopulateWithAbsenceStatus(context.Context, *PopulateRequest) (*PopulateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PopulateWithAbsenceStatus not implemented")
}
func (UnimplementedEmployeeServer) mustEmbedUnimplementedEmployeeServer() {}

// UnsafeEmployeeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmployeeServer will
// result in compilation errors.
type UnsafeEmployeeServer interface {
	mustEmbedUnimplementedEmployeeServer()
}

func RegisterEmployeeServer(s grpc.ServiceRegistrar, srv EmployeeServer) {
	s.RegisterService(&Employee_ServiceDesc, srv)
}

func _Employee_PopulateWithAbsenceStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PopulateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmployeeServer).PopulateWithAbsenceStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Employee_PopulateWithAbsenceStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmployeeServer).PopulateWithAbsenceStatus(ctx, req.(*PopulateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Employee_ServiceDesc is the grpc.ServiceDesc for Employee service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Employee_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "employee_v1.Employee",
	HandlerType: (*EmployeeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PopulateWithAbsenceStatus",
			Handler:    _Employee_PopulateWithAbsenceStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "employee.proto",
}
