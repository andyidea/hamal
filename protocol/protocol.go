package protocol

var (
	CodeOK               = 0
	CodeSystemErr        = 1
	CodeWrongParam       = 2
	CodeTokenCheckFailed = 3
	CodeTokenOutDate     = 4
	CodeUsernameNoExist  = 5
	CodePasswordError    = 6
	CodeNoPower          = 7
)

var (
	MsgOK               = "操作成功"
	MsgSystemErr        = "系统错误"
	MsgWrongParam       = "参数错误"
	MsgTokenCheckFailed = "token校验失败"
	MsgTokenOutDate     = "token已过期"
	MSgUsernameNoExist  = "用户名不存在"
	MsgPasswordError    = "密码错误"
	MsgNoPower          = "没有权限"
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
