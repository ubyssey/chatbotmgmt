package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocraft/web"

	"github.com/ubyssey/chatbot/models"
)

// GET /topics
func (ctx *RequestContext) ListTopics(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "GET /topics")
}

// POST /topics
func (ctx *RequestContext) CreateTopic(rw web.ResponseWriter, req *web.Request) {
	decoder := json.NewDecoder(req.Body)
	var t models.Topic
	if err := decoder.Decode(&t); err != nil {
		fmt.Fprint(rw, "JSON NO GOOD!!!")
		return
	}
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println("ENCODE JSON NOT GOOD:", err)
	}
	os.Stdout.Write(b)
	fmt.Fprint(rw, "GOT SOME YUMMY JSON!!")
}

// GET /topics/:uuid
func (ctx *ResourceRequestContext) GetTopic(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "GET /topics/%s", ctx.rid)
}

// PATCH /topics/:uuid
func (ctx *ResourceRequestContext) UpdateTopic(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "PATCH /topics/%s", ctx.rid)
}

// DELETE /topics/:uuid
func (ctx *ResourceRequestContext) DeleteTopic(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "DELETE /topics/%s", ctx.rid)
}
