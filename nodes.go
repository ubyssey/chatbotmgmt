package main

import (
	"fmt"
	"github.com/gocraft/web"
)

// POST /campaigns/:uuid/nodes
func (ctx *ResourceRequestContext) CreateNode(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "POST /campaigns/%s/nodes", ctx.rid)
}

// PATCH /campaigns/:uuid/nodes/:node_uuid
func (ctx *NodeRequestContext) UpdateNode(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "PATCH /campaigns/%s/nodes/%s", ctx.rid, ctx.nid)
}

// DELETE /campaigns/:uuid/nodes/:node_uuid
func (ctx *NodeRequestContext) DeleteNode(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "DELETE /campaigns/%s/nodes/%s", ctx.rid, ctx.nid)
}
