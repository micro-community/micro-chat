package handler

// it's standard to import the services own proto under the alias pb

import (
	"github.com/micro-community/micro-chat/model"
	pb "github.com/micro-community/micro-chat/proto"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/store"
)

const (
	chatStoreKeyPrefix    = "chats/"
	chatEventKeyPrefix    = "chats/"
	messageStoreKeyPrefix = "messages/"
)

// Chat satisfies the ChatHandler interface. You can see this interface defined in chat.pb.micro.go
type Chat struct {
	Namespace string
	repo      *model.Repository
}

//New Return Chat Handler
func New(namespace string) *Chat {
	return &Chat{
		Namespace: namespace,
		repo:      model.NewRepository("messsages"),
	}
}

// createMessage creates a message in the event stream. It handles the
// logic for ensuring client id is unique.
func (c *Chat) createMessage(msg *pb.Message) error {
	// a message was received from the client. validate it hasn't been received before
	if _, err := store.Read(messageStoreKeyPrefix + msg.ClientId); err == nil {
		// the message has already been processed
		return nil
	} else if err != store.ErrNotFound {
		// an unexpected error occurred
		return err
	}

	// send the message to the event stream
	if err := events.Publish(chatEventKeyPrefix+msg.ChatId, msg); err != nil {
		return err
	}

	// record the messages client id
	if err := store.Write(&store.Record{Key: messageStoreKeyPrefix + msg.ClientId}); err != nil {
		return err
	}

	return nil
}
