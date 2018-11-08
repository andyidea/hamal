//我佛慈悲，祝我发财
package strategy

import (
	"log"

	"time"

	"strconv"

	"math"

	"github.com/jinzhu/gorm"
	"github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/coinex"
	"hamal/db"
	"hamal/models"
	"net/http"
)

type StrategyCoinEXDuiQiao struct {
	StrategyInstanceID uint
	stop               bool
	trader             *coinex.CoinEx
}

func (s *StrategyCoinEXDuiQiao) ID() uint {
	return s.StrategyInstanceID
}

func (s *StrategyCoinEXDuiQiao) Launch() {
	s.stop = false

	go s.Loop()
}

func (s *StrategyCoinEXDuiQiao) Loop() {
	for {
		if s.stop {
			break
		}

		time.Sleep(5 * time.Second)

		err := s.RecardOrder()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		full, err := s.CheckTradeLimit()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if full {
			log.Println("该小时已达上限")
			continue
		}

		var cetSymbol = goex.Currency{"CET", ""}
		var coinSymbol = goex.BTC
		var currencySymbol = goex.USDT
		var pair = goex.CurrencyPair{coinSymbol, currencySymbol}

		pendingList, err := s.trader.GetPendingOrders(1, 100, pair)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		log.Println(pendingList)
		if len(pendingList) > 0 {
			log.Println("处理未完成的订单")
			item := pendingList[0]
			time.Sleep(5 * time.Second)

			info, err := s.trader.GetOneOrder(item.OrderID2, pair)
			if err != nil {
				log.Println(err)
				continue
			}

			if info.Status == goex.ORDER_FINISH {
				continue
			}

			_, err = s.trader.CancelOrder(item.OrderID2, pair)
			if err != nil {
				log.Println(err)
				continue
			}

		}

		//获取账户余额
		log.Println("获取账户余额")
		balance, err := s.trader.GetAccount()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		log.Println(balance)

		depth, err := s.trader.GetDepth(10, pair)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		//if depth.UTime.Unix() < now.Unix()-5 || depth.UTime.Unix() > now.Unix()-5 {
		//	log.Println("深度数据时间异常")
		//	continue
		//}

		ask1 := depth.AskList[len(depth.AskList)-1].Price //卖1
		bid1 := depth.BidList[0].Price                    //买1

		log.Println(ask1, bid1)

		var price = (bid1 + ask1) / 2

		var coinSymbolBalance float64
		var coinSymbolBalanceBaseCurrency float64
		var currencySymbolBalance float64
		for _, asset := range balance.SubAccounts {
			if coinSymbol == asset.Currency {
				coinSymbolBalance = asset.Amount
				coinSymbolBalanceBaseCurrency = asset.Amount * price
				log.Println(asset.Amount, price)
			}

			if currencySymbol == asset.Currency {
				currencySymbolBalance = asset.Amount
			}

			if cetSymbol == asset.Currency {
				if asset.Amount > 100 {
					s.SellCet()
					continue
				}
			}
		}

		var invUSD = coinSymbolBalanceBaseCurrency - currencySymbolBalance
		if invUSD > 0 {
			if invUSD > coinSymbolBalanceBaseCurrency*0.2 {
				//需要平仓
				_, err = s.trader.LimitSell(strconv.FormatFloat(invUSD/2/ask1, 'f', -1, 64), strconv.FormatFloat(ask1, 'f', -1, 64), pair)
				if err != nil {
					log.Println("平仓失败：", err.Error())
					continue
				}

				continue

			}
		} else if invUSD < 0 {

			invUSD = math.Abs(invUSD)
			if invUSD > coinSymbolBalanceBaseCurrency*0.2 {
				//需要平仓
				log.Println("发起平仓请求")
				_, err = s.trader.LimitBuy(strconv.FormatFloat(invUSD/2/bid1, 'f', -1, 64), strconv.FormatFloat(bid1, 'f', -1, 64), pair)
				if err != nil {
					log.Println("平仓失败：", err.Error())
					continue
				}

				continue

			}
		} else {
			continue
		}

		var amount float64
		//自成交
		log.Println("开始自成交")
		if invUSD > 0 {
			amount = currencySymbolBalance / price
		} else {
			amount = coinSymbolBalance
		}

		_, err = s.trader.LimitBuy(strconv.FormatFloat(amount, 'f', -1, 64), strconv.FormatFloat(price, 'f', -1, 64), pair)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		_, err = s.trader.LimitSell(strconv.FormatFloat(amount, 'f', -1, 64), strconv.FormatFloat(price, 'f', -1, 64), pair)
		if err != nil {
			log.Println(err.Error())
			continue
		}

	}
}

