package models

import (
	"context"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

const (
	topicCollection = "topics"
)

// represents a fully unmarshalled topic, complete with proper expiration time
type Topic struct {
	UUID        *string    `bson:"_id" json:"uuid"`
	VersionUUID *string    `json:"version_uuid"`
	Transient   *bool      `json:"transient"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Parent      *string    `json:"parent"`
	ExpiresAt   *time.Time `json:"expires_at"` // this Just Works(TM) because `time` implements the json.Unmarshaler interface :-)
}

// Validate the format of the fields of the passed topic.
func (t *Topic) ValidateFormat() error {
	if t.Title == nil || *t.Title == "" {
		return nil // TODO return a non-nil value
	}
	return nil
}

func (t *Topic) Validate() error {
	if err := t.ValidateFormat(); err != nil {
		return err
	}
	return nil // TODO: if the topic has a parent, make sure the parent exists!
}

func (t *Topic) GetById(ctx context.Context) error {
	if tid, ok := ctx.Value(0).(string); ok {
		log.Printf("using topic id %s\n", string(tid))
	}
	return db.C(topicCollection).FindId(ctx.Value(0)).One(t)
}

func (t *Topic) Save(ctx context.Context) error {
	if t.UUID == nil {
		t.UUID = new(string)
		*t.UUID = uuid.NewV4().String()
	}
	if t.VersionUUID == nil {
		t.VersionUUID = new(string)
	}
	*t.VersionUUID = uuid.NewV4().String()
	if err := t.Validate(); err != nil {
		return err
	}
	return db.C(topicCollection).Insert(t)
}
