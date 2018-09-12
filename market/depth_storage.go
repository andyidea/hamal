package market

import (
	"sync"

	"github.com/leek-box/sheep/huobi"
)

var currentDepthlMap sync.Map

func SetCurrentDepth(symbol string, depth *huobi.MarketDepth) {
	currentDepthlMap.Store(symbol, depth)
}

func GetCurrentDepth(symbol string) (*huobi.MarketDepth, bool) {
	data, ok := currentDepthlMap.Load(symbol)
	if ok {
		return data.(*huobi.MarketDepth), true
	}
	return nil, false
}
