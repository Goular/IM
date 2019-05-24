package im

import (
	"IM/handler"
	"github.com/gin-gonic/gin"
)

//
func Hello(c *gin.Context) {
	handler.SendResponse(c, nil, nil)
}
