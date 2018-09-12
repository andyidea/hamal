package market

import (
	"log"

	"time"

	"github.com/leek-box/sheep/fcoin"
	"github.com/leek-box/sheep/huobi"
	"github.com/leek-box/sheep/proto"
)

//
//type MarketTradeDetail struct {
//	Ch   string `json:"ch"`
//	Tick struct {
//		Data []struct {
//			Amount    float64 `json:"amount"`
//			Direction string  `json:"direction"`
//			Price     float64 `json:"price"`
//			TS        int64   `json:"ts"`
//		} `json:"data"`
//	} `json:"tick"`
//}
//
//func (m *MarketTradeDetail) String() string {
//	return fmt.Sprintln(m.Ch, "实时价格推送  价格:", m.Tick.Data[0].Price, " 数量:", m.Tick.Data[0].Amount, " 买卖：", m.Tick.Data[0].Direction)
//}
//
//type MarketDepth struct {
//	Ch   string `json:"ch"`
//	Tick struct {
//		Asks [][]float64 `json:"asks"`
//		Bids [][]float64 `json:"bids"`
//		TS   int64       `json:"ts"`
//	} `json:"tick"`
//}
//
//type Marketer struct {
//	huobiMarket *huobiapi.Market
//}
//
//func NewMarket() (*Marketer, error) {
//	// 创建行情实例
//	var m = Marketer{}
//	var err error
//	m.huobiMarket, err = huobiapi.NewMarket()
//	if err != nil {
//		return nil, err
//	}
//
//	return &m, nil
//}
//
//func (m *Marketer) Loop() {
//	m.huobiMarket.Loop()
//}
//
//func (m *Marketer) Close() {
//	m.huobiMarket.Close()
//}
//func (m *Marketer) Reconnect() {
//	m.huobiMarket.ReConnect()
//}
//
//// Listener 订阅事件监听器
//type DetailListener = func(symbol string, detail *MarketTradeDetail)
//
//func (m *Marketer) SubscribeDetail(listener DetailListener, symbols ...string) {
//	for _, symbol := range symbols {
//		m.huobiMarket.Subscribe("market."+symbol+".trade.detail", func(topic string, j *huobiapi.JSON) {
//			js, _ := j.MarshalJSON()
//			var mtd MarketTradeDetail
//			err := json.Unmarshal(js, &mtd)
//			if err != nil {
//				fmt.Println(err.Error())
//			}
//
//			ts := strings.Split(topic, ".")
//			listener(ts[1], &mtd)
//		})
//	}
//
//}
//
//// Listener 订阅事件监听器
//type DepthlListener = func(symbol string, depth *MarketDepth)
//
//func (m *Marketer) SubscribeDepth(listener DepthlListener, symbols ...string) {
//	for _, symbol := range symbols {
//		m.huobiMarket.Subscribe("market."+symbol+".depth.step0", func(topic string, j *huobiapi.JSON) {
//			js, _ := j.MarshalJSON()
//			var md = MarketDepth{}
//			err := json.Unmarshal(js, &md)
//			if err != nil {
//				fmt.Println(err.Error())
//			}
//
//			ts := strings.Split(topic, ".")
//			listener(ts[1], &md)
//		})
//	}
//
//}
//
//func (m *Marketer) UnSubscribeDepth(symbols ...string) {
//	for _, symbol := range symbols {
//		m.huobiMarket.Unsubscribe("market." + symbol + ".depth.step0")
//	}
//
//}
//
var HuobiAPI *huobi.Huobi

func LaunchMarket() {
	go func() {
		for {
			var pa2 proto.MarketDepthParams
			pa2.Symbol = "ftusdt"
			pa2.Level = "L20"
			marketDepth, err := fcoin.GetMarketDepth(&pa2)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			if marketDepth.Status != 0 {
				continue
			}

			var depth huobi.MarketDepth
			if len(marketDepth.Data.Bids) > 2 {
				depth.Tick.Bids = append(depth.Tick.Bids, []float64{marketDepth.Data.Bids[0], marketDepth.Data.Bids[1]})
			}

			if len(marketDepth.Data.Asks) > 2 {
				depth.Tick.Asks = append(depth.Tick.Bids, []float64{marketDepth.Data.Asks[0], marketDepth.Data.Asks[1]})
			}

			depth.Tick.TS = marketDepth.Data.TS
			SetCurrentDepth("ftusdt", &depth)
			time.Sleep(1 * time.Second)
		}
	}()
}

