package models

import (
	"bytes"

	mgo "gopkg.in/mgo.v2"
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

func CreateConnection() error {
	conn, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		return err
	}
	db = conn.DB("chatbotmgmt")
	return nil
}
