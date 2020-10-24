package repository

import (
	"errors"
	"time"

	"github.com/micro-community/micro-chat/model"
	message "github.com/micro-community/micro-chat/proto"
	"github.com/micro/micro/v3/service/store"
)

//Repository for message
type Repository struct {
	messsages model.Table
}

//New return a message repo
func New() *Repository {
	nameIndex := model.ByEquality("name")
	nameIndex.Unique = true
	emailIndex := model.ByEquality("email")
	emailIndex.Unique = true

	return &Repository{
		messsages: model.NewTable(store.DefaultStore, "messsages", model.Indexes(nameIndex, emailIndex), nil),
	}
}

//Create for
func (repo *Repository) Create(msg *message.Message, salt string) error {
	msg.SentAt = time.Now().Unix()
	err := repo.messsages.Save(msg)
	if err != nil {
		return err
	}
	return repo.messsages.Save(message.Message{
		Id: msg.Id,
	})
}

//Delete messages
func (repo *Repository) Delete(id string) error {
	return repo.messsages.Delete(model.Equals("id", id))
}

//Update messages
func (repo *Repository) Update(msg *message.Message) error {
	msg.SentAt = time.Now().Unix()
	return repo.messsages.Save(msg)
}

//Read messages
func (repo *Repository) Read(id string) (*message.Message, error) {
	messsage := &message.Message{}
	return messsage, repo.messsages.Read(model.Equals("id", id), messsage)
}

//Search messages
func (repo *Repository) Search(username, email string, limit, offset int64) ([]*message.Message, error) {
	var query model.Query
	if len(username) > 0 {
		query = model.Equals("name", username)
	} else if len(email) > 0 {
		query = model.Equals("email", email)
	} else {
		return nil, errors.New("username and email cannot be blank")
	}

	messsages := []*message.Message{}
	return messsages, repo.messsages.List(query, &messsages)
}
