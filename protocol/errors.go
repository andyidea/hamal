package protocol

import "fmt"

type Error struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code:%d msg:%s detail:%s", e.Code, e.Msg, e.Detail)
}

func NewError(code int, msg string, detail string) Error {
	return Error{Code: code, Msg: msg, Detail: detail}
}

func ErrorRsp(reserr error) Error {
	err, ok := reserr.(Error)
	if ok {
		return err
	}

	return NewError(CodeSystemErr, MsgSystemErr, err.Error())

}
