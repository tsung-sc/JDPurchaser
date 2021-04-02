package controller

import (
	"JD_Purchase/api"
	"github.com/emicklei/go-restful"
	"github.com/jinzhu/gorm"
	"log"
)

type PurchaserCtl struct {
	*gorm.DB
	*api.Api
}

//func (ctl PurchaserCtl) Register(container *restful.Container) {
//	ws := new(restful.WebService)
//	ws.
//		Path("/core").
//		Consumes(restful.MIME_XML, restful.MIME_JSON).
//		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well
//
//	ws.Route(ws.GET("list").To(ctl.list))
//	ws.Route(ws.POST("addTask").To(ctl.addTask))
//	ws.Route(ws.GET("runTask").To(ctl.runTask))
//	//ws.Route(ws.PUT("/{user-id}").To(u.createUser))
//	//ws.Route(ws.DELETE("/{user-id}").To(u.removeUser))
//
//	container.Add(ws)
//}

func (ctl PurchaserCtl) WebService(ws *restful.WebService) {
	ws.Route(ws.GET("list").To(ctl.list))
	ws.Route(ws.POST("addTask").To(ctl.addTask))
	ws.Route(ws.GET("runTask").To(ctl.runTask))
	//ws.Route(ws.PUT("/{user-id}").To(u.createUser))
	//ws.Route(ws.DELETE("/{user-id}").To(u.removeUser))
}

func (ctl PurchaserCtl) list(request *restful.Request, response *restful.Response) {
	log.Printf(request.Request.RemoteAddr)
	response.Write([]byte("Hello World"))
}

func (ctl PurchaserCtl) addTask(request *restful.Request, response *restful.Response) {
	log.Printf(request.Request.RemoteAddr)
	response.Write([]byte("Hello World"))
}

func (ctl PurchaserCtl) runTask(request *restful.Request, response *restful.Response) {
	skuIDs := "730618,4080291:2"
	area := "18_1482_48938_52586"
	result, err := ctl.LoginByQRCode()
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	log.Println(result)
	err = ctl.BuyItemInStock(skuIDs, area, false, 5, 3, 5)
	if err != nil {
		log.Printf("%+v", err)
		return
	}
}
