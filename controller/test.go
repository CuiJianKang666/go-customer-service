package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/taoshihan1991/imaptool/models"
)

func Test(c *gin.Context) {
	kefuId := c.Query("kefuId")
	title := c.Query("title")
	models.FindReplyItemByUserIdTitle(kefuId, title)
}
