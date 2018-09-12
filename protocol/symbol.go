package protocol

type GetSymbolsRspItem struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	BaseCurrencyID  string  `json:"base_currency_id"`
	QuoteCurrencyID string  `json:"quote_currency_id"`
	LastetPrice     float64 `json:"lastet_price"`
	CreatedAt       int64   `json:"created_at"`
}

type GetSymbolsRsp struct {
	Count int                 `json:"count"`
	Items []GetSymbolsRspItem `json:"items"`
}

type AddOrUpdateSymbolsReq struct {
	ID              string `form:"id"`
	Name            string `form:"name"`
	BaseCurrencyID  string `form:"base_currency_id"`
	QuoteCurrencyID string `form:"quote_currency_id"`
}
