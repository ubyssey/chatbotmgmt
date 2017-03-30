package models

import (
	mgo "gopkg.in/mgo.v2"
)

var (
	Connection mgo.Session
)

type Model struct {
	UUID        *string `json:"uuid"`
	VersionUUID *string `json:"version_uuid"`
}
