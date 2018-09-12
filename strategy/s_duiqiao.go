//我佛慈悲，祝我发财
package strategy

import (
	"hamal/models"

	"hamal/db"

	"log"

	"time"

	"strconv"

	"math"

	"github.com/leek-box/sheep/bibox"
)

const (
	OrderSideBuy  = "1"
	OrderSideSell = "2"
)

type StrategyDuiQiao struct {
	StrategyInstanceID uint
	Si                 *models.StrategyInstance
	Ue                 *models.UserExchange
	stop               bool
	trader             *bibox.Bibox
}

func (s *StrategyDuiQiao) Launch() {
	s.stop = false

	s.Loop()
}

func (s *StrategyDuiQiao) Loop() {
	for {
		if s.stop {
			break
		}

		time.Sleep(5 * time.Second)

		var pair = "BTC_USDT"
		var coinSymbol = "BTC"
		var currencySymbol = "USDT"
		var accountType = "0"

		pendingList, err := s.trader.GetOrderPendingList(pair, accountType, "1", "100", coinSymbol, currencySymbol, "")
		if err != nil {
			log.Println(err.Error())
			continue
		}

		log.Println(pendingList)

		if len(pendingList.Result) <= 0 {
			log.Println("异常")
			continue
		}
		if len(pendingList.Result[0].Result.Items) > 0 {
			log.Println("处理未完成的订单")
			item := pendingList.Result[0].Result.Items[0]
			time.Sleep(5 * time.Second)

			info, err := s.trader.GetOrderInfo(strconv.FormatInt(item.ID, 10))
			if err != nil {
				log.Println(err)
				continue
			}

			if info.Result[0].Result.Status == 3 {
				continue
			}

			err = s.trader.OrderCancel(strconv.FormatInt(item.ID, 10))
			if err != nil {
				log.Println(err)
				continue
			}

		}

		//获取账户余额
		log.Println("获取账户余额")
		balance, err := s.trader.GetAccountBalabce()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		log.Println(balance)

		if len(balance.Result) <= 0 {
			log.Println("异常")
			continue
		}
		var coinSymbolBalance float64
		var coinSymbolBalanceUSD float64
		var currencySymbolBalance float64
		var currencySymbolBalanceUSD float64
		for _, asset := range balance.Result[0].Result.AssetsList {
			if coinSymbol == asset.CoinSymbol {
				coinSymbolBalance, _ = strconv.ParseFloat(asset.Balance, 64)
				coinSymbolBalanceUSD, _ = strconv.ParseFloat(asset.USDValue, 64)
			}

			if currencySymbol == asset.CoinSymbol {
				currencySymbolBalance, _ = strconv.ParseFloat(asset.Balance, 64)
				currencySymbolBalanceUSD, _ = strconv.ParseFloat(asset.USDValue, 64)
			}
		}

		log.Println(coinSymbolBalance, currencySymbolBalance, coinSymbolBalanceUSD, currencySymbolBalanceUSD)

		depth, err := bibox.GetMarketDepth(pair)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		now := time.Now()
		if depth.Result.UpdateTime < now.Unix()-5 || depth.Result.UpdateTime > now.Unix()-5 {
			log.Println("深度数据时间异常")
		}

		ask1, err := strconv.ParseFloat(depth.Result.Asks[0].Price, 64) //卖1
		if err != nil {
			log.Println(err.Error())
			continue
		}

		bid1, err := strconv.ParseFloat(depth.Result.Bids[0].Price, 64) //买1
		if err != nil {
			log.Println(err.Error())
			continue
		}

		var invUSD = coinSymbolBalanceUSD - currencySymbolBalanceUSD
		if invUSD > 0 {
			if invUSD > coinSymbolBalanceUSD*0.2 {
				//需要平仓
				_, err = s.trader.OrderPlace(pair, accountType, "2", OrderSideSell, strconv.FormatFloat(ask1, 'f', -1, 64), strconv.FormatFloat(invUSD/2/ask1, 'f', -1, 64))
				if err != nil {
					log.Println("平仓失败：", err.Error())
					continue
				}

				continue

			}
		} else if invUSD < 0 {
			log.Println("发起平仓请求")
			invUSD = math.Abs(invUSD)
			if invUSD > currencySymbolBalanceUSD*0.2 {
				//需要平仓
				_, err = s.trader.OrderPlace(pair, accountType, "2", OrderSideBuy, strconv.FormatFloat(bid1, 'f', -1, 64), strconv.FormatFloat(invUSD/2/bid1, 'f', -1, 64))
				if err != nil {
					log.Println("平仓失败：", err.Error())
					continue
				}

				continue

			}
		} else {
			continue
		}

		var price = (bid1 + ask1) / 2
		var amount float64
		//自成交
		log.Println("开始自成交")
		if invUSD > 0 {
			amount = currencySymbolBalance / price
		} else {
			amount = coinSymbolBalance
		}

		_, err = s.trader.OrderPlace(pair, accountType, "2", OrderSideBuy, strconv.FormatFloat(price, 'f', -1, 64), strconv.FormatFloat(amount, 'f', -1, 64))
		if err != nil {
			log.Println(err.Error())
			continue
		}
		_, err = s.trader.OrderPlace(pair, accountType, "2", OrderSideSell, strconv.FormatFloat(price, 'f', -1, 64), strconv.FormatFloat(amount, 'f', -1, 64))
		if err != nil {
			log.Println(err.Error())
			continue
		}

	}
}

//对账
func (s *StrategyDuiQiao) SyncOrder() {
	//var pair = "BTC_USDT"
	//var coinSymbol = "BTC"
	//var currencySymbol = "USDT"
	//var accountType = "0"

	//orders, err := db.GetOrders(s.Ue.ID)
	//if err != nil {
	//	log.Println(err.Error())
	//	return
	//}
	//
	//var m = make(map[string]string)
	//
	//var i int
	//for i = 1; i < 1000; i++ {
	//	historyList, err := s.trader.GetOrderHistoryList(pair, accountType, strconv.Itoa(i), "100", coinSymbol, currencySymbol, "")
	//	if err != nil {
	//		log.Println(err.Error())
	//		return
	//	}
	//
	//	for _, item := range historyList.Result[0].Result.Items {
	//
	//	}
	//
	//	log.Println(historyList)
	//}
	//
	//historyList, err := s.trader.GetOrderHistoryList(pair, accountType, "1", "100", coinSymbol, currencySymbol, "")
	//if err != nil {
	//	log.Println(err.Error())
	//	return
	//}
	//
	//log.Println(historyList)

}

func (s *StrategyDuiQiao) Stop() {
	s.stop = true
}

func NewStrategyDuiQiao(strategyInstanceID uint) (*StrategyDuiQiao, error) {
	s := StrategyDuiQiao{StrategyInstanceID: strategyInstanceID}

	si, err := db.GetStrategyInstanceModel(strategyInstanceID)
	if err != nil {
		return nil, err
	}

	s.Si = si

	ue, err := db.GetUserExchange(si.UserExchangeID)
	if err != nil {
		return nil, err
	}

	s.Ue = ue

	s.trader, err = bibox.NewBibox(ue.AccessKey, ue.SecretKey)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
