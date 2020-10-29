package model

import (
	"errors"
	"time"

	message "github.com/micro-community/micro-chat/proto"
	"github.com/micro/dev/model"
	"github.com/micro/micro/v3/service/store"
)

//Repository for message
type Repository struct {
	Name      string
	messsages model.Table
}

//NewRepository return a message repo
func NewRepository(repoName string) *Repository {

	clientIndex := model.ByEquality("ClientId")
	clientIndex.Unique = true
	//	userIndex := model.ByEquality("UserId")
	//	userIndex.Unique = true

	return &Repository{
		Name:      repoName,
		messsages: model.NewTable(store.DefaultStore, repoName, model.Indexes(clientIndex), nil),
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