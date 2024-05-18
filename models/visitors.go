package models

import (
	"time"
)

type Visitor struct {
	Model
	Name      string `json:"name"`
	Avator    string `json:"avator"`
	SourceIp  string `json:"source_ip"`
	ToId      string `json:"to_id"`
	VisitorId string `json:"visitor_id"`
	Status    uint   `json:"status"`
	Refer     string `json:"refer"`
	City      string `json:"city"`
	ClientIp  string `json:"client_ip"`
	Extra     string `json:"extra"`
}

func CreateVisitor(name, avator, sourceIp, toId, visitorId, refer, city, clientIp, extra string) {
	v := &Visitor{
		Name:      name,
		Avator:    avator,
		SourceIp:  sourceIp,
		ToId:      toId,
		VisitorId: visitorId,
		Status:    1,
		Refer:     refer,
		City:      city,
		ClientIp:  clientIp,
		Extra:     extra,
	}
	v.UpdatedAt = time.Now()
	OldDB.Create(v)
}
func FindVisitorByVistorId(visitorId string) Visitor {
	var v Visitor
	OldDB.Where("visitor_id = ?", visitorId).First(&v)
	return v
}
func FindVisitors(page uint, pagesize uint) []Visitor {
	offset := (page - 1) * pagesize
	if offset < 0 {
		offset = 0
	}
	var visitors []Visitor
	OldDB.Offset(offset).Limit(pagesize).Order("status desc, updated_at desc").Find(&visitors)
	return visitors
}
func UpdateToKefuId(to_id string, target_id string) {
	DB.Where(&Visitor{ToId: to_id}).Update("to_id", target_id)
}
func FindVisitorsByKefuId(page uint, pagesize uint, kefuId string) []Visitor {
	offset := (page - 1) * pagesize
	if offset <= 0 {
		offset = 0
	}
	var visitors []Visitor
	//sql := fmt.Sprintf("select * from visitor where id>=(select id from visitor where  to_id='%s' order by updated_at desc limit %d,1) and to_id='%s' order by updated_at desc limit %d ", kefuId, offset, kefuId, pagesize)
	//OldDB.Raw(sql).Scan(&visitors)
	OldDB.Where("to_id=?", kefuId).Offset(offset).Limit(pagesize).Order("updated_at desc").Find(&visitors)
	return visitors
}
func FindVisitorsOnline() []Visitor {
	var visitors []Visitor
	OldDB.Where("status = ?", 1).Find(&visitors)
	return visitors
}
func UpdateVisitorStatus(visitorId string, status uint) {
	visitor := Visitor{}
	OldDB.Model(&visitor).Where("visitor_id = ?", visitorId).Update("status", status)
}
func UpdateVisitor(name, avator, visitorId string, status uint, clientIp string, sourceIp string, refer, extra string, toId string) {
	visitor := &Visitor{
		Status:   status,
		ClientIp: clientIp,
		SourceIp: sourceIp,
		Refer:    refer,
		Extra:    extra,
		Name:     name,
		Avator:   avator,
		ToId:     toId,
	}
	visitor.UpdatedAt = time.Now()
	OldDB.Model(visitor).Where("visitor_id = ?", visitorId).Update(visitor)
}
func UpdateVisitorKefu(visitorId string, kefuId string) {
	visitor := Visitor{}
	OldDB.Model(&visitor).Where("visitor_id = ?", visitorId).Update("to_id", kefuId)
}

// 查询条数
func CountVisitors() uint {
	var count uint
	OldDB.Model(&Visitor{}).Count(&count)
	return count
}

// 查询条数
func CountVisitorsByKefuId(kefuId string) uint {
	var count uint
	OldDB.Model(&Visitor{}).Where("to_id=?", kefuId).Count(&count)
	return count
}

// 查询每天条数
type EveryDayNum struct {
	Day string `json:"day"`
	Num int64  `json:"num"`
}

func CountVisitorsEveryDay(toId string) []EveryDayNum {
	var results []EveryDayNum
	OldDB.Raw("select DATE_FORMAT(updated_at,'%y-%m-%d') as day ,"+
		"count(*) as num from visitor where to_id=? group by day order by day desc limit 7",
		toId).Scan(&results)
	return results
}

func CountRootVisitorsEveryDay() []EveryDayNum {
	var results []EveryDayNum
	OldDB.Raw("select DATE_FORMAT(updated_at,'%y-%m-%d') as day ," +
		"count(*) as num from visitor group by day order by day desc limit 7").Scan(&results)
	return results
}
