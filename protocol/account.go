package protocol

type GetAccountBalanceRspItem struct {
	Currency      string  `json:"currency"` // 币种
	Balance       float64 `json:"balance"`
	TradeBalance  float64 `json:"trade_balance"` // 结余
	FrozenBalance float64 `json:"frozen_balance"`
	USD           float64 `json:"usd"`
}

type GetAccountBalanceRsp struct {
	Count int                        `json:"count"`
	Items []GetAccountBalanceRspItem `json:"items"`
}
