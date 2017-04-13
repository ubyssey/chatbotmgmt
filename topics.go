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
func (reqctx *RequestContext) ListTopics(rw web.ResponseWriter, req *web.Request) {
	ctx := context.Background()
	var topics []models.Topic
	if err := models.GetAllTopics(ctx, &topics); err != nil {
		rw.WriteHeader(500)
		return
	}
	j, err := json.Marshal(map[string]([]models.Topic){
		"results": topics,
	})
	if err != nil {
		log.Print("list topics: failed to encode as json: ", err)
		rw.WriteHeader(500)
		return
	}
	fmt.Fprint(rw, string(j))
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
func (reqctx *ResourceRequestContext) UpdateTopic(rw web.ResponseWriter, req *web.Request) {
	ctx := context.Background()
	decoder := json.NewDecoder(req.Body)
	var t models.Topic
	if err := decoder.Decode(&t); err != nil {
		rw.WriteHeader(400)
		fmt.Fprint(rw, "the request body could not be parsed as json or contained an improperly formatted field")
		return
	}

	// enforce update rules, ensure t.UUID is populated
	if t.UUID != nil {
		if *t.UUID != reqctx.rid { // enforce body/url uuid matching
			rw.WriteHeader(400)
			fmt.Fprint(rw, "the uuid values in the url and request body must match")
			return
		}
	} else {
		t.UUID = new(string)
		*t.UUID = reqctx.rid
	}
	if t.VersionUUID == nil {
		rw.WriteHeader(400)
		fmt.Fprint(rw, "a version uuid must be provided")
		return
	}

	// save the record
	if err := t.Save(ctx); err != nil {
		rw.WriteHeader(500)
		log.Print(err)
		return
	}

	// return the full record as the response body
	if err := t.GetById(ctx, reqctx.rid); err != nil {
		rw.WriteHeader(500)
		log.Print(err)
		return
	}
	j, err := json.Marshal(t)
	if err != nil {
		log.Print("update topic: failed to encode as json: ", err)
		rw.WriteHeader(500)
		return
	}
	fmt.Fprint(rw, string(j))
}

// DELETE /topics/:uuid
func (reqctx *ResourceRequestContext) DeleteTopic(rw web.ResponseWriter, req *web.Request) {
	ctx := context.Background()

	decoder := json.NewDecoder(req.Body)
	var t models.Topic
	if err := decoder.Decode(&t); err != nil {
		rw.WriteHeader(400)
		fmt.Fprint(rw, "the request body could not be parsed as json or contained an improperly formatted field")
		return
	}

	if t.UUID == nil {
		t.UUID = &reqctx.rid
	} else {
		if *t.UUID != reqctx.rid {
			rw.WriteHeader(400)
			fmt.Fprint(rw, "the uuid values in the body and url must match")
			return
		}
	}

	err := t.Delete(ctx)
	if err != nil {
		if err == mgo.ErrNotFound {
			rw.WriteHeader(404)
			return
		}
		switch err.(type) {
		case *models.ValidationError:
			rw.WriteHeader(400)
			fmt.Fprint(rw, err)
		case *models.DependentResourceError:
			rw.WriteHeader(412)
			fmt.Fprint(rw, err) // TODO this error message should be formatted according to spec
		default:
			rw.WriteHeader(500)
		}
		return
	}
	rw.WriteHeader(204)
}
