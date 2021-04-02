package server

import (
	"JD_Purchase/config"
	"JD_Purchase/server/controller"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	_ "net/http/pprof"
)

var listenPort string

func ListenAndServe() {
	wsContailer := restful.NewContainer()
	wsContailer.Router(restful.CurlyRouter{})
	controller.InitController(wsContailer)
	if config.Get().ListenPort != "" {
		listenPort = config.Get().ListenPort
	} else {
		listenPort = ":9527"
	}
	log.Printf("start listening on localhost%s", listenPort)
	server := &http.Server{Addr: listenPort, Handler: wsContailer}
	log.Fatal(server.ListenAndServe())
}
