package models

import (
	mgo "gopkg.in/mgo.v2"
)

var (
	db *mgo.Database
)

type Model struct {
	UUID        *string `bson:"_id" json:"uuid"`
	VersionUUID *string `json:"version_uuid"`
}

func CreateConnection() error {
	conn, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		return err
	}
	db = conn.DB("chatbotmgmt")
	return nil
}
