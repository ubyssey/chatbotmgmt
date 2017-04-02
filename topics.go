package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocraft/web"

	"github.com/ubyssey/chatbot/models"
	mgo "gopkg.in/mgo.v2"
)

// GET /topics
func (ctx *RequestContext) ListTopics(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "GET /topics")
}

// POST /topics
func (reqctx *RequestContext) CreateTopic(rw web.ResponseWriter, req *web.Request) {
	ctx := context.Background()

	decoder := json.NewDecoder(req.Body)
	var t models.Topic
	if err := decoder.Decode(&t); err != nil {
		rw.WriteHeader(400)
		fmt.Fprint(rw, "the request body could not be parsed as json or contained an improperly formatted field")
		return
	}
	if err := t.Save(ctx); err != nil {
		switch err.(type) {
		case *models.ValidationError:
			rw.WriteHeader(400)
			fmt.Fprint(rw, err)
		default:
			rw.WriteHeader(500)
		}
		return
	}
	top_url := fmt.Sprintf("/topics/%s", *t.UUID)
	rw.Header().Set("Location", top_url)
	rw.WriteHeader(201)
}

// GET /topics/:uuid
func (reqctx *ResourceRequestContext) GetTopic(rw web.ResponseWriter, req *web.Request) {
	ctx := context.Background()
	var t models.Topic
	err := t.GetById(ctx, reqctx.rid)
	if err != nil {
		if err == mgo.ErrNotFound {
			rw.WriteHeader(404)
		} else {
			rw.WriteHeader(500)
			log.Print(err)
		}
		return
	}
	j, err := json.Marshal(t)
	if err != nil {
		log.Print("get topic: failed to encode as json: ", err)
		rw.WriteHeader(500)
		return
	}
	fmt.Fprint(rw, string(j))
}

// PATCH /topics/:uuid
func (ctx *ResourceRequestContext) UpdateTopic(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "PATCH /topics/%s", ctx.rid)
}

// DELETE /topics/:uuid
func (reqctx *ResourceRequestContext) DeleteTopic(rw web.ResponseWriter, req *web.Request) {
	ctx := context.Background()
	err := (&models.Topic{}).DeleteById(ctx, reqctx.rid)
	if err != nil {
		switch err {
		case mgo.ErrNotFound:
			rw.WriteHeader(404)
		default:
			rw.WriteHeader(500)
		}
		return
	}
	rw.WriteHeader(204)
}
