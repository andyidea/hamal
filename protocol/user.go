package protocol

type UserRegisterReq struct {
}

type UserRegisterRsp struct {
}

type LoginReq struct {
	AccountName string `json:"account_name" form:"account_name"`
	Password    string `json:"password" form:"password"`
}

type LoginRsp struct {
	Token string `json:"token"`
}

type GetUserInfoRsp struct {
	Name   string   `json:"name"`
	Avatar string   `json:"avatar"`
	Roles  []string `json:"roles"`
}

type UpdateUserInfoReq struct {
	Name   string `form:"name"`
	Avatar string `form:"avatar"`
}

type GetUserExchangesRspItem struct {
	ID        uint   `json:"id"`
	Exchange  string `json:"exchange"`
	Tag       string `json:"tag"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type GetUserExchangesRsp struct {
	Count int                       `json:"count"`
	Items []GetUserExchangesRspItem `json:"items"`
}

type AddOrUpdateUserExchangeReq struct {
	ID        uint   `form:"id"`
	Exchange  string `form:"exchange"`
	Tag       string `form:"tag"`
	AccessKey string `form:"access_key"`
	SecretKey string `form:"secret_key"`
}

type GetUserExchangesReq struct {
	Exchange string `form:"exchange"`
}

type DeleteUserExchangeReq struct {
	ID uint `form:"id"`
}
