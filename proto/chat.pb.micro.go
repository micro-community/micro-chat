// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: chat.proto

package chat

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/micro/v3/service/api"
	client "github.com/micro/micro/v3/service/client"
	server "github.com/micro/micro/v3/service/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Chat service

func NewChatEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Chat service

type ChatService interface {
	//支持幂等调用的接口,对相同的人用户创建会话，默认总是返回同一个chat，除非指定强制创建为新会话的参数
	New(ctx context.Context, in *NewRequest, opts ...client.CallOption) (*NewResponse, error)
	//删除一个聊天，除非指定彻底销毁
	Remove(ctx context.Context, in *RemoveRequest, opts ...client.CallOption) (*RemoveResponse, error)
	// 查询某个chat中的历史消息
	History(ctx context.Context, in *HistoryRequest, opts ...client.CallOption) (*HistoryResponse, error)
	// Send a single message to the chat，当发送消息给服务端时，消息会被添加到历史消息，同时转发给其他的连接用户。
	// 客户端的消息参数中应当携带chat_id和user_id 以便服务端识别会话和用户
	Send(ctx context.Context, in *SendRequest, opts ...client.CallOption) (*SendResponse, error)
	// 双向stream的方式，连接到某一个会话(或者聊天室)，
	Connect(ctx context.Context, opts ...client.CallOption) (Chat_ConnectService, error)
}

type chatService struct {
	c    client.Client
	name string
}

func NewChatService(name string, c client.Client) ChatService {
	return &chatService{
		c:    c,
		name: name,
	}
}

func (c *chatService) New(ctx context.Context, in *NewRequest, opts ...client.CallOption) (*NewResponse, error) {
	req := c.c.NewRequest(c.name, "Chat.New", in)
	out := new(NewResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) Remove(ctx context.Context, in *RemoveRequest, opts ...client.CallOption) (*RemoveResponse, error) {
	req := c.c.NewRequest(c.name, "Chat.Remove", in)
	out := new(RemoveResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) History(ctx context.Context, in *HistoryRequest, opts ...client.CallOption) (*HistoryResponse, error) {
	req := c.c.NewRequest(c.name, "Chat.History", in)
	out := new(HistoryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) Send(ctx context.Context, in *SendRequest, opts ...client.CallOption) (*SendResponse, error) {
	req := c.c.NewRequest(c.name, "Chat.Send", in)
	out := new(SendResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) Connect(ctx context.Context, opts ...client.CallOption) (Chat_ConnectService, error) {
	req := c.c.NewRequest(c.name, "Chat.Connect", &Message{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &chatServiceConnect{stream}, nil
}

type Chat_ConnectService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Message) error
	Recv() (*Message, error)
}

type chatServiceConnect struct {
	stream client.Stream
}

func (x *chatServiceConnect) Close() error {
	return x.stream.Close()
}

func (x *chatServiceConnect) Context() context.Context {
	return x.stream.Context()
}

func (x *chatServiceConnect) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *chatServiceConnect) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *chatServiceConnect) Send(m *Message) error {
	return x.stream.Send(m)
}

func (x *chatServiceConnect) Recv() (*Message, error) {
	m := new(Message)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Chat service

type ChatHandler interface {
	//支持幂等调用的接口,对相同的人用户创建会话，默认总是返回同一个chat，除非指定强制创建为新会话的参数
	New(context.Context, *NewRequest, *NewResponse) error
	//删除一个聊天，除非指定彻底销毁
	Remove(context.Context, *RemoveRequest, *RemoveResponse) error
	// 查询某个chat中的历史消息
	History(context.Context, *HistoryRequest, *HistoryResponse) error
	// Send a single message to the chat，当发送消息给服务端时，消息会被添加到历史消息，同时转发给其他的连接用户。
	// 客户端的消息参数中应当携带chat_id和user_id 以便服务端识别会话和用户
	Send(context.Context, *SendRequest, *SendResponse) error
	// 双向stream的方式，连接到某一个会话(或者聊天室)，
	Connect(context.Context, Chat_ConnectStream) error
}

func RegisterChatHandler(s server.Server, hdlr ChatHandler, opts ...server.HandlerOption) error {
	type chat interface {
		New(ctx context.Context, in *NewRequest, out *NewResponse) error
		Remove(ctx context.Context, in *RemoveRequest, out *RemoveResponse) error
		History(ctx context.Context, in *HistoryRequest, out *HistoryResponse) error
		Send(ctx context.Context, in *SendRequest, out *SendResponse) error
		Connect(ctx context.Context, stream server.Stream) error
	}
	type Chat struct {
		chat
	}
	h := &chatHandler{hdlr}
	return s.Handle(s.NewHandler(&Chat{h}, opts...))
}

type chatHandler struct {
	ChatHandler
}

func (h *chatHandler) New(ctx context.Context, in *NewRequest, out *NewResponse) error {
	return h.ChatHandler.New(ctx, in, out)
}

func (h *chatHandler) Remove(ctx context.Context, in *RemoveRequest, out *RemoveResponse) error {
	return h.ChatHandler.Remove(ctx, in, out)
}

func (h *chatHandler) History(ctx context.Context, in *HistoryRequest, out *HistoryResponse) error {
	return h.ChatHandler.History(ctx, in, out)
}

func (h *chatHandler) Send(ctx context.Context, in *SendRequest, out *SendResponse) error {
	return h.ChatHandler.Send(ctx, in, out)
}

func (h *chatHandler) Connect(ctx context.Context, stream server.Stream) error {
	return h.ChatHandler.Connect(ctx, &chatConnectStream{stream})
}

type Chat_ConnectStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Message) error
	Recv() (*Message, error)
}

type chatConnectStream struct {
	stream server.Stream
}

func (x *chatConnectStream) Close() error {
	return x.stream.Close()
}

func (x *chatConnectStream) Context() context.Context {
	return x.stream.Context()
}

func (x *chatConnectStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *chatConnectStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *chatConnectStream) Send(m *Message) error {
	return x.stream.Send(m)
}

func (x *chatConnectStream) Recv() (*Message, error) {
	m := new(Message)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
