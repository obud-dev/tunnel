package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeOK        = 0
	BaseCodeError = -1
	MsgOk         = "OK"
)

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

func New(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func Response(c *gin.Context, data any, err error) {
	if err != nil {
		codeError, ok := err.(*CodeError)
		if !ok {
			codeError = &CodeError{Code: BaseCodeError, Msg: err.Error()}
		}
		c.JSON(http.StatusOK, Result{
			Code: codeError.Code,
			Msg:  codeError.Msg,
		})
		return
	}
	c.JSON(http.StatusOK, Result{
		Code: CodeOK,
		Msg:  MsgOk,
		Data: data,
	})
}
