package response

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) { c.JSON(200, Result{Code: 200, Msg: "成功", Data: data}) }
func OKEmpty(c *gin.Context)              { OK(c, nil) }
func Fail(c *gin.Context, code int, msg string, data interface{}) {
	status := 400
	switch code {
	case CodeSystemError:
		status = 500
	case CodeArticleNotExist, CodeNotFound:
		status = 404
	case CodeConflict:
		status = 409
	case CodeTooLarge:
		status = 413
	}
	c.JSON(status, Result{Code: code, Msg: msg, Data: data})
}
func Internal(c *gin.Context, err error) {
	log.Printf("API error %s %s: %v", c.Request.Method, c.FullPath(), err)
	Fail(c, CodeSystemError, "服务暂时不可用，请稍后重试", nil)
}
func HTTPError(c *gin.Context, status int, msg string) {
	c.JSON(status, Result{Code: status, Msg: msg})
}

const (
	CodeParamInvalid    = 10001
	CodeParamBlank      = 10002
	CodeArticleNotExist = 50004
	CodeNotFound        = 50005
	CodeUploadError     = 80001
	CodeSystemError     = 40001
	CodeConflict        = 40009
	CodeTooLarge        = 40013
)
