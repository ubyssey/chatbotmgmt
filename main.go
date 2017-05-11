package main

import (
	"fmt"
	"github.com/gocraft/web"
	"github.com/satori/go.uuid"
	"github.com/ubyssey/chatbotmgmt/models"
	"log"
	"net/http"
)

type RequestContext struct{}

// Context for a request referencing a single resource identified by a string ID
type ResourceRequestContext struct {
	*RequestContext
	rid string // resource ID
}

type NodeRequestContext struct {
	*ResourceRequestContext
	nid string // node ID
}

func ReadUuidParam(req *web.Request, panme string) (string, error) {
	return "", nil // TODO implement
}

func (ctx *ResourceRequestContext) LoadUuid(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc, param string, dest *string) {
	_, err := uuid.FromString(req.PathParams[param])
	if err != nil {
		rw.WriteHeader(400)
		fmt.Fprintf(rw, "%s doesn't look like a UUID to me!", req.PathParams[param])
		return
	}
	*dest = req.PathParams[param]
	next(rw, req)
}

func (ctx *ResourceRequestContext) LoadResourceId(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	ctx.LoadUuid(rw, req, next, "uuid", &ctx.rid)
}

func (ctx *NodeRequestContext) LoadNodeId(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	ctx.LoadUuid(rw, req, next, "nodeuuid", &ctx.nid)
}

func (ctx *RequestContext) AuthenticateRequest(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	// TODO real authentication
	next(rw, req)
}

// I want to keep this one around because it's cute
func (ctx *RequestContext) SendHelloWorld(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "Hello World! I am Gopher")
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC) // use UTC for logging

	err := models.CreateConnection()
	if err != nil {
		log.Print("[FATAL] Could not connect to mongo!")
		log.Fatal(err)
	}

	router := web.New(RequestContext{}).
		Middleware((*RequestContext).AuthenticateRequest).
		Get("/", (*RequestContext).SendHelloWorld)

	// topic routes
	router.Get("/topics", (*RequestContext).ListTopics)
	router.Post("/topics", (*RequestContext).CreateTopic)

	topicRouter := router.Subrouter(ResourceRequestContext{}, "/topics")
	topicRouter.Middleware((*ResourceRequestContext).LoadResourceId)
	topicRouter.Get("/:uuid", (*ResourceRequestContext).GetTopic)
	topicRouter.Patch("/:uuid", (*ResourceRequestContext).UpdateTopic)
	topicRouter.Delete("/:uuid", (*ResourceRequestContext).DeleteTopic)

	// campaign routes
	router.Get("/campaigns", (*RequestContext).ListCampaigns)
	router.Post("/campaigns", (*RequestContext).CreateCampaign)

	campaignRouter := router.Subrouter(ResourceRequestContext{}, "/campaigns")
	campaignRouter.Middleware((*ResourceRequestContext).LoadResourceId)
	campaignRouter.Get("/:uuid", (*ResourceRequestContext).GetCampaign)
	campaignRouter.Patch("/:uuid", (*ResourceRequestContext).UpdateCampaign)
	campaignRouter.Delete("/:uuid", (*ResourceRequestContext).DeleteCampaign)
	campaignRouter.Post("/:uuid/nodes", (*ResourceRequestContext).CreateNode)

	nodeRouter := campaignRouter.Subrouter(NodeRequestContext{}, "/:uuid/nodes")
	nodeRouter.Middleware((*NodeRequestContext).LoadNodeId)
	nodeRouter.Patch("/:nodeuuid", (*NodeRequestContext).UpdateNode)
	nodeRouter.Delete("/:nodeuuid", (*NodeRequestContext).DeleteNode)

	log.Print("Binding server to localhost:3000")
	http.ListenAndServe("localhost:3000", router)
}
