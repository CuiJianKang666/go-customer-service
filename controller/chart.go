package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/taoshihan1991/imaptool/models"
	"github.com/taoshihan1991/imaptool/tools"
	"time"
)

func GetChartStatistic(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")

	dayNumMap := make(map[string]string)
	var result []models.EveryDayNum
	if kefuName == "root" {
		result = models.CountRootVisitorsEveryDay()
	} else {
		result = models.CountVisitorsEveryDay(kefuName.(string))
	}
	fmt.Println(result)

	for _, item := range result {
		dayNumMap[item.Day] = tools.Int2Str(item.Num)
	}

	nowTime := time.Now()
	list := make([]map[string]string, 0)
	for i := 0; i > -7; i-- {
		getTime := nowTime.AddDate(0, 0, i)   //年，月，日   获取一天前的时间
		resTime := getTime.Format("06-01-02") //获取的时间的格式
		tmp := make(map[string]string)
		tmp["day"] = resTime
		if dayNumMap[resTime] == "" {
			tmp["num"] = "0"
		} else {
			tmp["num"] = dayNumMap[resTime]
		}
		fmt.Println(tmp)
		list = append(list, tmp)
	}

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": list,
	})

}
