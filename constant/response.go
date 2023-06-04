package constant

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

const (
	SUCCESS        = http.StatusOK
	INVALID_PARAMS = http.StatusBadRequest
	ERROR          = http.StatusInternalServerError
)

const (
	SUCCESS_MSG = "Ok"
	ERROR_MSG   = "Fail"
)

func ResponseWithData(c *gin.Context, httpCode, respCode int, msg string, data interface{}, info ...interface{}) {
	resp := Response{
		Code:    respCode,
		Message: msg,
		Data:    data,
	}

	c.JSON(httpCode, resp)
}
