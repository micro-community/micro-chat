package handler

import (
	"context"

	pb "github.com/micro-community/micro-chat/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

// History returns the historical messages in a chat
func (c *Chat) History(ctx context.Context, req *pb.HistoryRequest, rsp *pb.HistoryResponse) error {
	// as per the New function, in a real world application we would authorize the request to ensure
	// the authenticated user is part of the chat they're attempting to read the history of

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

	// lookup the historical messages for the chat using the event store. lots of packages in micro
	// support options, in this case we'll pass the ReadLimit option to restrict the number of messages
	// we'll load from the events store.
	messages, err := events.Read(chatEventKeyPrefix+req.ChatId, events.ReadLimit(50))
	if err != nil {
		logger.Errorf("Error reading from the event store. Chat ID: %v. Error: %v", req.ChatId, err)
		return errors.InternalServerError("chat.History.Unknown", "Error reading from the event store")
	}

	// we've loaded the messages from the event store. next we need to serialize them and return them
	// to the client. The message is stored in the event payload, to retrieve it we need to unmarshal
	// the event into a message struct.
	rsp.Messages = make([]*pb.Message, len(messages))
	for i, ev := range messages {
		var msg pb.Message
		if err := ev.Unmarshal(&msg); err != nil {
			logger.Errorf("Error unmarshaling event: %v", err)
			return errors.InternalServerError("chat.History.Unknown", "Error unmarshaling event")
		}
		rsp.Messages[i] = &msg
	}

	return nil
}
