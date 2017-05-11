package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"log"

	"github.com/ubyssey/chatbotmgmt/models"

	mgo "gopkg.in/mgo.v2"
)

// GET /campaigns
func (reqctx *RequestContext) ListCampaigns(rw web.ResponseWriter, req *web.Request) {
	ctx := context.Background()
	campaigns := make([]models.Campaign, 0)
	if err := models.GetAllCampaigns(ctx, &campaigns); err != nil {
		rw.WriteHeader(500)
		return
	}
	// nil out values we don't care about
	for i, _ := range campaigns {
		campaigns[i].Nodes = nil
		campaigns[i].RootNode = nil
	}
	j, err := json.Marshal(map[string](interface{}){
		"results": campaigns,
	})
	if err != nil {
		log.Print("list campaigns: failed to encode as json: ", err)
		rw.WriteHeader(500)
		return
	}
	fmt.Fprint(rw, string(j))
}

// POST /campaigns
func (reqctx *RequestContext) CreateCampaign(rw web.ResponseWriter, req *web.Request) {
	ctx := context.Background()
	decoder := json.NewDecoder(req.Body)
	var c models.Campaign
	if err := decoder.Decode(&c); err != nil {
		rw.WriteHeader(400)
		fmt.Fprint(rw, "the request body could not be parsed as json or contained an improperly formatted field")
		log.Print("create campaign: failed to parse body: ", err)
		return
	}
	if err := c.Save(ctx); err != nil {
		switch err.(type) {
		case *models.ValidationError:
			rw.WriteHeader(400)
			fmt.Fprint(rw, err)
		default:
			rw.WriteHeader(500)
			log.Print("create campaign: failed to save: ", err)
		}
		return
	}
	cpg_url := fmt.Sprintf("/campaigns/%s", *c.UUID)
	rw.Header().Set("Location", cpg_url)
	rw.WriteHeader(201)
}

// GET /campaigns/:uuid
func (reqctx *ResourceRequestContext) GetCampaign(rw web.ResponseWriter, req *web.Request) {
	ctx := context.Background()
	var c models.Campaign
	err := c.GetById(ctx, reqctx.rid)
	if err != nil {
		if err == mgo.ErrNotFound {
			rw.WriteHeader(404)
		} else {
			rw.WriteHeader(500)
			log.Print(err)
		}
		return
	}
	j, err := json.Marshal(c)
	if err != nil {
		log.Print("get campaign: failed to encode as json: ", err)
		rw.WriteHeader(500)
		return
	}
	fmt.Fprint(rw, string(j))
}

// PATCH /campaigns/:uuid
func (ctx *ResourceRequestContext) UpdateCampaign(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "PATCH /campaigns/%s", ctx.rid)
}

// DELETE /campaigns/:uuid
func (ctx *ResourceRequestContext) DeleteCampaign(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "DELETE /campaigns/%s", ctx.rid)
}
