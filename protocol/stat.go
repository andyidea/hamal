package protocol

import "time"

type GetDailyStatsRspItem struct {
	Date        time.Time `json:"date"`
	UpdatedAt   time.Time `json:"updated_at"`
	BTCBalance  float64   `json:"btc_balance"`
	USDTBalance float64   `json:"usdt_balance"`
}

type GetDailyStatsRsp struct {
	Count int                    `json:"count"`
	Items []GetDailyStatsRspItem `json:"items"`
}
