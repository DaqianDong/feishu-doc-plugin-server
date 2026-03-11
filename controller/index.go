package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 标准返回结构
type Response struct {
	Code    int         // 状态码，表示类似 state、status、ret 之类的值
	Message string      // 信息，可用于直接返回字符串的情况
	Data    interface{} // 成功后返回的数据
}

// Handler 处理器1
type Handler func(*gin.Context) (rsp any, err error)

func Wrap(handler Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if ok {
					c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
				} else {
					c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprint(r)})
				}
			}
		}()

		rsp, err := handler(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: err.Error()})
		} else {
			c.JSON(http.StatusOK, Response{Data: rsp})
		}
		log.Printf("api: %s, rsp: %v, err: %v \n", c.Request.URL, rsp, err)
	}
}
