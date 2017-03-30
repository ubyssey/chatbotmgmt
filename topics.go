package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gocraft/web"

	"github.com/ubyssey/chatbot/models"
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
		fmt.Fprint(rw, "JSON NO GOOD!!!")
		return
	}
	if err := t.Save(ctx); err != nil {
		fmt.Fprint(rw, "Failed to save!!")
		return
	}
	top_url := fmt.Sprintf("/topics/%s", *t.UUID)
	rw.Header().Set("Location", top_url)
}

// GET /topics/:uuid
func (reqctx *ResourceRequestContext) GetTopic(rw web.ResponseWriter, req *web.Request) {
	ctx := context.WithValue(context.Background(), 0, reqctx.rid)
	t := models.Topic{}
	err := t.GetById(ctx)
	if err != nil {
		fmt.Fprintf(rw, "got an error getting the topic!")
		fmt.Println(err)
		return
	}
	j, err := json.Marshal(t)
	if err != nil {
		fmt.Println("ENCODE JSON NOT GOOD:", err)
		return
	}
	fmt.Fprint(rw, string(j))
}

// PATCH /topics/:uuid
func (ctx *ResourceRequestContext) UpdateTopic(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "PATCH /topics/%s", ctx.rid)
}

// DELETE /topics/:uuid
func (ctx *ResourceRequestContext) DeleteTopic(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "DELETE /topics/%s", ctx.rid)
}
