package models

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type ReplyItem struct {
	Id       string `json:"item_id"`
	Content  string `json:"item_content"`
	GroupId  string `json:"group_id"`
	ItemName string `json:"item_name"`
	UserId   string `json:"user_id"`
}
type ReplyGroup struct {
	Id        string       `json:"group_id"`
	GroupName string       `json:"group_name"`
	UserId    string       `json:"user_id"`
	Items     []*ReplyItem `json:"items"`
}

type Reply_group struct {
	Id        string `json:"id"`
	GroupName string `json:"group_name"`
	UserId    string `json:"user_id"`
}

func (Reply_group) TableName() string {
	return "reply_group"
}

func FindReplyItemByUserIdTitle(userId interface{}, title string) ReplyItem {
	var reply []ReplyItem
	//OldDB.Where("user_id = ? and item_name = ?", userId, title).Find(&reply)
	OldDB.Where("user_id = ?", userId).Find(&reply)
	//使用大模型来做文本相似度对比
	var itemArr []string
	for _, item := range reply {
		itemArr = append(itemArr, item.ItemName)
	}
	client := resty.New()
	requestBody := map[string]interface{}{
		"sentences1": []string{title},
		"sentences2": itemArr,
	}
	// 发送 POST 请求
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post("http://localhost:8010/text_similarity")

	if err != nil {
		fmt.Println("Error:", err)
		return ReplyItem{}
	}

	// 解析响应 JSON 数据
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return ReplyItem{}
	}
	maxScore := result["res"].(map[string]interface{})["max_score"].(string)
	maxIndex := int(result["res"].(map[string]interface{})["max_index"].(float64))
	if maxScore > "0.9" {
		fmt.Println(reply[maxIndex])
		return reply[maxIndex]
	}
	return ReplyItem{}
}

func FindReplyByUserId(userId interface{}) []*ReplyGroup {
	var replyGroups []*ReplyGroup
	//OldDB.Raw("select a.*,b.* from reply_group a left join reply_item b on a.id=b.group_id where a.user_id=? ", userId).Scan(&replyGroups)
	var replyItems []*ReplyItem
	OldDB.Where("user_id = ?", userId).Find(&replyGroups)
	OldDB.Where("user_id = ?", userId).Find(&replyItems)
	temp := make(map[string]*ReplyGroup)
	for _, replyGroup := range replyGroups {
		replyGroup.Items = make([]*ReplyItem, 0)
		temp[replyGroup.Id] = replyGroup
	}
	for _, replyItem := range replyItems {
		temp[replyItem.GroupId].Items = append(temp[replyItem.GroupId].Items, replyItem)
	}
	return replyGroups
}
func FindReplyTitleByUserId(userId interface{}) []*ReplyGroup {
	var replyGroups []*ReplyGroup
	//OldDB.Raw("select a.*,b.* from reply_group a left join reply_item b on a.id=b.group_id where a.user_id=? ", userId).Scan(&replyGroups)
	var replyItems []*ReplyItem
	OldDB.Where("user_id = ?", userId).Find(&replyGroups)
	OldDB.Select("item_name,group_id").Where("user_id = ?", userId).Find(&replyItems)
	temp := make(map[string]*ReplyGroup)
	for _, replyGroup := range replyGroups {
		replyGroup.Items = make([]*ReplyItem, 0)
		temp[replyGroup.Id] = replyGroup
	}
	for _, replyItem := range replyItems {
		temp[replyItem.GroupId].Items = append(temp[replyItem.GroupId].Items, replyItem)
	}
	return replyGroups
}
func CreateReplyGroup(groupName string, userId string) {
	g := &ReplyGroup{
		GroupName: groupName,
		UserId:    userId,
	}
	OldDB.Create(g)
}
func CreateReplyContent(groupId string, userId string, content, itemName string) {
	g := &ReplyItem{
		GroupId:  groupId,
		UserId:   userId,
		Content:  content,
		ItemName: itemName,
	}
	OldDB.Create(g)
}
func UpdateReplyContent(id, userId, title, content string) {
	r := &ReplyItem{
		ItemName: title,
		Content:  content,
	}
	OldDB.Model(&ReplyItem{}).Where("user_id = ? and id = ?", userId, id).Update(r)
}
func DeleteReplyContent(id string, userId string) {
	OldDB.Where("user_id = ? and id = ?", userId, id).Delete(ReplyItem{})
}
func DeleteReplyGroup(id string, userId string) {
	OldDB.Where("user_id = ? and id = ?", userId, id).Delete(ReplyGroup{})
	OldDB.Where("user_id = ? and group_id = ?", userId, id).Delete(ReplyItem{})
}
func FindReplyBySearcch(userId interface{}, search string) []*ReplyGroup {
	var replyGroups []*ReplyGroup
	var replyItems []*ReplyItem
	OldDB.Where("user_id = ?", userId).Find(&replyGroups)
	OldDB.Where("user_id = ? and content like ?", userId, "%"+search+"%").Find(&replyItems)
	temp := make(map[string]*ReplyGroup)
	for _, replyGroup := range replyGroups {
		replyGroup.Items = make([]*ReplyItem, 0)
		temp[replyGroup.Id] = replyGroup
	}
	for _, replyItem := range replyItems {
		temp[replyItem.GroupId].Items = append(temp[replyItem.GroupId].Items, replyItem)
	}
	var newReplyGroups = make([]*ReplyGroup, 0)
	for _, replyGroup := range replyGroups {
		if len(replyGroup.Items) != 0 {
			newReplyGroups = append(newReplyGroups, replyGroup)
		}
	}
	return newReplyGroups
}

func GetReplyGroup(groupName string, kefuId string) Reply_group {
	var replyGroup Reply_group
	DB.Where(&Reply_group{}).Where(&Reply_group{GroupName: groupName, UserId: kefuId}).Find(&replyGroup)
	return replyGroup
}

func UpdateReplyGroup(groupName string, kefuId string, newGroupName string) {
	DB.Model(&Reply_group{}).Where(&Reply_group{GroupName: groupName, UserId: kefuId}).Update("group_name", newGroupName)
}
