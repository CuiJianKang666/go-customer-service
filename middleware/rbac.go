package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/taoshihan1991/imaptool/models"
	"strings"
)

func RbacAuth(c *gin.Context) {
	for _, api := range apiWhiteList {
		if api == c.Request.URL.Path {
			c.Next()
			return
		}
	}
	roleId, _ := c.Get("role_id")
	role := models.FindRole(roleId)
	var flag bool
	rPaths := strings.Split(c.Request.RequestURI, "?")
	uriParam := fmt.Sprintf("%s:%s", c.Request.Method, rPaths[0])
	//权限以“，”分割
	if role.Method != "*" || role.Path != "*" {
		paths := strings.Split(role.Path, ",")
		for _, p := range paths {
			if uriParam == p {
				flag = true
				break
			}
		}
		if !flag {
			c.JSON(200, gin.H{
				"code": 403,
				"msg":  "没有权限:" + uriParam,
			})
			c.Abort()
			return
		}
	}
}
