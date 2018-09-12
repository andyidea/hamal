package handler

import (
	"hamal/protocol"
	"log"

	"fmt"

	"net/http"

	"hamal/db"

	"hamal/models"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req protocol.UserRegisterReq
	err := c.Bind(&req)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	db.UserRegister()

	c.JSON(http.StatusOK, protocol.OKRsp(nil))

}

// Login 用户登录
func Login(c *gin.Context) {
	var req protocol.LoginReq
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	data, err := db.UserLogin(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(data))
	return
}

func GetUserInfo(c *gin.Context) {
	user := getUser(c)

	data, err := db.GetUserInfo(user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(data))
	return
}

func UpdateUserInfo(c *gin.Context) {
	var req protocol.UpdateUserInfoReq
	err := c.Bind(&req)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	user := getUser(c)

	user.Name = req.Name
	user.Avatar = req.Avatar

	err = db.UpdateUserInfo(user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(nil))
	return

}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, protocol.OKRsp(nil))
	return
}

func GetUserExchanges(c *gin.Context) {
	var req protocol.GetUserExchangesReq
	err := c.Bind(&req)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	user := getUser(c)

	ues, err := db.GetUserExchages(user.ID, req.Exchange)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	var rsp protocol.GetUserExchangesRsp
	rsp.Count = len(ues)

	for _, ue := range ues {
		var item protocol.GetUserExchangesRspItem
		item.ID = ue.ID
		item.Exchange = ue.Exchange
		item.Tag = ue.Tag
		item.AccessKey = ue.AccessKey
		//item.SecretKey = ue.SecretKey

		rsp.Items = append(rsp.Items, item)
	}
	c.JSON(http.StatusOK, protocol.OKRsp(rsp))
}

func AddOrUpdateUserExchange(c *gin.Context) {
	var req protocol.AddOrUpdateUserExchangeReq
	err := c.Bind(&req)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	log.Println(req)
	user := getUser(c)

	var exchange models.UserExchange
	exchange.ID = req.ID
	exchange.Tag = req.Tag
	exchange.Exchange = req.Exchange
	exchange.AccessKey = req.AccessKey
	exchange.SecretKey = req.SecretKey
	exchange.UserID = user.ID

	if req.ID == 0 {
		err = db.AddUserExchange(&exchange)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, protocol.ErrorRsp(err))
			return
		}
	}

	c.JSON(http.StatusOK, protocol.OKRsp(nil))
}

func DeleteUserExchange(c *gin.Context) {
	var req protocol.DeleteUserExchangeReq
	err := c.Bind(&req)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, protocol.NewError(protocol.CodeWrongParam, protocol.MsgWrongParam, err.Error()))
		return
	}

	user := getUser(c)

	err = db.DeleteUserExchange(req.ID, user.ID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	c.JSON(http.StatusOK, protocol.OKRsp(nil))
}
