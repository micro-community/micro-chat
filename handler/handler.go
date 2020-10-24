package handler

// it's standard to import the services own proto under the alias pb

import (
	message "github.com/micro-community/micro-chat/service"
)

const (
	chatStoreKeyPrefix    = "chats/"
	chatEventKeyPrefix    = "chats/"
	messageStoreKeyPrefix = "messages/"
)

// Chat satisfies the ChatHandler interface. You can see this interface defined in chat.pb.micro.go
type Chat struct {
	Message *message.Chat
}
