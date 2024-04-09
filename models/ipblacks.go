package models

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Ipblack struct {
	ID       uint      `gorm:"primary_key" json:"id"`
	IP       string    `json:"ip"`
	KefuId   string    `json:"kefu_id"`
	CreateAt time.Time `json:"create_at"`
}

func CreateIpblack(ip string, kefuId string) uint {
	black := &Ipblack{
		IP:       ip,
		KefuId:   kefuId,
		CreateAt: time.Now(),
	}
	OldDB.Create(black)
	return black.ID
}
func DeleteIpblackByIp(ip string) {
	OldDB.Where("ip = ?", ip).Delete(Ipblack{})
}
func FindIp(ip string) Ipblack {
	var ipblack Ipblack
	OldDB.Where("ip = ?", ip).First(&ipblack)
	return ipblack
}
func FindIpsByKefuId(id string) []Ipblack {
	var ipblack []Ipblack
	OldDB.Where("kefu_id = ?", id).Find(&ipblack)
	return ipblack
}
func FindIps(query interface{}, args []interface{}, page uint, pagesize uint) []Ipblack {
	offset := (page - 1) * pagesize
	if offset < 0 {
		offset = 0
	}
	var ipblacks []Ipblack
	if query != nil {
		OldDB.Where(query, args...).Offset(offset).Limit(pagesize).Find(&ipblacks)
	} else {
		OldDB.Offset(offset).Limit(pagesize).Find(&ipblacks)
	}
	return ipblacks
}

// 查询条数
func CountIps(query interface{}, args []interface{}) uint {
	var count uint
	if query != nil {
		OldDB.Model(&Visitor{}).Where(query, args...).Count(&count)
	} else {
		OldDB.Model(&Visitor{}).Count(&count)
	}
	return count
}

func GetIpblack(ip string, kefuId string) (Ipblack, error) {
	var ipline Ipblack
	err := OldDB.Where(&Ipblack{IP: ip, KefuId: kefuId}).Find(&ipline).Error
	if err != nil {
		logrus.Error(err.Error())
	}
	return ipline, err
}
