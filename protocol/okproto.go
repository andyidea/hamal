package protocol

type OK struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Detail string      `json:"detail"`
	Data   interface{} `json:"data"`
}

func OKRsp(data interface{}) OK {
	return OK{Code: CodeOK, Msg: MsgOK, Data: data}
}

func OKRspWithMsg(msg string, data interface{}) OK {
	return OK{Code: CodeOK, Msg: msg, Data: data}
}
