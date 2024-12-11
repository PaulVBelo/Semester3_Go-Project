// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: delivery_system.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DeliverySystem_SendBooking_FullMethodName = "/deliverysystem.DeliverySystem/SendBooking"
)

// DeliverySystemClient is the client API for DeliverySystem service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Определяем сервис
type DeliverySystemClient interface {
	SendBooking(ctx context.Context, in *BookingEvent, opts ...grpc.CallOption) (*BookingResponse, error)
}

type deliverySystemClient struct {
	cc grpc.ClientConnInterface
}

func NewDeliverySystemClient(cc grpc.ClientConnInterface) DeliverySystemClient {
	return &deliverySystemClient{cc}
}

func (c *deliverySystemClient) SendBooking(ctx context.Context, in *BookingEvent, opts ...grpc.CallOption) (*BookingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BookingResponse)
	err := c.cc.Invoke(ctx, DeliverySystem_SendBooking_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeliverySystemServer is the server API for DeliverySystem service.
// All implementations must embed UnimplementedDeliverySystemServer
// for forward compatibility.
//
// Определяем сервис
type DeliverySystemServer interface {
	SendBooking(context.Context, *BookingEvent) (*BookingResponse, error)
	mustEmbedUnimplementedDeliverySystemServer()
}

// UnimplementedDeliverySystemServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDeliverySystemServer struct{}

func (UnimplementedDeliverySystemServer) SendBooking(context.Context, *BookingEvent) (*BookingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendBooking not implemented")
}
func (UnimplementedDeliverySystemServer) mustEmbedUnimplementedDeliverySystemServer() {}
func (UnimplementedDeliverySystemServer) testEmbeddedByValue()                        {}

// UnsafeDeliverySystemServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeliverySystemServer will
// result in compilation errors.
type UnsafeDeliverySystemServer interface {
	mustEmbedUnimplementedDeliverySystemServer()
}

func RegisterDeliverySystemServer(s grpc.ServiceRegistrar, srv DeliverySystemServer) {
	// If the following call pancis, it indicates UnimplementedDeliverySystemServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DeliverySystem_ServiceDesc, srv)
}

func _DeliverySystem_SendBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookingEvent)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeliverySystemServer).SendBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeliverySystem_SendBooking_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeliverySystemServer).SendBooking(ctx, req.(*BookingEvent))
	}
	return interceptor(ctx, in, info, handler)
}

// DeliverySystem_ServiceDesc is the grpc.ServiceDesc for DeliverySystem service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeliverySystem_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "deliverysystem.DeliverySystem",
	HandlerType: (*DeliverySystemServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendBooking",
			Handler:    _DeliverySystem_SendBooking_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "delivery_system.proto",
}
