package handler

import (
	"hamal/protocol"
	"log"
	"net/http"

	"hamal/models"

	"hamal/db"

	"github.com/gin-gonic/gin"
)

func GetSymbols(c *gin.Context) {
	data, err := db.GetSymbols()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(data))

}

func AddOrUpdateSymbols(c *gin.Context) {
	var req protocol.AddOrUpdateSymbolsReq
	err := c.Bind(&req)
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

	var symbol models.Symbol
	symbol.ID = req.ID
	symbol.Name = req.Name
	symbol.BaseCurrencyID = req.BaseCurrencyID
	symbol.QuoteCurrencyID = req.QuoteCurrencyID

	err = db.AddOrUpdateSymbol(&symbol)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(nil))
}
