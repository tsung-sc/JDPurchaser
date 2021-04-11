package controller

import (
	"JD_Purchase/api"
	"github.com/emicklei/go-restful"
	"github.com/jinzhu/gorm"
)

type SystemCtl struct {
	*api.Api
	*gorm.DB
}

func (ctl SystemCtl) WebService(ws *restful.WebService) {

}
