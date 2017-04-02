package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	topicCollection = "topics"
)

// represents a fully unmarshalled topic, complete with proper expiration time
type Topic struct {
	Model       `bson:",inline"`
	Transient   *bool      `json:"transient"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Parent      *string    `json:"parent"`
	ExpiresAt   *time.Time `json:"expires_at"` // this Just Works(TM) because `time` implements the json.Unmarshaler interface :-)
}

// Validate the format of the fields of the passed topic.
func (t *Topic) ValidateFormat() error {
	if t.Title == nil || *t.Title == "" {
		return &ValidationError{"validate topic format: title is required"}
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

// validate that the topic can be deleted (that no other resources depend on it)
func (t *Topic) ValidateDelete(ctx context.Context) error {
	var results []Topic
	err := db.C(topicCollection).Find(bson.M{"parent": *t.UUID}).Select(bson.M{"_id": 1}).All(&results)
	if err != nil {
		return err
	}
	if len(results) > 0 {
		var err DependentResourceError
		for _, v := range results {
			err.resources = append(err.resources, fmt.Sprintf("topic:%s", *v.UUID))
		}
		return &err
	}
	return nil
}

func (t *Topic) GetById(ctx context.Context, tid string) error {
	return db.C(topicCollection).FindId(tid).One(t)
}

// Delete a Topic by ID
func (t *Topic) DeleteById(ctx context.Context, tid string) error {
	t.UUID = &tid // necessary for ValidateDelete()
	if err := t.ValidateDelete(ctx); err != nil {
		return err
	}
	if err := db.C(topicCollection).RemoveId(tid); err != nil {
		if err != mgo.ErrNotFound {
			log.Print("delete topic: db error: ", err)
		}
		return err
	}
	return nil
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
