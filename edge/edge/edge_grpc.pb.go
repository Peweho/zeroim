// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: edge.proto

package edge

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

// EdgeClient is the client API for Edge service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EdgeClient interface {
	Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type edgeClient struct {
	cc grpc.ClientConnInterface
}

func NewEdgeClient(cc grpc.ClientConnInterface) EdgeClient {
	return &edgeClient{cc}
}

func (c *edgeClient) Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/edge.Edge/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EdgeServer is the server API for Edge service.
// All implementations must embed UnimplementedEdgeServer
// for forward compatibility
type EdgeServer interface {
	Ping(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedEdgeServer()
}

// UnimplementedEdgeServer must be embedded to have forward compatible implementations.
type UnimplementedEdgeServer struct {
}

func (UnimplementedEdgeServer) Ping(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedEdgeServer) mustEmbedUnimplementedEdgeServer() {}

// UnsafeEdgeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EdgeServer will
// result in compilation errors.
type UnsafeEdgeServer interface {
	mustEmbedUnimplementedEdgeServer()
}

func RegisterEdgeServer(s grpc.ServiceRegistrar, srv EdgeServer) {
	s.RegisterService(&Edge_ServiceDesc, srv)
}

func _Edge_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/edge.Edge/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeServer).Ping(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Edge_ServiceDesc is the grpc.ServiceDesc for Edge service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Edge_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "edge.Edge",
	HandlerType: (*EdgeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Edge_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "edge.proto",
}
