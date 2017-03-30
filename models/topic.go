package models

import (
	"context"
	"time"

	mgo "gopkg.in/mgo.v2"
)

// represents a fully unmarshalled topic, complete with proper expiration time
type Topic struct {
	Model
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

}
