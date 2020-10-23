package handler

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/micro-community/micro-chat/proto"
	"github.com/micro/micro/v3/service/errors"
)

// Send a single message to the chat, designed for ease of use via the API / CLI
func (c *Chat) Send(ctx context.Context, req *pb.SendRequest, rsp *pb.SendResponse) error {
	// validate the request
	if len(req.ChatId) == 0 {
		return errors.BadRequest("chat.Send.MissingChatID", "ChatID is missing")
	}
	if len(req.UserId) == 0 {
		return errors.BadRequest("chat.Send.MissingUserID", "UserID is missing")
	}
	if len(req.Text) == 0 {
		return errors.BadRequest("chat.Send.MissingText", "Text is missing")
	}

	// construct the message
	msg := &pb.Message{
		Id:       uuid.New().String(),
		ClientId: req.ClientId,
		ChatId:   req.ChatId,
		UserId:   req.UserId,
		Subject:  req.Subject,
		Text:     req.Text,
	}

	// default the client id if not provided
	if len(msg.ClientId) == 0 {
		msg.ClientId = uuid.New().String()
	}

	// create the message
	return c.createMessage(msg)
}
