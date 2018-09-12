package handler

import (
	"hamal/protocol"
	"log"
	"net/http"

	"hamal/strategy"
	"strconv"

	"hamal/db"

	"errors"
	"github.com/gin-gonic/gin"
	"hamal/common"
)

func GetStrategies(c *gin.Context) {
	var req protocol.GetStrategiesReq
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	data, err := db.GetStrategies(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(data))
}

func GetStrategyInstances(c *gin.Context) {
	var req protocol.GetStrategyInstancesReq
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	user := getUser(c)
	data, err := db.GetStrategyInstances(user, req.StrategyID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(data))
}

func AddStrategyInstance(c *gin.Context) {
	var req protocol.AddStrategyInstanceReq
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	user := getUser(c)

	data, err := db.AddStrategyInstance(&req, user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(data))

}

func UpdateStrategyInstance(c *gin.Context) {
	var req protocol.UpdateStrategyInstanceReq
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, ""))
		return
	}

	user := getUser(c)

	id, _ := strconv.Atoi(idStr)

	si, err := db.GetStrategyInstanceModel(uint(id))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	if si.UserID != user.ID {
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeNoPower, protocol.MsgNoPower, ""))
		return
	}

	si.StrategyID = req.StrategyID
	si.BaseCurrencyID = req.BaseCurrencyID
	si.TargetCurrencyID = req.TargetCurrencyID
	si.Amount = req.Amount
	si.Interval = req.Interval

	err = db.UpdateStrategyInstance(si)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(nil))
}

func DeleteStrategyInstance(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, ""))
		return
	}

	id, _ := strconv.Atoi(idStr)

	err := db.DeleteStrategyInstance(uint(id))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(nil))

}

func LaunchStrategyInstance(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, ""))
		return
	}

	user := getUser(c)

	var err error
	id, _ := strconv.ParseUint(idStr, 10, 32)

	sim, err := db.GetStrategyInstanceModel(uint(id))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	if sim.UserID != user.ID {
		c.JSON(http.StatusOK, protocol.ErrorRsp(errors.New("user error")))
		return
	}

	si, ok := strategy.GetStrategyInstance(uint(id))
	if ok {
		si.Launch()
	} else {
		si, err = strategy.NewStrategyInstance(uint(id))
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, protocol.ErrorRsp(err))
			return
		}
		si.Launch()
		strategy.AddStrategyInstance(si)
	}

	sim.Status = common.StrategyInstanceStatusRunning

	err = db.UpdateStrategyInstance(sim)
	if err != nil {
		si.Stop()
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(nil))
}

func StopStrategyInstance(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, ""))
		return
	}
	user := getUser(c)

	id, _ := strconv.ParseUint(idStr, 10, 32)

	sim, err := db.GetStrategyInstanceModel(uint(id))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	if sim.UserID != user.ID {
		c.JSON(http.StatusOK, protocol.ErrorRsp(errors.New("user error")))
		return
	}

	si, ok := strategy.GetStrategyInstance(uint(id))
	if ok {
		si.Stop()
	}

	sim.Status = common.StrategyInstanceStatusRest

	err = db.UpdateStrategyInstance(sim)
	if err != nil {
		si.Launch()
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(nil))

}
