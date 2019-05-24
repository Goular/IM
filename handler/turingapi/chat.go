package turingapi

import (
	"IM/handler"
	"IM/util/turingapi"
	"github.com/gin-gonic/gin"
)

// 与机器人进行文字聊天
func ChatBot(c *gin.Context) {
	text := c.PostForm("text")
	err, result := turingapi.ChatRobotWithText(text, nil)
	if err != nil {
		handler.SendResponse(c, err, nil)
	} else {
		handler.SendResponse(c, nil, turingapi.GetResponseTxt(result))
	}
}
