package models

import (
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

const (
	campaignCollection = "campaigns"
)

type Campaign struct {
	Model     `bson:",inline"`
	ArticleId *string          `json:"article_id"`
	PublishAt *time.Time       `json:"publish_at"`
	Topics    *[]string        `json:"topics"`
	RootNode  *string          `json:"root_node"`
	Nodes     *map[string]Node `json:"nodes"`
	Name      *string          `json:"name"`
}

type Node struct {
	Model   `bson:",inline"`
	Effect  *string                 `json:"effect"`
	Content *map[string]interface{} `json:"content"`
	Actions *[]UserAction           `json:"user_actions"`
}

type UserAction struct {
	Model  `bson:",inline"`
	Type   *string `json:"type"`
	Label  *string `json:"label"`
	Target *string `json:"target"`
}

func (c *Campaign) Validate(ctx context.Context) error {
	if c.Nodes == nil {
		return &ValidationError{"validate campaign: \"nodes\" is required"}
	}
	if c.PublishAt == nil {
		return &ValidationError{"validate campaign: a publication date is required"}
	}
	if c.Topics == nil || len(*c.Topics) < 1 {
		return &ValidationError{"validate campaign: a topics array is required"}
	}
	if c.RootNode == nil {
		return &ValidationError{"validate campaign: a root node is required"}
	}
	if c.Name == nil || *c.Name == "" {
		return &ValidationError{"validate campaign: a name is required"}
	}
	if _, ok := (*c.Nodes)[*c.RootNode]; !ok {
		return &ValidationError{"validate campaign: the root node must refer to an existant node"}
	}
	// TODO validate each node - eventually
	return nil
}

// validate that resources referenced (read: topics) exist and make sense
func (c *Campaign) ValidateReferences(ctx context.Context) error {
	if c.Topics == nil || len(*c.Topics) < 1 {
		return &ValidationError{"validate campaign: \"nodes\" is required"}
	}
	for _, t := range *c.Topics {
		n, err := db.C(topicCollection).FindId(t).Count()
		if err != nil {
			return err
		}
		if n != 1 {
			return &ValidationError{fmt.Sprintf("no topic with the uuid %s exists", t)}
		}
	}
	return nil
}

func (c *Campaign) GetById(ctx context.Context, cid string) error {
	return db.C(campaignCollection).FindId(cid).One(c)
}

func (c *Campaign) Save(ctx context.Context) error {
	newrec := c.UUID == nil // new record if uuid is nil
	if newrec {
		c.UUID = new(string)
		*c.UUID = uuid.NewV4().String()
	} else {
		log.Panic("not implemented")
	}
	if c.VersionUUID == nil {
		c.VersionUUID = new(string)
	}
	*c.VersionUUID = uuid.NewV4().String()
	if err := c.Validate(ctx); err != nil {
		return err
	}
	if err := c.ValidateReferences(ctx); err != nil {
		return err
	}
	if newrec {
		return db.C(campaignCollection).Insert(c)
	} else {
		log.Panic("not implemented")
	}
	return nil
}
