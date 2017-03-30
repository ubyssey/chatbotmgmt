package main

import (
	"fmt"
	"github.com/gocraft/web"
)

// GET /campaigns
func (ctx *RequestContext) ListCampaigns(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "GET /campaigns")
}

// POST /campaigns
func (ctx *RequestContext) CreateCampaign(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "POST /campaigns")
}

// GET /campaigns/:uuid
func (ctx *ResourceRequestContext) GetCampaign(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "GET /campaigns/%s", ctx.rid)
}

// PATCH /campaigns/:uuid
func (ctx *ResourceRequestContext) UpdateCampaign(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "PATCH /campaigns/%s", ctx.rid)
}

// DELETE /campaigns/:uuid
func (ctx *ResourceRequestContext) DeleteCampaign(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "DELETE /campaigns/%s", ctx.rid)
}
