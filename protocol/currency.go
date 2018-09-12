package protocol

type GetCurrenciesReq struct {
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	ID    string `form:"id"`
}

type GetCurrenciesRspItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

type GetCurrenciesRsp struct {
	Count int                    `json:"count"`
	Items []GetCurrenciesRspItem `json:"items"`
}

type AddOrUpdateCurrencyReq struct {
	ID   string
	Name string `form:"name"`
}

type AddOrUpdateCurrencyRsp struct {
}
