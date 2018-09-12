//我佛慈悲，祝我发财
package strategy

import (
	"hamal/models"
	"hamal/util"
	"log"
	"strings"
	"time"

	"hamal/db"
	"net/http"

	"github.com/leek-box/GoEx"
	"github.com/leek-box/GoEx/huobi"
	"github.com/leek-box/sheep/proto"
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
	stop     bool
	interval float64
	amount float64
	trader   *huobi.HuoBiPro
}

func (s *StrategyHighFrequencyLang) Init(si *models.StrategyInstance) {
	s.StrategyInstanceID = si.ID
	s.UserID = si.UserID

}

func (s *StrategyHighFrequencyLang) ID()uint {
	return s.StrategyInstanceID
}

//syA string     //流通货币
//syB string     //交易货币
//amount float64 //单笔交易数量
func (s *StrategyHighFrequencyLang) Launch() {
	// 初始设置
	//syA, syB, amount = CHATBTCConfig()
	s.stop = false
	log.Println("启动策略")

	// 各种状态
	var buyOrderID string  //买单ID
	var sellOrderID string //卖单ID
	var sellPrice float64  //当前卖价
	var buyPrice float64   //当前买家
	//var ask1 float64                //卖1
	var bid1 float64                    //买1
	var buyOrderState goex.TradeStatus  //买入订单状态
	var sellOrderState goex.TradeStatus //卖出订单状态
	var status = S1StatusWaitChance     //策略状态
	var partialAmount = -1.00           //部分成交
	var recardNextCheck = time.Now().Unix()
	var coinSymbol = goex.BTC
	var currencySymbol = goex.USDT
	var pair = goex.CurrencyPair{coinSymbol, currencySymbol}

	syA := currencySymbol.Symbol
	syB := coinSymbol.Symbol

	for {
		if s.stop {

			if status == S1StatusBuyApplied && buyOrderID != "" {
				log.Println("撤销订单")
				var params proto.OrderCancelParams
				params.OrderID = buyOrderID
				_, err := s.trader.CancelOrder(buyOrderID, pair)
				if err != nil {
					log.Println(err.Error())
					break
				}
				log.Println(syB, "撤销成功了")

			}
			break
		}
		time.Sleep(time.Second)

		now := time.Now()

		//检查记录
		if recardNextCheck <= now.Unix() {
			s.CheckRecard()

			recardNextCheck = now.Add(time.Minute * 30).Unix()
		}

		depth, err := s.trader.GetDepth(10, pair)
		if err != nil {
			log.Println(err.Error())
		}

		//depth, ok := market.GetCurrentDepth(syB + syA)
		//if ok {
		//	if depth.Tick.TS/1000 < now.Unix()-5 {
		//		log.Println(syB+syA, "时间大于5秒了，为：", now.Unix()-depth.Tick.TS/1000, "秒")
		//		continue
		//	}
		//} else {
		//	log.Println(syB+syA, "没有深度数据")
		//	continue
		//}

		//ask1 = depth.Tick.Asks[0][0]
		bid1 = depth.BidList[0].Price

		log.Println(bid1)

		switch status {
		case S1StatusWaitChance:
			//观察机会

			//获取当前委托的订单
			orders, err := s.trader.GetUnfinishOrders(pair)
			if err != nil {
				log.Println(err.Error())
				break
			}

			price := float64(int(bid1/s.interval)) * s.interval

			//判断是否有对应买价的单子
			if len(orders) > 0 {
				var orderID string
				for _, order := range orders {
					if order.Side == goex.BUY {
						orderID = order.OrderID2
						break
					}
				}

				if orderID != "" {
					log.Println(syB, "发现买单了")
					buyOrderID = orderID
					buyPrice =price
					sellPrice = price + s.interval
					//recard, err = service.GetStrategyInstanceRecardByBuyOrderID(buyOrderID)
					//if err != nil {
					//	//if err == gorm.ErrRecordNotFound {
					//	//	recard = &models.StrategyInstanceRecard{}
					//	//	recard.BuyOrderID = buyOrderID
					//	//	recard.BuyPrice = buyPrice
					//	//	recard.Amount = amount
					//	//	recard.Status = common.StrategyInstanceRecardStatusBuySuccess
					//	//	if err := service.AddStrategyInstanceRecard(recard); err != nil {
					//	//		log.Println(err.Error())
					//	//	}
					//	//}
					//	log.Println(err.Error())
					//}
					status = S1StatusBuyApplied
					break
				}
			}

			//判断是否有对应卖价的单子
			if len(orders) > 0 {
				var entrustPrice float64
				for _, order := range orders {
					if order.Side == goex.SELL {
						entrustPrice = orders[0].Price
						break
					}
				}
				if len(orders) > 0 && entrustPrice <= price+s.interval {
					log.Println(syB, "有对应卖单了，不用买")
					break
				}
			}

			log.Println(syB, "发现机会，尝试购买")

			var placeParams proto.OrderPlaceParams
			placeParams.Price = price //价格只支持到4位小数
			placeParams.Amount = s.amount
			placeParams.Type = proto.OrderPlaceTypeBuyLimit
			placeParams.BaseCurrencyID = syB
			placeParams.QuoteCurrencyID = syA
			ac, err := s.trader.GetAccount()
			if err !=  nil {
				log.Println(err.Error())
			}
			log.Println(ac)
			log.Println(util.Float2String(s.amount) , util.Float2String(price), pair)
			ret, err := s.trader.LimitBuy(util.Float2String(s.amount), util.Float2String(price), pair)
			if err != nil {
				log.Println(err.Error())
				break
			}
			buyOrderID = ret.OrderID2
			log.Println(syB, "买入订单提交成功, ID为", buyOrderID, "买入价：", util.Float2String(price))

			//设置状态
			buyPrice = price
			sellPrice = price + s.interval
			status = S1StatusBuyApplied
		case S1StatusBuyApplied:
			//等待韭菜卖给我
			log.Println(syB, "查询订单信息，订单ID", buyOrderID)
			var params proto.OrderInfoParams
			params.OrderID = buyOrderID
			buyOrderInfo, err := s.trader.GetOneOrder(buyOrderID, pair)
			if err != nil {
				log.Println(err.Error())
				break
			}

			log.Println(syB, "购买的订单查询成功，订单状态", buyOrderInfo.Status)
			buyOrderState = buyOrderInfo.Status

			switch buyOrderState {
			case goex.ORDER_FINISH:
				//购买完成，状态变更
				log.Println(syB, "购买的订单已成交")

				//recard = &models.StrategyInstanceRecard{}
				//recard.StrategyInstanceID = s.StrategyInstanceID
				//recard.UserID = s.UserID
				//recard.BuyOrderID = buyOrderID
				//recard.BuyPrice = buyPrice
				//recard.Amount = amount
				//recard.Status = common.StrategyInstanceRecardStatusBuySuccess
				//if err := db.AddStrategyInstanceRecard(recard); err != nil {
				//	log.Println(err.Error())
				//}

				buyOrderID = ""
				status = S1StatusBuySuccess
			case goex.ORDER_PART_FINISH:
				//部分完成，状态变更
				log.Println(syB, "购买的订单已部分成交")
				partialAmount = buyOrderInfo.DealAmount
				if err != nil {
					log.Println(err.Error())
					break
				}
				log.Println("剩下的就撤销吧")
				_, err := s.trader.CancelOrder(buyOrderID, pair)
				if err != nil {
					log.Println(err.Error())
					break
				}
				log.Println(syB, "撤销成功了")
				buyOrderID = ""
				status = S1StatusBuySuccess
			case goex.ORDER_CANCEL:
				//既然已经撤销就别耽误时间，赶紧下一波
				status = S1StatusWaitChance
			//case common.OrderStatePartialFilled:
			//部分成交
			default:
				price := float64(int(bid1/s.interval)) * s.interval
				if price >= buyPrice+s.interval {
					log.Println(syB, "fuck，价格涨了，申请撤销了。")
					var params proto.OrderCancelParams
					params.OrderID = buyOrderID
					_, err := s.trader.CancelOrder(buyOrderID, pair)
					if err != nil {
						log.Println(err.Error())
						break
					}
					log.Println(syB, "撤销成功了")

					//撤单申请成功。就当是撤单成功了 吧
					status = S1StatusWaitChance

				}
				break
			}
		case S1StatusBuySuccess:
			//买好了，该挂卖单了

			if sellPrice <= buyPrice {
				sellPrice = buyPrice + s.interval
			}

			log.Println(syB, "尝试卖出")
			var err error
			var orderAmount float64
			if partialAmount > 0 {
				orderAmount = partialAmount
			} else {
				orderAmount = s.amount
			}

			orderAmount = orderAmount

			var placeParams proto.OrderPlaceParams
			placeParams.Price = sellPrice
			placeParams.Amount = orderAmount
			placeParams.Type = proto.OrderPlaceTypeSellLimit
			placeParams.BaseCurrencyID = syB
			placeParams.QuoteCurrencyID = syA
			ret, err := s.trader.LimitSell(util.Float2String(orderAmount), util.Float2String(sellPrice), pair)
			if err != nil {
				log.Println(err.Error())
				if strings.Contains(err.Error(), "余额不足") {
					status = S1StatusWaitChance
					partialAmount = -1.00
					break
				} else if strings.Contains(err.Error(), "限价单交易数额错误") {
					log.Println(orderAmount)
				} else if err.Error() == "3006" {
					log.Println(sellPrice)
					log.Println(partialAmount)
					partialAmount = -1

				}
				break
			}
			sellOrderID = ret.OrderID2

			//if recard != nil {
			//	recard.SellOrderID = sellOrderID
			//	recard.SellPrice = sellPrice
			//	recard.Status = common.StrategyInstanceRecardStatusSellApplied
			//	err = db.UpdateStrategyInstanceRecard(recard)
			//	if err != nil {
			//		log.Println(err.Error())
			//	}
			//}

			partialAmount = -1.00
			log.Println(syB, "卖出订单提交成功，卖出价为", sellPrice)
			status = S1StatusSellApplied
		case S1StatusSellApplied:
			//等待韭菜买货
			status = S1StatusWaitChance
			break

			log.Println(syB, "查询订单信息，订单ID", sellOrderID)
			var params proto.OrderInfoParams
			params.OrderID = sellOrderID
			sellOrderInfo, err := s.trader.GetOneOrder(sellOrderID, pair)
			if err != nil {
				log.Println(err.Error())
				break
			}

			log.Println(syB, "卖出的订单查询成功，订单状态", sellOrderInfo.Status)
			sellOrderState = sellOrderInfo.Status

			switch sellOrderState {
			case goex.ORDER_FINISH:
				log.Println(syB, "卖的订单已成交，交易完成！")
				sellOrderID = ""
				log.Println(syB, "恭喜发财！买价：", buyPrice, "卖价：", sellPrice, "交易数量：", s.amount)
				status = S1StatusWaitChance
			case goex.ORDER_CANCEL:
				status = S1StatusWaitChance
			default:
				if bid1+s.interval < sellPrice {
					status = S1StatusWaitChance
				}
				break

			}
		default:
			log.Println("error:程序出问题了，不应该来这里！")
		}

	}

	log.Println("已停止")

}

