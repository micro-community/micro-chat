package handler

import (
	"context"

	pb "github.com/micro-community/micro-chat/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

// Remove a chat, is
func (c *Chat) Remove(ctx context.Context, req *pb.RemoveRequest, rsp *pb.RemoveResponse) error {

	// validate the request
	if len(req.ChatId) == 0 {
		return errors.BadRequest("chat.History.MissingChatID", "ChatID is missing")
	}

	// lookup the chat from the store to ensure it's valid
	if _, err := store.Read(chatStoreKeyPrefix + req.ChatId); err == store.ErrNotFound {
		return errors.BadRequest("chat.History.InvalidChatID", "Chat not found with this ID")
	} else if err != nil {
		logger.Errorf("Error reading from the store. Chat ID: %v. Error: %v", req.ChatId, err)
		return errors.InternalServerError("chat.History.Unknown", "Error reading from the store")
	}

	return nil
}
