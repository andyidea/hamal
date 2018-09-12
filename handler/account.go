package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leek-box/GoEx/builder"
	"hamal/db"
	"hamal/protocol"
	"log"
	"net/http"
)

func GetAccountBalance(c *gin.Context) {
	user := getUser(c)
	userExchanges, err := db.GetUserExchages(user.ID, "")
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	var rsp protocol.GetAccountBalanceRsp

	for _, ue := range userExchanges {
		builder := builder.NewAPIBuilder()
		builder.APIKey(ue.AccessKey)
		log.Println(ue.AccessKey, ue.SecretKey)
		builder.APISecretkey(ue.SecretKey)
		api := builder.Build(ue.Exchange)
		account, err := api.GetAccount()
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, protocol.ErrorRsp(err))
			return
		}

		for _, v := range account.SubAccounts {
			if v.Amount > 0 {
				var item protocol.GetAccountBalanceRspItem
				item.Currency = v.Currency.String()
				item.Balance = v.Amount
				item.FrozenBalance = v.ForzenAmount
				item.TradeBalance = v.Amount - v.ForzenAmount

				rsp.Items = append(rsp.Items, item)
			}

		}
	}

	c.JSON(http.StatusOK, protocol.OKRsp(rsp))
}
