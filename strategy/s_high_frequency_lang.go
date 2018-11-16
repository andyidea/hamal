//我佛慈悲，祝我发财
package strategy

import (
	"hamal/config"
	"hamal/models"
	"hamal/util"
	"log"
	"time"

	"hamal/db"
	"net/http"

	"github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/huobi"
)

const (
	S1StatusWaitChance  = "wait_chance"  //等待时机
	S1StatusBuyApplied  = "buy_applied"  //购买申请成功
	S1StatusBuySuccess  = "buy_success"  //购买成功
	S1StatusSellApplied = "sell_applied" //出售申请成功
)

type StrategyHighFrequencyLang struct {
	StrategyInstanceID uint
	UserID             uint
	models.Strategy
	stop          bool
	interval      float64
	amount        float64
	trader        *huobi.HuoBiPro
	PrincipalUsdt float64
}

func (s *StrategyHighFrequencyLang) Init(si *models.StrategyInstance) {
	s.StrategyInstanceID = si.ID
	s.UserID = si.UserID

}

func (s *StrategyHighFrequencyLang) ID() uint {
	return s.StrategyInstanceID
}

func (s *StrategyHighFrequencyLang) Launch() {
	for {
		if s.stop {
			break
		}

		time.Sleep(5 * time.Second)

		var coinSymbol = goex.BTC
		var currencySymbol = goex.USDT
		var pair = goex.CurrencyPair{coinSymbol, currencySymbol}
		var interval = s.interval
		var amount = s.amount

		//获取深度
		depth, err := s.trader.GetDepth(10, pair)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		ask1 := depth.AskList[len(depth.AskList)-1].Price //卖1
		bid1 := depth.BidList[0].Price                    //买1

		log.Println(pair, "卖1：", ask1, "买1：", bid1)

		balance, err := s.trader.GetAccount()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		log.Println(pair, "账户信息：",
			currencySymbol, util.Float2String(balance.SubAccounts[currencySymbol].Amount), "冻结余额：", util.Float2String(balance.SubAccounts[currencySymbol].ForzenAmount),
			coinSymbol, util.Float2String(balance.SubAccounts[coinSymbol].Amount), "冻结余额：", util.Float2String(balance.SubAccounts[coinSymbol].ForzenAmount))

		profit := (balance.SubAccounts[coinSymbol].Amount+balance.SubAccounts[coinSymbol].ForzenAmount)*bid1 + balance.SubAccounts[currencySymbol].Amount + balance.SubAccounts[currencySymbol].ForzenAmount - s.PrincipalUsdt

		log.Println(pair, "本金：", s.PrincipalUsdt, "盈利：", profit)
		unfinishOrders, err := s.trader.GetUnfinishOrders(pair)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		//log.Println(unfinishOrders)

		//historyOrders, err := s.trader.GetOrderHistorys(pair, 0, 10)
		//if err != nil {
		//	log.Println(err.Error())
		//	continue
		//}

		//log.Println(historyOrders)

		if len(unfinishOrders) > 0 {
			//是否需要撤销
			var buyCount = 0
			var sellCount = 0
			var minSellPrice = 999999999.99
			for _, uforder := range unfinishOrders {
				if uforder.Side == goex.BUY {
					buyCount++
					log.Println(pair, "正在购买中，购买价：", uforder.Price)
					if uforder.Price+interval < bid1 {
						_, err = s.trader.CancelOrder(uforder.OrderID2, pair)
						if err != nil {
							log.Println(err)
							continue
						}
					}
				} else {
					sellCount++
					//log.Println(uforder)
					if minSellPrice > uforder.Price {
						minSellPrice = uforder.Price
					}
				}
			}

			log.Println(pair, "最小卖出价：", minSellPrice)

			if buyCount == 0 {
				if balance.SubAccounts[coinSymbol].Amount > amount {
					_, err = s.trader.LimitSell(util.Float2String(amount), util.Float2String(ask1+interval), pair)
					if err != nil {
						log.Println(err.Error())
						continue
					}
				} else if minSellPrice-interval*2 > bid1 {
					_, err = s.trader.LimitBuy(util.Float2String(amount), util.Float2String(bid1), pair)
					if err != nil {
						log.Println(err.Error())
						continue
					}
				} else {
					log.Println(pair, "距离下次购买差价：", bid1-(minSellPrice-interval*2))
				}

			}
		} else {
			if balance.SubAccounts[coinSymbol].Amount > amount {
				_, err = s.trader.LimitSell(util.Float2String(amount), util.Float2String(ask1+interval), pair)
				if err != nil {
					log.Println(err.Error())
					continue
				}
			} else {
				_, err = s.trader.LimitBuy(util.Float2String(amount), util.Float2String(bid1), pair)
				if err != nil {
					log.Println(err.Error())
					continue
				}
			}

		}

	}

}

func (s *StrategyHighFrequencyLang) Stop() {
	s.stop = true
}

func NewStrategyHighFrequencyLang(strategyInstanceID uint) (*StrategyHighFrequencyLang, error) {
	s := StrategyHighFrequencyLang{StrategyInstanceID: strategyInstanceID}

	si, err := db.GetStrategyInstanceModel(strategyInstanceID)
	if err != nil {
		return nil, err
	}

	ue, err := db.GetUserExchange(si.UserExchangeID)
	if err != nil {
		return nil, err
	}

	s.interval = si.Interval
	s.amount = si.Amount

	_client := http.DefaultClient
	transport := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 4 * time.Second,
	}
	_client.Transport = transport
	log.Println(ue.AccessKey, ue.SecretKey)
	huobi := huobi.NewHuoBiProSpot(_client, ue.AccessKey, ue.SecretKey)
	s.trader = huobi

	return &s, nil
}

func NewStrategyHighFrequencyLang2(c *config.Config) (*StrategyHighFrequencyLang, error) {
	s := StrategyHighFrequencyLang{}

	s.interval = c.Interval
	s.amount = c.Amount
	s.PrincipalUsdt = c.PrincipalUsdt

	_client := http.DefaultClient
	transport := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 4 * time.Second,
	}
	_client.Transport = transport
	huobi := huobi.NewHuoBiProSpot(_client, c.ApiKey, c.ApiSecret)
	s.trader = huobi

	return &s, nil
}
