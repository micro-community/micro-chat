package handler

import (
	"context"
	"sort"
	"strings"

	"github.com/google/uuid"
	pb "github.com/micro-community/micro-chat/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

// New creates a chat for a group of users. The RPC is idempotent so if it's called multiple times
// for the same users, the same response will be returned. It's good practice to design APIs as
// idempotent since this enables safe retries.
func (c *Chat) New(ctx context.Context, req *pb.NewRequest, rsp *pb.NewResponse) error {
	// in a real world application we would authorize the request to ensure the authenticated user
	// is part of the chat they're attempting to create. We could do this by getting the user id from
	// auth.AccountFromContext(ctx) and then validating the presence of their id in req.UserIds. If
	// the user is not part of the request then we'd return a Forbidden error, which the micro api
	// would transform to a 403 status code.

	// validate the request
	if len(req.UserIds) == 0 {
		// Return a bad request error to the client, the first argument is a unique id which the client
		// can check for. The second argument is a human readable description. Returning the correct type
		// of error is important as it's used by the network to know if a request should be retried. Only
		// 500 (InternalServerError) and 408 (Timeout) errors are retried.
		return errors.BadRequest("chat.New.MissingUserIDs", "One or more user IDs are required")
	}

	// construct a key to identify the chat, we'll do this by sorting the user ids alphabetically and
	// then joining them. When a service calls the store, the data returned will be automatically scoped
	// to the service however it's still advised to use a prefix when writing data since this allows
	// other types of keys to be written in the future. We'll make a copy of the req.UserIds object as
	// it's a good practice to not mutate the request object.
	sortedIDs := make([]string, len(req.UserIds))
	copy(sortedIDs, req.UserIds)
	sort.Strings(sortedIDs)

	// key to lookup the chat in the store using, e.g. "chat/usera-userb-userc"
	key := chatStoreKeyPrefix + strings.Join(sortedIDs, "-")

	// read from the store to check if a chat with these users already exists
	recs, err := store.Read(key)
	if err == nil {
		// if an error wasn't returned, at least one record was found. The value returned by the store
		// is the bytes representation of the chat id. We'll convert this back into a string and return
		// it to the client.
		rsp.ChatId = string(recs[0].Value)
		return nil
	} else if err != store.ErrNotFound {
		// if no records were found then we'd expect to get a store.ErrNotFound error returned. If this
		// wasn't the case, the service could've experienced an issue connecting to the store so we should
		// log the error and return an InternalServerError to the client, indicating the request should
		// be retried
		logger.Errorf("Error reading from the store. Key: %v. Error: %v", key, err)
		return errors.InternalServerError("chat.New.Unknown", "Error reading from the store")
	}

	// no chat id was returned so we'll generate one, write it to the store and then return it to the
	// client
	chatID := uuid.New().String()
	record := store.Record{Key: chatStoreKeyPrefix + chatID, Value: []byte(chatID)}
	if err := store.Write(&record); err != nil {
		logger.Errorf("Error writing to the store. Key: %v. Error: %v", record.Key, err)
		return errors.InternalServerError("chat.New.Unknown", "Error writing to the store")
	}

	// The chat was successfully created so we'll log the event and then return the id to the client.
	// Note that we'll use logger.Infof here vs the Errorf above.
	logger.Infof("New chat created with ID %v", chatID)
	rsp.ChatId = chatID
	return nil
}
