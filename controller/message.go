package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/taoshihan1991/imaptool/common"
	"github.com/taoshihan1991/imaptool/models"
	"github.com/taoshihan1991/imaptool/tools"
	"github.com/taoshihan1991/imaptool/ws"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func SendMessageV2(c *gin.Context) {
	fromId := c.PostForm("from_id")
	toId := c.PostForm("to_id")
	content := c.PostForm("content")
	cType := c.PostForm("type")
	if content == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "内容不能为空",
		})
		return
	}
	//限流
	if !tools.LimitFreqSingle("sendmessage:"+c.ClientIP(), 1, 2) {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  c.ClientIP() + "发送频率过快",
		})
		return
	}
	var kefuInfo models.User
	var vistorInfo models.Visitor
	if cType == "kefu" {
		kefuInfo = models.FindUser(fromId)
		vistorInfo = models.FindVisitorByVistorId(toId)
	} else if cType == "visitor" {
		vistorInfo = models.FindVisitorByVistorId(fromId)
		kefuInfo = models.FindUser(toId)
	}

	if kefuInfo.ID == 0 || vistorInfo.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "用户不存在",
		})
		return
	}

	models.CreateMessage(kefuInfo.Name, vistorInfo.VisitorId, content, cType)
	//var msg TypeMessage
	if cType == "kefu" {
		guest, ok := ws.ClientList[vistorInfo.VisitorId]

		if guest != nil && ok {
			//客服发消息给访问者
			ws.VisitorMessage(vistorInfo.VisitorId, content, kefuInfo)
		}
		//主要作用是为了客服端页面展示
		ws.KefuMessage(vistorInfo.VisitorId, content, kefuInfo)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "ok",
		})
	}
	if cType == "visitor" {
		//访客发消息给客服
		guest, ok := ws.ClientList[vistorInfo.VisitorId]
		if ok && guest != nil {
			guest.UpdateTime = time.Now()
		}
		//发送给客服的消息
		msg := ws.TypeMessage{
			Type: "message",
			Data: ws.ClientMessage{
				Avator:  vistorInfo.Avator,
				Id:      vistorInfo.VisitorId,
				Name:    vistorInfo.Name,
				ToId:    kefuInfo.Name,
				Content: content,
				Time:    time.Now().Format("2006-01-02 15:04:05"),
				IsKefu:  "no",
			},
		}
		str, _ := json.Marshal(msg)
		//发送消息给客服
		ws.OneKefuMessage(kefuInfo.Name, str)
		kefu, ok := ws.KefuList[kefuInfo.Name]
		if !ok || kefu == nil {
			//客服不在线时，发送提醒邮件
			go SendNoticeEmail(content+"|"+vistorInfo.Name, content)
		}
		//客服自动回复访问者
		go ws.VisitorAutoReply(vistorInfo, kefuInfo, content)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "ok",
		})
	}
}

func SendVisitorNotice(c *gin.Context) {
	notice := c.Query("msg")
	if notice == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "msg不能为空",
		})
		return
	}
	msg := ws.TypeMessage{
		Type: "notice",
		Data: notice,
	}
	str, _ := json.Marshal(msg)
	for _, visitor := range ws.ClientList {
		visitor.Conn.WriteMessage(websocket.TextMessage, str)
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func SendCloseMessageV2(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	if visitorId == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "visitor_id不能为空",
		})
		return
	}

	oldUser, ok := ws.ClientList[visitorId]
	if oldUser != nil || ok {
		msg := ws.TypeMessage{
			Type: "force_close",
			Data: visitorId,
		}
		str, _ := json.Marshal(msg)
		err := oldUser.Conn.WriteMessage(websocket.TextMessage, str)
		oldUser.Conn.Close()
		delete(ws.ClientList, visitorId)
		tools.Logger().Println("close_message", oldUser, err)
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func UploadImg(c *gin.Context) {
	f, err := c.FormFile("imgfile")
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	} else {
		fileExt := strings.ToLower(path.Ext(f.Filename))
		if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "上传失败!只允许png,jpg,gif,jpeg文件",
			})
			return
		}
		isMainUploadExist, _ := tools.IsFileExist(common.Upload)
		if !isMainUploadExist {
			os.Mkdir(common.Upload, os.ModePerm)
		}
		fileName := tools.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := fmt.Sprintf("%s%d%s/", common.Upload, time.Now().Year(), time.Now().Month().String())
		isExist, _ := tools.IsFileExist(fildDir)
		if !isExist {
			os.Mkdir(fildDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "上传成功!",
			"result": gin.H{
				"path": filepath,
			},
		})
	}
}
func UploadFile(c *gin.Context) {
	f, err := c.FormFile("realfile")
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	} else {

		fileExt := strings.ToLower(path.Ext(f.Filename))
		if f.Size >= 90*1024*1024 {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "上传失败!不允许超过90M",
			})
			return
		}

		fileName := tools.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := fmt.Sprintf("%s%d%s/", common.Upload, time.Now().Year(), time.Now().Month().String())
		isExist, _ := tools.IsFileExist(fildDir)
		if !isExist {
			os.Mkdir(fildDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "上传成功!",
			"result": gin.H{
				"path": filepath,
				"ext":  fileExt,
				"size": f.Size,
				"name": f.Filename,
			},
		})
	}
}
func GetMessagesV2(c *gin.Context) {
	visitorId := c.Query("visitorId")
	messages := models.FindAllKefuMessageByVisitorId("message.visitor_id = ?", visitorId)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": messages,
	})
}
func GetMessagespages(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pagesize", "10"))
	if pageSize > 20 {
		pageSize = 20
	}
	count := models.CountMessage("visitor_id = ?", visitorId)
	list := models.FindMessageByPage(uint(page), uint(pageSize), "message.visitor_id = ?", visitorId)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"count":    count,
			"page":     page,
			"list":     list,
			"pagesize": pageSize,
		},
	})
}
