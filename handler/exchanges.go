package handler

import (
	"hamal/protocol"
	"log"
	"net/http"

	"hamal/db"

	"github.com/gin-gonic/gin"
)

func GetExchanges(c *gin.Context) {
	var req protocol.GetExchagesReq
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	exchanges, err := db.GetExchanges()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	var rsp protocol.GetExchangesRsp
	rsp.Count = len(exchanges)

	for _, ex := range exchanges {
		var item protocol.GetExchangesRspItem
		item.Name = ex.Name

		rsp.Items = append(rsp.Items, item)
	}

	c.JSON(http.StatusOK, protocol.OKRsp(rsp))

}