func (s *StrategyHighFrequencyLang) CheckRecard() {
	//var params db.GetStrategyInstanceRecardsParams
	//params.StrategyInstanceID = s.StrategyInstanceID
	//params.Status = common.StrategyInstanceRecardStatusSellApplied
	//
	//recards, err := db.GetStrategyInstanceRecards(&params)
	//if err != nil {
	//	log.Println(err.Error())
	//	return
	//}
	//for _, recard := range recards {
	//	log.Println(recard.ID)
	//	var params proto.OrderInfoParams
	//	params.OrderID = recard.SellOrderID
	//	sellOrderInfo, err := s.trader.GetOrderInfo(&params)
	//	if err != nil {
	//		log.Println(err.Error())
	//		break
	//	}
	//	sellOrderState := sellOrderInfo.State
	//
	//	switch sellOrderState {
	//	case common.OrderStateFilled:
	//		recard.Status = common.StrategyInstanceRecardStatusDone
	//		fee := 0.002
	//		feeDiscount := 0.25
	//		recard.Profit = (recard.SellPrice * recard.Amount) - (recard.BuyPrice * recard.Amount) - (recard.BuyPrice*recard.Amount)*fee*feeDiscount - (recard.SellPrice*recard.Amount)*fee*feeDiscount
	//		err = db.UpdateStrategyInstanceRecard(&recard)
	//		if err != nil {
	//			log.Println(err.Error())
	//		}
	//		break
	//	case common.OrderStateCanceled:
	//		recard.Status = common.StrategyInstanceRecardStatusSellOrderCancel
	//		err = db.UpdateStrategyInstanceRecard(&recard)
	//		if err != nil {
	//			log.Println(err.Error())
	//		}
	//		break
	//	default:
	//		break
	//	}
	//}
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
	log.Println( ue.AccessKey, ue.SecretKey)
	huobi := huobi.NewHuoBiProSpot(_client, ue.AccessKey, ue.SecretKey)
	s.trader = huobi

	return &s, nil
}
