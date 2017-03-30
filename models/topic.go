package models

import (
	"context"
	"errors"
	"github.com/satori/go.uuid"
	"time"

	mgo "gopkg.in/mgo.v2"
)

const (
	topicCollection = "topics"
)

// represents a fully unmarshalled topic, complete with proper expiration time
type Topic struct {
	persisted   bool
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
		return errors.New("validate topic format: title is required")
	}
	if t.Transient == nil {
		return errors.New("validate topic format: transient is required")
	}
	if *t.Transient == false && t.ExpiresAt != nil {
		return errors.New("validate topic format: topics with an expiratin date must be transient")
	}
	return nil
}

func (t *Topic) Validate(ctx context.Context) error {
	if err := t.ValidateFormat(); err != nil {
		return err
	}
	if t.Parent != nil {
		pt := new(Topic)
		if err := pt.GetById(ctx, *t.Parent); err != nil {
			if err == mgo.ErrNotFound {
				return errors.New("validate topic parent: the parent topic could not be found")
			}
			return err
		}
	}
	return nil
}

func (t *Topic) GetById(ctx context.Context, tid string) error {
	t.persisted = true
	return db.C(topicCollection).FindId(tid).One(t)
}

func (t *Topic) Save(ctx context.Context) error {
	if t.UUID != nil && t.VersionUUID != nil {
		st := new(Topic)
		st.GetById(ctx, *t.UUID)
		if *t.VersionUUID != *st.VersionUUID {
			return errors.New("Version UUID mismatch!")
		}
	}
	if t.UUID == nil {
		t.UUID = new(string)
		*t.UUID = uuid.NewV4().String()
	}
	if t.VersionUUID == nil {
		t.VersionUUID = new(string)
	}
	*t.VersionUUID = uuid.NewV4().String()
	if err := t.Validate(ctx); err != nil {
		return err
	}
	return db.C(topicCollection).Insert(t)
}
