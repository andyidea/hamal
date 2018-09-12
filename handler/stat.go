package handler

import (
	"log"
	"net/http"

	"hamal/protocol"

	"hamal/db"

	"github.com/gin-gonic/gin"
)

func GetDailyStats(c *gin.Context) {
	user := getUser(c)

	dailyStats, err := db.GetDailyStats(user.ID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, protocol.ErrorRsp(err))
		return
	}

	var data protocol.GetDailyStatsRsp
	for _, d := range dailyStats {
		var item protocol.GetDailyStatsRspItem
		item.Date = d.Date
		item.UpdatedAt = d.UpdatedAt
		item.USDTBalance = d.USDTBalance
		item.BTCBalance = d.BTCBalance

		data.Items = append(data.Items, item)
	}

	c.JSON(http.StatusOK, protocol.OKRsp(data))

}
