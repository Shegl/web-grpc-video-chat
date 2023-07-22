// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: chat.proto

package chat

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
	Chat_GetHistory_FullMethodName  = "/Chat/GetHistory"
	Chat_SendMessage_FullMethodName = "/Chat/SendMessage"
	Chat_Listen_FullMethodName      = "/Chat/Listen"
)

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatClient interface {
	GetHistory(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*HistoryResponse, error)
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*Empty, error)
	Listen(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (Chat_ListenClient, error)
}

type chatClient struct {
	cc grpc.ClientConnInterface
}

func NewChatClient(cc grpc.ClientConnInterface) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) GetHistory(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*HistoryResponse, error) {
	out := new(HistoryResponse)
	err := c.cc.Invoke(ctx, Chat_GetHistory_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Chat_SendMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) Listen(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (Chat_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &Chat_ServiceDesc.Streams[0], Chat_Listen_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &chatListenClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Chat_ListenClient interface {
	Recv() (*ChatMessage, error)
	grpc.ClientStream
}

type chatListenClient struct {
	grpc.ClientStream
}

func (x *chatListenClient) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatServer is the server API for Chat service.
// All implementations must embed UnimplementedChatServer
// for forward compatibility
type ChatServer interface {
	GetHistory(context.Context, *AuthRequest) (*HistoryResponse, error)
	SendMessage(context.Context, *SendMessageRequest) (*Empty, error)
	Listen(*AuthRequest, Chat_ListenServer) error
	mustEmbedUnimplementedChatServer()
}

// UnimplementedChatServer must be embedded to have forward compatible implementations.
type UnimplementedChatServer struct {
}

func (UnimplementedChatServer) GetHistory(context.Context, *AuthRequest) (*HistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHistory not implemented")
}
func (UnimplementedChatServer) SendMessage(context.Context, *SendMessageRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedChatServer) Listen(*AuthRequest, Chat_ListenServer) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}
func (UnimplementedChatServer) mustEmbedUnimplementedChatServer() {}

// UnsafeChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServer will
// result in compilation errors.
type UnsafeChatServer interface {
	mustEmbedUnimplementedChatServer()
}

func RegisterChatServer(s grpc.ServiceRegistrar, srv ChatServer) {
	s.RegisterService(&Chat_ServiceDesc, srv)
}

func _Chat_GetHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).GetHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_GetHistory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).GetHistory(ctx, req.(*AuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AuthRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChatServer).Listen(m, &chatListenServer{stream})
}

type Chat_ListenServer interface {
	Send(*ChatMessage) error
	grpc.ServerStream
}

type chatListenServer struct {
	grpc.ServerStream
}

func (x *chatListenServer) Send(m *ChatMessage) error {
	return x.ServerStream.SendMsg(m)
}

// Chat_ServiceDesc is the grpc.ServiceDesc for Chat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetHistory",
			Handler:    _Chat_GetHistory_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _Chat_SendMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Listen",
			Handler:       _Chat_Listen_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chat.proto",
}
