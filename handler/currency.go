package handler

import (
	"hamal/protocol"
	"log"
	"net/http"

	"hamal/db"

	"github.com/gin-gonic/gin"
)

func GetCurrencies(c *gin.Context) {
	var req protocol.GetCurrenciesReq
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	data, err := db.GetCurrencies(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(data))
}

func AddOrUpdateCurrency(c *gin.Context) {
	var req protocol.AddOrUpdateCurrencyReq
	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	req.ID = c.Param("id")
	if req.ID == "" {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	data, err := db.AddOrUpdateCurrency(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(data))
}

func DeleteCurrency(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, ""))
		return
	}

	err := db.DeleteCurrency(id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(nil))
}
