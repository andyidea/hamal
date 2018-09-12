package protocol

type GetExchagesReq struct {
}

type GetExchangesRspItem struct {
	Name string `json:"name"`
}

type GetExchangesRsp struct {
	Count int                   `json:"count"`
	Items []GetExchangesRspItem `json:"items"`
}
