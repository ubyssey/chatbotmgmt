package models

import (
	"context"
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

// validate topic. does not care about version uuids
func (t *Topic) Validate(ctx context.Context) error {
	if err := t.Model.ValidateUUIDFormat(ctx); err != nil {
		return err
	}
	if t.Title == nil || *t.Title == "" {
		return &ValidationError{"validate topic format: title is required"}
	}
	if t.Transient == nil {
		return &ValidationError{"validate topic format: transient is required"}
	}
	if *t.Transient == false && t.ExpiresAt != nil {
		return &ValidationError{"validate topic format: topics with an expiration date must be transient"}
	}
	if err := t.ValidateReferences(ctx); err != nil {
		return err
	}
	return nil
}

func (t *Topic) NormalizeUUIDFormat(ctx context.Context) error {
	if t.Parent != nil {
		pid, err := uuid.FromString(*t.Parent)
		if err != nil {
			return err
		}
		*t.Parent = pid.String()
	}
	return nil
}

// validate references are valid
func (t *Topic) ValidateReferences(ctx context.Context) error {
	if t.Parent != nil {
		pt := new(Topic)
		if err := pt.GetById(ctx, *t.Parent); err != nil {
			if err == mgo.ErrNotFound {
				return &ValidationError{"validate topic: the parent topic could not be found"}
			}
			return err
		}
	}
	return nil
}

// validate that the topic can be deleted (that no other resources depend on it)
func (t *Topic) ValidateDelete(ctx context.Context) error {
	if t.UUID == nil {
		return &ValidationError{"validate delete topic: no uuid provided"}
	}
	if t.VersionUUID == nil {
		return &ValidationError{"validate delete topic: no version uuid provided"}
	}
	ot := new(Topic)
	if err := db.C(topicCollection).FindId(*t.UUID).One(ot); err != nil {
		return err
	}
	if err := ValidateVersionUUID(ctx, t.Model, ot.Model); err != nil {
		return err
	}

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
	// TODO check for campaigns referencing the topic
	// TODO check for sub mgmt nodes referencing the topic
	return nil
}

func GetAllTopics(ctx context.Context, t *[]Topic) error {
	return db.C(topicCollection).Find(nil).All(t)
}

func (t *Topic) GetById(ctx context.Context, tid string) error {
	return db.C(topicCollection).FindId(tid).One(t)
}

// Delete a Topic by ID
func (t *Topic) Delete(ctx context.Context) error {
	if err := t.ValidateDelete(ctx); err != nil {
		return err
	}
	if err := db.C(topicCollection).RemoveId(*t.UUID); err != nil {
		if err != mgo.ErrNotFound {
			log.Print("delete topic: db error: ", err)
		}
		return err
	}
	return nil
}

func (t *Topic) Save(ctx context.Context) error {
	newrec := t.UUID == nil // it's a new record if its uuid is nil
	if !newrec {
		ot := new(Topic)
		if err := db.C(topicCollection).FindId(*t.UUID).One(ot); err != nil {
			return err
		}
		if err := ValidateVersionUUID(ctx, ot.Model, t.Model); err != nil {
			return err
		}
		// fuse existant and fresh fields
		// this is gross and should probably use reflection so it's not tightly coupled, but that's a job for later
		if t.Transient != nil {
			ot.Transient = t.Transient
		}
		if t.Title != nil {
			ot.Title = t.Title
		}
		if t.Description != nil {
			ot.Description = t.Description
		}
		if t.Parent != nil {
			ot.Parent = t.Parent
		}
		if t.ExpiresAt != nil {
			ot.ExpiresAt = t.ExpiresAt
		}
		t = ot // override the topic we're working on with our fresh one
	} else {
		t.UUID = new(string)
		*t.UUID = uuid.NewV4().String()
	}
	if t.VersionUUID == nil {
		t.VersionUUID = new(string)
	}
	*t.VersionUUID = uuid.NewV4().String()
	t.NormalizeUUIDFormat(ctx)
	if err := t.Validate(ctx); err != nil {
		return err
	}
	if err := t.ValidateReferences(ctx); err != nil {
		return err
	}
	if newrec {
		return db.C(topicCollection).Insert(t)
	} else {
		return db.C(topicCollection).Update(bson.M{"_id": *t.UUID}, t) // use record ID as primary key
	}
}
