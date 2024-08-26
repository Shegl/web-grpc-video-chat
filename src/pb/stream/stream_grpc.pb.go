// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: stream.proto

package stream

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
	Stream_StreamState_FullMethodName = "/Stream/StreamState"
	Stream_ChangeState_FullMethodName = "/Stream/ChangeState"
	Stream_AVStream_FullMethodName    = "/Stream/AVStream"
)

// StreamClient is the client API for Stream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamClient interface {
	StreamState(ctx context.Context, in *User, opts ...grpc.CallOption) (Stream_StreamStateClient, error)
	ChangeState(ctx context.Context, in *User, opts ...grpc.CallOption) (*Ack, error)
	AVStream(ctx context.Context, in *User, opts ...grpc.CallOption) (Stream_AVStreamClient, error)
}

type streamClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamClient(cc grpc.ClientConnInterface) StreamClient {
	return &streamClient{cc}
}

func (c *streamClient) StreamState(ctx context.Context, in *User, opts ...grpc.CallOption) (Stream_StreamStateClient, error) {
	stream, err := c.cc.NewStream(ctx, &Stream_ServiceDesc.Streams[0], Stream_StreamState_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &streamStreamStateClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Stream_StreamStateClient interface {
	Recv() (*StateMessage, error)
	grpc.ClientStream
}

type streamStreamStateClient struct {
	grpc.ClientStream
}

func (x *streamStreamStateClient) Recv() (*StateMessage, error) {
	m := new(StateMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamClient) ChangeState(ctx context.Context, in *User, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, Stream_ChangeState_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamClient) AVStream(ctx context.Context, in *User, opts ...grpc.CallOption) (Stream_AVStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Stream_ServiceDesc.Streams[1], Stream_AVStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &streamAVStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Stream_AVStreamClient interface {
	Recv() (*AVFrameData, error)
	grpc.ClientStream
}

type streamAVStreamClient struct {
	grpc.ClientStream
}

func (x *streamAVStreamClient) Recv() (*AVFrameData, error) {
	m := new(AVFrameData)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamServer is the server API for Stream service.
// All implementations must embed UnimplementedStreamServer
// for forward compatibility
type StreamServer interface {
	StreamState(*User, Stream_StreamStateServer) error
	ChangeState(context.Context, *User) (*Ack, error)
	AVStream(*User, Stream_AVStreamServer) error
	mustEmbedUnimplementedStreamServer()
}

// UnimplementedStreamServer must be embedded to have forward compatible implementations.
type UnimplementedStreamServer struct {
}

func (UnimplementedStreamServer) StreamState(*User, Stream_StreamStateServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamState not implemented")
}
func (UnimplementedStreamServer) ChangeState(context.Context, *User) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeState not implemented")
}
func (UnimplementedStreamServer) AVStream(*User, Stream_AVStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method AVStream not implemented")
}
func (UnimplementedStreamServer) mustEmbedUnimplementedStreamServer() {}

// UnsafeStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamServer will
// result in compilation errors.
type UnsafeStreamServer interface {
	mustEmbedUnimplementedStreamServer()
}

func RegisterStreamServer(s grpc.ServiceRegistrar, srv StreamServer) {
	s.RegisterService(&Stream_ServiceDesc, srv)
}

func _Stream_StreamState_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(User)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServer).StreamState(m, &streamStreamStateServer{stream})
}

type Stream_StreamStateServer interface {
	Send(*StateMessage) error
	grpc.ServerStream
}

type streamStreamStateServer struct {
	grpc.ServerStream
}

func (x *streamStreamStateServer) Send(m *StateMessage) error {
	return x.ServerStream.SendMsg(m)
}

func _Stream_ChangeState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamServer).ChangeState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Stream_ChangeState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamServer).ChangeState(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stream_AVStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(User)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServer).AVStream(m, &streamAVStreamServer{stream})
}

type Stream_AVStreamServer interface {
	Send(*AVFrameData) error
	grpc.ServerStream
}

type streamAVStreamServer struct {
	grpc.ServerStream
}

func (x *streamAVStreamServer) Send(m *AVFrameData) error {
	return x.ServerStream.SendMsg(m)
}

// Stream_ServiceDesc is the grpc.ServiceDesc for Stream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Stream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Stream",
	HandlerType: (*StreamServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ChangeState",
			Handler:    _Stream_ChangeState_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamState",
			Handler:       _Stream_StreamState_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "AVStream",
			Handler:       _Stream_AVStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "stream.proto",
}