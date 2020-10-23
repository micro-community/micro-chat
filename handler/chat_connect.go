package handler

import (
	"context"

	pb "github.com/micro-community/micro-chat/proto"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

// Connect to a chat using a bidirectional stream enabling the client to send and recieve messages
// over a single RPC. When a message is sent on the stream, it will be added to the chat history
// and sent to the other connected users. When opening the connection, the client should provide
// the chat_id and user_id in the context so the server knows which messages to stream.
func (c *Chat) Connect(ctx context.Context, stream pb.Chat_ConnectStream) error {
	// the client passed the chat id and user id in the request context. we'll load that information
	// now and validate it. If any information is missing we'll return a BadRequest error to the client
	userID, ok := metadata.Get(ctx, "user-id")
	if !ok {
		return errors.BadRequest("chat.Connect.MissingUserID", "UserID missing in context")
	}
	chatID, ok := metadata.Get(ctx, "chat-id")
	if !ok {
		return errors.BadRequest("chat.Connect.MissingChatID", "ChatId missing in context")
	}

	// lookup the chat from the store to ensure it's valid
	if _, err := store.Read(chatStoreKeyPrefix + chatID); err == store.ErrNotFound {
		return errors.BadRequest("chat.Connect.InvalidChatID", "Chat not found with this ID")
	} else if err != nil {
		logger.Errorf("Error reading from the store. Chat ID: %v. Error: %v", chatID, err)
		return errors.InternalServerError("chat.Connect.Unknown", "Error reading from the store")
	}

	// as per the New and Connect functions, at this point in a real world application we would
	// authorize the request to ensure the authenticated user is part of the chat they're attempting
	// to read the history of

	// create a new context which can be cancelled, in the case either the consumer of publisher errors
	// we don't want one to keep running in the background
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// create a channel to send errors on, because the subscriber / publisher will run in seperate go-
	// routines, they need a way of returning errors to the client
	errChan := make(chan error)

	// create an event stream to consume messages posted by other users into the chat. we'll use the
	// user id as a queue to ensure each user recieves the message
	evStream, err := events.Consume(chatEventKeyPrefix+chatID, events.WithGroup(userID))
	if err != nil {
		logger.Errorf("Error streaming events. Chat ID: %v. Error: %v", chatID, err)
		return errors.InternalServerError("chat.Connect.Unknown", "Error connecting to the event stream")
	}
	go func() {
		for {
			select {
			case <-cancelCtx.Done():
				// the context has been cancelled or timed out, stop subscribing to new messages
				return
			case ev := <-evStream:
				// received a message, unmarshal it into a message struct. if an error occurs log it and
				// cancel the context
				var msg pb.Message
				if err := ev.Unmarshal(&msg); err != nil {
					logger.Errorf("Error unmarshaling message. ChatID: %v. Error: %v", chatID, err)
					errChan <- err
					return
				}

				// ignore any messages published by the current user
				if msg.UserId == userID {
					continue
				}

				// publish the message to the stream
				if err := stream.Send(&msg); err != nil {
					logger.Errorf("Error sending message to stream. ChatID: %v. Message ID: %v. Error: %v", chatID, msg.Id, err)
					errChan <- err
					return
				}
			}
		}
	}()

	// transform the stream.Recv into a channel which can be used in the select statement below
	msgChan := make(chan *pb.Message)
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				errChan <- err
				close(msgChan)
				return
			}
			msgChan <- msg
		}
	}()

	for {
		select {
		case <-cancelCtx.Done():
			// the context has been cancelled or timed out, stop subscribing to new messages
			return nil
		case err := <-errChan:
			// an error occurred in another goroutine, terminate the stream
			return err
		case msg := <-msgChan:
			// set the defaults
			msg.UserId = userID
			msg.ChatId = chatID

			// create the message
			if err := c.createMessage(msg); err != nil {
				return err
			}
		}
	}
}
