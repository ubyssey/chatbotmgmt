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
	RootNode  *string          `json:"root_node,omitempty"`
	Nodes     *map[string]Node `json:"nodes,omitempty"`
	Name      *string          `json:"name"`
}

func GetAllCampaigns(ctx context.Context, c *[]Campaign) error {
	return db.C(campaignCollection).Find(nil).All(c)
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
	for _, node := range *c.Nodes {
		if err := node.Validate(ctx, c); err != nil {
			return err
		}
	}
	return nil
}

// validate the format of, and normalize the format of, all uuids in the campaign
func (c *Campaign) NormalizeUUIDFormat(ctx context.Context) error {
	if c.Topics != nil {
		for i, tid := range *c.Topics {
			ntid, err := uuid.FromString(tid)
			if err != nil {
				return err
			}
			(*c.Topics)[i] = ntid.String()
		}
	}
	newNodes := make(map[string]Node)
	for k, v := range *c.Nodes {
		nid, err := uuid.FromString(k)
		if err != nil {
			return err
		}
		switch *v.Effect {
		case "message":
			if v.Actions != nil {
				for _, v := range *v.Actions {
					if *v.Type == "node" || *v.Type == "campaign" {
						tid, err := uuid.FromString(*v.Target)
						if err != nil {
							return err
						}
						*v.Target = tid.String()
					}
				}
			}
		case "subscribe_topic":
		case "unsubscribe_topic":
			tid, err := uuid.FromString(*v.TopicUuid)
			if err != nil {
				return err
			}
			*v.TopicUuid = tid.String()
		}
		newNodes[nid.String()] = v
	}
	c.Nodes = &newNodes
	if c.RootNode != nil {
		rnid, err := uuid.FromString(*c.RootNode)
		if err != nil {
			return err
		}
		*c.RootNode = rnid.String()
	}
	return nil
}

// validate that resources referenced (read: topics) exist and make sense
func (c *Campaign) ValidateReferences(ctx context.Context) error {
	// TODO validate that this campaign doesn't refer to itself?
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
	c.NormalizeUUIDFormat(ctx)
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
