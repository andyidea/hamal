package stat

import (
	"hamal/util"
	"time"
)

func BalanceStat() {
	//users, err := db.GetUsers()
	//if err != nil {
	//	log.Println(err.Error())
	//	return
	//}
	//
	//now := time.Now()
	//
	//for _, user := range users {
	//	if user.HuobiAccessKey == "" {
	//		continue
	//	}
	//
	//	stat, err := db.GetOrAddDailyStat(user.ID, time.Now())
	//	if err != nil {
	//		log.Println(err.Error())
	//		continue
	//	}
	//
	//	trader, err := huobi.NewHuobi(user.HuobiAccessKey, user.HuobiSecretKey)
	//	if err != nil {
	//		log.Println(err.Error())
	//		continue
	//	}
	//
	//	var balance []proto.AccountBalance
	//	balance, err = trader.GetAccountBalance()
	//	if err != nil {
	//		log.Println(err.Error())
	//		return
	//
	//	}
	//
	//	var btcBalance float64
	//	var usdtBalance float64
	//
	//	log.Println("开始执行数据统计")
	//	for _, b := range balance {
	//		if b.Balance == "0.000000000000000000" {
	//			continue
	//		}
	//
	//		b.Currency = strings.ToUpper(b.Currency)
	//
	//		fb, _ := util.String2Float(b.Balance)
	//
	//		log.Println(b.Currency, market.GetLatestBTCPrice(b.Currency))
	//
	//		log.Println("处理数据", b.Currency, "btc价格:", market.GetLatestBTCPrice(b.Currency)*fb, "usdt价格：", market.GetLatestUSDTPrice(b.Currency)*fb)
	//		btcBalance += market.GetLatestBTCPrice(b.Currency) * fb
	//		usdtBalance += market.GetLatestUSDTPrice(b.Currency) * fb
	//	}
	//
	//	stat.BTCBalance = btcBalance
	//	stat.USDTBalance = usdtBalance
	//
	//	var params db.GetStrategyInstanceRecardsParams
	//	params.UserID = user.ID
	//	params.Status = common.StrategyInstanceRecardStatusDone
	//	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	//	params.UpdateTimeBegin = &date
	//
	//	recards, err := db.GetStrategyInstanceRecards(&params)
	//	if err != nil {
	//		log.Println(err.Error())
	//		return
	//	}
	//
	//	var profit float64
	//	for _, recard := range recards {
	//		profit += recard.Profit
	//	}
	//
	//	stat.Profit = profit
	//
	//	err = db.AddOrUpdateDailyStat(stat)
	//	if err != nil {
	//		log.Println(err.Error())
	//	}
	//
	//	log.Println("数据统计执行完成")
	//
	//}

}

// 启动统计
func LaunchStat() {
	util.StartTimer(time.Minute*10, BalanceStat)
}
