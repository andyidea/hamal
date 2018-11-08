//我佛慈悲，祝我发财
package strategy

import (
	"hamal/models"
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
	stop     bool
	interval float64
	amount   float64
	trader   *huobi.HuoBiPro
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

		//获取深度
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

func NewStrategyHighFrequencyLang2() (*StrategyHighFrequencyLang, error) {
	s := StrategyHighFrequencyLang{}

	s.interval = 500
	s.amount = 0.00001

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
