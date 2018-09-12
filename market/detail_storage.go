package market

import (
	"strings"
	"sync"

	"github.com/leek-box/sheep/huobi"
)

var currentDetailMap sync.Map

func SetCurrentDetail(symbol string, detail *huobi.MarketTradeDetail) {
	currentDetailMap.Store(symbol, detail)
}

func GetCurrentDetail(symbol string) (*huobi.MarketTradeDetail, bool) {
	data, ok := currentDetailMap.Load(symbol)
	if ok {
		return data.(*huobi.MarketTradeDetail), true
	}
	return nil, false
}

func GetLatestPrice(symbolID string) float64 {
	current, ok := GetCurrentDetail(symbolID)
	if ok {
		if len(current.Tick.Data) > 0 {
			return current.Tick.Data[0].Price
		}
	}

	return 0.00
}

func GetLatestBTCUSDTPrice() float64 {
	current, ok := GetCurrentDetail("btcusdt")
	if ok {
		if len(current.Tick.Data) > 0 {
			return current.Tick.Data[0].Price
		}
	}

	return 0.00
}

func GetLatestUSDTPrice(currencyID string) float64 {
	var usdt = "usdt"
	var btc = "btc"
	var lowcur = strings.ToLower(currencyID)
	if lowcur == usdt {
		return 1.00
	}
	current, ok := GetCurrentDetail(lowcur + usdt)
	if ok {
		if len(current.Tick.Data) > 0 {
			return current.Tick.Data[0].Price
		}
	}

	current, ok = GetCurrentDetail(lowcur + btc)
	if ok {
		if len(current.Tick.Data) > 0 {
			return current.Tick.Data[0].Price * GetLatestBTCUSDTPrice()
		}
	}

	return 0.00
}

func GetLatestBTCPrice(currencyID string) float64 {
	var usdt = "usdt"
	var btc = "btc"
	var lowcur = strings.ToLower(currencyID)
	if lowcur == btc {
		return 1.00
	}
	if lowcur == usdt {
		btcusdtPrice := GetLatestBTCUSDTPrice()
		return 1.00 / btcusdtPrice
	}

	current, ok := GetCurrentDetail(lowcur + btc)
	if ok {
		if len(current.Tick.Data) > 0 {
			return current.Tick.Data[0].Price
		}
	}

	return 0.00
}
