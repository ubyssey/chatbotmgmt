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

type ValidationError struct {
	msg string
}

func (e *ValidationError) Error() string {
	return e.msg
}

type DependentResourceError struct {
	resources []string
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

func CreateConnection() error {
	conn, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		return err
	}
	db = conn.DB("chatbotmgmt")
	return nil
}