//记录订单信息系
func (s *StrategyCoinEXDuiQiao) RecardOrder() error {

	var coinSymbol = goex.BTC
	var currencySymbol = goex.USDT
	var pair = goex.CurrencyPair{coinSymbol, currencySymbol}

	orders, err := s.trader.GetOrderHistorys(pair, 1, 10000)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	cetPrice, err := s.GetCetPrice()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, order := range orders {
		_, err := db.GetStrategyInstanceRecardByOrderID(order.OrderID2)
		if err == nil {
			continue
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		orderTime := time.Unix(int64(order.OrderTime), 0)

		if order.Status == goex.ORDER_FINISH || order.Status == goex.ORDER_PART_FINISH {

			var fee float64
			if order.Side == goex.SELL {
				fee = order.Fee / cetPrice
			} else {
				fee = order.Fee * order.AvgPrice / cetPrice
			}
			var recard models.StrategyInstanceRecard
			recard.OrderID = order.OrderID2
			recard.Amount = order.Amount
			recard.Price = order.AvgPrice
			recard.Fee = fee
			recard.StrategyInstanceID = s.StrategyInstanceID
			recard.RecardType = "TRANSACTION"
			recard.OrderTime = orderTime
			recard.OrderStatus = int(order.Status)
			err = db.AddStrategyInstanceRecard(&recard)
			if err != nil {
				log.Println(err.Error())
				continue
			}
		}
	}
	return nil
}

func (s *StrategyCoinEXDuiQiao) SellCet() {
	var coinSymbol = goex.Currency{"CET", ""}
	var currencySymbol = goex.USDT
	var pair = goex.CurrencyPair{coinSymbol, currencySymbol}

	pendingList, err := s.trader.GetPendingOrders(1, 100, pair)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(pendingList)
	if len(pendingList) > 0 {
		log.Println("处理未完成的订单")
		item := pendingList[0]
		time.Sleep(5 * time.Second)

		info, err := s.trader.GetOneOrder(item.OrderID2, pair)
		if err != nil {
			log.Println(err)
			return
		}

		if info.Status == goex.ORDER_FINISH {
			return
		}

		_, err = s.trader.CancelOrder(item.OrderID2, pair)
		if err != nil {
			log.Println(err)
			return
		}

	}

	//获取账户余额
	log.Println("获取账户余额")
	balance, err := s.trader.GetAccount()
	if err != nil {
		log.Println(err.Error())
		return
	}

	depth, err := s.trader.GetDepth(10, pair)
	if err != nil {
		log.Println(err.Error())
		return
	}

	price := depth.BidList[0].Price //买1

	for _, asset := range balance.SubAccounts {
		if coinSymbol == asset.Currency {
			log.Println("准备卖出cet，数量：", asset.Amount, "，单价：", price)
			_, err = s.trader.LimitSell(strconv.FormatFloat(asset.Amount, 'f', -1, 64), strconv.FormatFloat(price, 'f', -1, 64), pair)
			if err != nil {
				log.Println(err.Error())
				return
			}
			break
		}
	}

}

//对账
func (s *StrategyCoinEXDuiQiao) CheckTradeLimit() (bool, error) {

	md, err := s.trader.GetMiningDifficulty()
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	var params db.GetStrategyInstanceRecardsParams
	params.StrategyInstanceID = s.StrategyInstanceID

	recards, err := db.GetStrategyInstanceRecards(&params)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	var now = time.Now()
	var yearDay = now.YearDay()
	var hour = time.Now().Hour()

	var feeAmount float64

	for _, recard := range recards {
		orderTime := recard.OrderTime

		if yearDay != orderTime.YearDay() || hour != orderTime.Hour() {
			continue
		}
		if recard.OrderStatus == goex.ORDER_FINISH || recard.OrderStatus == goex.ORDER_PART_FINISH {

			feeAmount += recard.Fee
		}
	}

	log.Println("本小时已挖到cet:", feeAmount, ". 总量：", md.Difficulty)

	if feeAmount > md.Difficulty {
		return true, nil
	}

	return false, nil

}

func (s *StrategyCoinEXDuiQiao) GetCetPrice() (float64, error) {
	var coinSymbol = goex.Currency{"CET", ""}
	var currencySymbol = goex.USDT
	var pair = goex.CurrencyPair{coinSymbol, currencySymbol}

	ticker, err := s.trader.GetTicker(pair)
	if err != nil {
		return 0, nil
	}

	return ticker.Last, nil

}

func (s *StrategyCoinEXDuiQiao) Stop() {
	s.stop = true
}

func NewStrategyCoinEXDuiQiao(strategyInstanceID uint) (*StrategyCoinEXDuiQiao, error) {
	s := StrategyCoinEXDuiQiao{StrategyInstanceID: strategyInstanceID}

	si, err := db.GetStrategyInstanceModel(strategyInstanceID)
	if err != nil {
		return nil, err
	}

	ue, err := db.GetUserExchange(si.UserExchangeID)
	if err != nil {
		return nil, err
	}

	_client := http.DefaultClient
	transport := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 4 * time.Second,
	}
	_client.Transport = transport
	s.trader = coinex.New(_client, ue.AccessKey, ue.SecretKey)

	return &s, nil
}
