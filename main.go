package main

import (
	"hamal/config"
	_ "hamal/config"
	"hamal/strategy"
)

type Config struct {
	ApiKey    string  `json:"api_key"`
	ApiSecret string  `json:"api_secret"`
	Interval  float64 `json:"interval"`
	Amount    float64 `json:"amount"`
}

func main() {

	//log.SetFlags(log.Lshortfile | log.LstdFlags)

	//db.InitDB()

	//strategy.LaunchStrategyInstance()
	//
	//gin.SetMode(gin.DebugMode)
	//r := gin.Default()
	//
	//r.Use(middleware.Cors())
	//r.Use(middleware.AuthVerify())
	//
	//v1 := r.Group("v1")
	//
	////user
	////v1.POST("register", handler.Register)
	//v1.POST("user/login", handler.Login)
	//v1.GET("user/info", handler.GetUserInfo)
	//v1.PUT("user/info", handler.UpdateUserInfo)
	//v1.POST("user/logout", handler.Logout)
	//
	////user exchange
	//v1.GET("/user/exchanges", handler.GetUserExchanges)
	//v1.PUT("/user/exchanges", handler.AddOrUpdateUserExchange)
	//v1.DELETE("/user/exchanges", handler.DeleteUserExchange)
	//
	////account
	//v1.GET("account/balance", handler.GetAccountBalance)
	//
	////currency
	//v1.GET("currencies", handler.GetCurrencies)
	//v1.PUT("currencies/:id", handler.AddOrUpdateCurrency)
	//v1.DELETE("currencies/:id", handler.DeleteCurrency)
	//
	////symbol
	//v1.GET("symbols", handler.GetSymbols)
	//v1.PUT("symbols/:id", handler.AddOrUpdateSymbols)
	//
	////strategy
	//v1.GET("strategies", handler.GetStrategies)
	//v1.GET("strategy/instances", handler.GetStrategyInstances)
	//v1.POST("strategy/instances", handler.AddStrategyInstance)
	//v1.PUT("strategy/instances/:id", handler.UpdateStrategyInstance)
	//v1.DELETE("strategy/instances/:id", handler.DeleteStrategyInstance)
	//v1.PUT("strategy/instances/:id/launch", handler.LaunchStrategyInstance)
	//v1.PUT("strategy/instances/:id/stop", handler.StopStrategyInstance)
	//
	////exchage
	//v1.GET("exchanges", handler.GetExchanges)
	//
	////stat
	//v1.GET("stat/daily", handler.GetDailyStats)

	//s, err := strategy.NewStrategyCoinEXDuiQiao(1, "620EC939904341BCB6EC71A249F6D997", "0E6621086A6145D58A0196B6EFD75CA6CD6288E759F1DC84")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//s.Launch()

	//r.Run("0.0.0.0:8888")

	s, _ := strategy.NewStrategyHighFrequencyLang2(&config.Instance)
	s.Launch()
}
