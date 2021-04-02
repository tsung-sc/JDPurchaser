package controller

import (
	"JD_Purchase/api"
	"JD_Purchase/db"
	"github.com/emicklei/go-restful"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type instance struct {
	Api *api.Api
	DB  *gorm.DB
}

func InitInstance() *instance {
	var err error
	client := &http.Client{}
	instance := new(instance)
	instance.Api, err = api.NewApi(client)
	if err != nil {
		log.Fatalln(err)
	}
	instance.DB = db.Instance()
	return instance
}

type NS struct {
	Controller CtlInterface
	Path       string
}

type CtlInterface interface {
	WebService(ws *restful.WebService)
}

func RegisterController() []*NS {
	instance := InitInstance()
	routers := []*NS{
		{
			Controller: PurchaserCtl{instance.DB, instance.Api},
			Path:       "/core",
		},
	}
	return routers
}

func InitController(wsContainer *restful.Container) {
	for _, controller := range RegisterController() {
		ws := new(restful.WebService)
		ws.
			Path(controller.Path).
			Consumes(restful.MIME_XML, restful.MIME_JSON).
			Produces(restful.MIME_JSON, restful.MIME_XML)
		controller.Controller.WebService(ws)
		wsContainer.Add(ws)
	}
}
