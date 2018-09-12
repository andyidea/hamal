package middleware

import (
	"encoding/json"
	"hamal/util"
	"log"

	"hamal/protocol"
	"net/http"

	"time"

	"hamal/models"

	"hamal/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var writeUrls = []string{
	"/v1/user/login",
}

func ClaimsToStruct(as map[string]interface{}) models.User {
	result := models.User{}

	data := map[string]interface{}{}
	for k, v := range as {
		switch v.(type) {
		case int:

			data[k] = v
		case float64:

			data[k] = int64(v.(float64))
		default:
			data[k] = v
		}
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("ClaimsToStruct", err)
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Println("ClaimsToStruct", string(bytes), err)
	}
	return result
}

func AuthVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, url := range writeUrls {
			if url == c.Request.URL.Path {
				c.Next()
				return
			}
		}
		token := c.GetHeader("X-Token")
		t, err := util.ParseToken(token)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusOK, protocol.NewError(protocol.CodeTokenCheckFailed, protocol.MsgTokenCheckFailed, err.Error()))
			return
		}

		claims := t.Claims.(jwt.MapClaims)

		creadAt := int64(claims["token_created_at"].(float64))
		if creadAt+60*60*24*30 < time.Now().Unix() {
			c.AbortWithStatusJSON(http.StatusOK, protocol.NewError(protocol.CodeTokenOutDate, protocol.MsgTokenOutDate, ""))
			return
		}

		uid := uint(claims["ID"].(float64))

		user, err := db.GetUserInfoModel(uid)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusOK, protocol.NewError(protocol.CodeTokenCheckFailed, protocol.MsgTokenCheckFailed, err.Error()))
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
