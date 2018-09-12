package protocol

type GetStrategiesReq struct {
}

type GetStrategiesRspItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

type GetStrategiesRsp struct {
	Count int                    `json:"count"`
	Items []GetStrategiesRspItem `json:"items"`
}

type GetStrategyInstancesRspItem struct {
	ID               uint    `json:"id"`
	StrategyID       string  `json:"strategy_id"`
	StrategyName     string  `json:"strategy_name"`
	BaseCurrencyID   string  `json:"base_currency_id"`
	TargetCurrencyID string  `json:"target_currency_id"`
	UserExchangeID   uint    `json:"user_exchange_id"`
	Amount           float64 `json:"amount"`
	Interval         float64 `json:"interval"`
	Status           string  `json:"status"`
	CreatedAt        int64   `json:"created_at"`
}

type GetStrategyInstancesRsp struct {
	Count int                           `json:"count"`
	Items []GetStrategyInstancesRspItem `json:"items"`
}

type AddStrategyInstanceReq struct {
	StrategyID       string  `form:"strategy_id"`
	BaseCurrencyID   string  `form:"base_currency_id"`
	TargetCurrencyID string  `form:"target_currency_id"`
	UserExchangeID   uint    `form:"user_exchange_id"`
	Amount           float64 `form:"amount"`
	Interval         float64 `form:"interval"`
}

type AddStrategyInstanceRsp struct {
	GetStrategyInstancesRspItem
}

type UpdateStrategyInstanceReq struct {
	StrategyID       string  `form:"strategy_id"`
	BaseCurrencyID   string  `form:"base_currency_id"`
	TargetCurrencyID string  `form:"target_currency_id"`
	Amount           float64 `form:"amount"`
	Interval         float64 `form:"interval"`
}

type GetStrategyInstanceRsp struct {
	GetStrategyInstancesRspItem
}

type GetStrategyInstancesReq struct {
	StrategyID string `form:"strategy_id"`
}
