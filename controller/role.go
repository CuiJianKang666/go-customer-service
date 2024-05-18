package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/taoshihan1991/imaptool/models"
	"strconv"
)

func GetRoleList(c *gin.Context) {
	roles := models.FindRoles()
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "获取成功",
		"result": roles,
	})
}
func PostRole(c *gin.Context) {
	roleId := c.PostForm("id")
	method := c.PostForm("method")
	name := c.PostForm("name")
	path := c.PostForm("path")
	if roleId == "" || method == "" || name == "" || path == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数不能为空",
		})
		return
	}
	models.SaveRole(roleId, name, method, path)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}
func CreateRole(c *gin.Context) {
	method := c.PostForm("method")
	name := c.PostForm("name")
	path := c.PostForm("path")
	if method == "" || name == "" || path == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数不能为空",
		})
		return
	}
	models.CreateRole(name, method, path)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "添加成功",
	})
}

func DeleteRole(c *gin.Context) {
	role_id := c.Query("role_id")
	if role_id == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数不能为空",
		})
		return
	}
	id, _ := strconv.Atoi(role_id)
	models.DeleteRole(id)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
