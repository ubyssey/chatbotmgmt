package models

import (
	"context"
	"fmt"
	"net/url"

	mgo "gopkg.in/mgo.v2"
)

type Node struct {
	Model        `bson:",inline"`
	Effect       *string                 `json:"effect"`
	Content      *map[string]interface{} `json:"content"`
	Actions      *[]UserAction           `json:"user_actions"`
	TopicUuid    *string                 `json:"topic"`
	Confirmation *string                 `json:"confirmation"`
	NextNode     *string                 `json:"next"`
}

type UserAction struct {
	Model  `bson:",inline"`
	Type   *string `json:"type"`
	Label  *string `json:"label"`
	Target *string `json:"target"`
}

// validate the node in the context of the campaign to which it belongs.
func (n *Node) Validate(ctx context.Context, c *Campaign) error {
	if n.Effect == nil {
		return &ValidationError{"validate campaign node: an effect is required"}
	}
	switch *n.Effect {
	case "message":
		if n.Content == nil {
			return &ValidationError{"validate campaign node: a message node must include content"}
		}
		if n.Actions != nil && len(*n.Actions) > 0 {
			for _, a := range *n.Actions {
				if err := a.Validate(ctx, c); err != nil {
					return err
				}
			}
		}
	case "subscribe_topic":
	case "unsubscribe_topic":
		if n.NextNode != nil {
			if _, ok := (*c.Nodes)[*n.NextNode]; !ok {
				return &ValidationError{fmt.Sprintf("validate campaign node: no node with the uuid %s exists in this campaign", *n.NextNode)}
			}
		}
		if n.TopicUuid == nil {
			return &ValidationError{"validate campaign node: a topic subscription node must have a topic UUID"}
		}
		tt := new(Topic)
		if err := tt.GetById(ctx, *n.TopicUuid); err != nil {
			if err == mgo.ErrNotFound {
				return &ValidationError{fmt.Sprintf("validate campaign node: no topic with the uuid %s exists", *n.TopicUuid)}
			}
			return err
		}
	default:
		return &ValidationError{"validate campaign node: a node's effect must be one of: (\"message\", \"subscribe_topic\", \"unsubscribe_topic\")"}
	}
	return nil
}

// validate the action in the context of the campaign to which it belongs.
func (a *UserAction) Validate(ctx context.Context, c *Campaign) error {
	if a.Type == nil {
		return &ValidationError{"validate campaign node action: a type is required"}
	}
	if a.Label == nil {
		return &ValidationError{"validate campaign node action: a label is required"}
	}
	if a.Target == nil {
		return &ValidationError{"validate campaign node action: a target is required"}
	}
	switch *a.Type {
	case "node":
		if _, ok := (*c.Nodes)[*a.Target]; !ok {
			return &ValidationError{fmt.Sprintf("validate campaign node action: no node with the uuid %s exists in this campaign", *a.Target)}
		}
	case "campaign":
		cpg := new(Campaign)
		if err := cpg.GetById(ctx, *a.Target); err != nil {
			if err == mgo.ErrNotFound {
				return &ValidationError{fmt.Sprintf("validate campaign node action: no campaign with the uuid %s exists", *a.Target)}
			}
			return err
		}
	case "link":
		if _, err := url.Parse(*a.Target); err != nil {
			return &ValidationError{fmt.Sprintf("validate campaign node action: the url \"%s\" appears to be malformed", *a.Target)}
		}
	default:
		return &ValidationError{"validate campaign node action: a user action's type must be one of: (\"node\", \"campaign\", \"link\")"}
	}
	return nil
}
