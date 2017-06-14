package models

import (
	"bytes"
	"context"
	"log"

	mgo "gopkg.in/mgo.v2"

	"github.com/satori/go.uuid"
)

var (
	db *mgo.Database
)

type Model struct {
	UUID        *string `bson:"_id" json:"uuid"`
	VersionUUID *string `json:"version_uuid"`
}

// generic validation error
type ValidationError struct {
	msg string
}

func (e *ValidationError) Error() string {
	return e.msg
}

// used when a resource is depended on by one or more other resources
type DependentResourceError struct {
	resources []string // describes resources that depend on the resource in question in the format "<resource type>:<resource uuid>"
}

func (e *DependentResourceError) Error() string {
	var msgb bytes.Buffer
	for i, v := range e.resources {
		msgb.WriteString(v)
		if i < (len(e.resources) - 1) {
			msgb.WriteString(" ")
		}
	}
	return msgb.String()
}

type VersionUUIDMismatchError struct {
	msg string
}

func (e *VersionUUIDMismatchError) Error() string {
	return e.msg
}

// validate uuid format and ensure string representations are lowercased
// only validates present fields
func (m *Model) ValidateUUIDFormat(ctx context.Context) error {
	if m.UUID != nil {
		uid, err := uuid.FromString(*m.UUID)
		if err != nil {
			return err
		}
		*m.UUID = uid.String()
	}
	if m.VersionUUID != nil {
		vuid, err := uuid.FromString(*m.VersionUUID)
		if err != nil {
			return err
		}
		*m.VersionUUID = vuid.String()
	}
	return nil
}

// the models on which this is invoked must have a non-nil uuid and versionuuid
// arg naming: oo = old object, no = new object
// the order of the arguments is irrelevant; both orderings behave identically
func ValidateVersionUUID(ctx context.Context, oo Model, no Model) error {
	if oo.VersionUUID == nil || no.VersionUUID == nil {
		log.Panic("validate version uuid: one or more objects do(es) not have a version uuid!")
	}
	if *oo.VersionUUID != *no.VersionUUID {
		return &VersionUUIDMismatchError{"validate version uuid: version uuid values do not match"}
	}
	return nil
}

func CreateConnection() error {
	conn, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		return err
	}
	db = conn.DB("chatbotmgmt")
	return nil
}
