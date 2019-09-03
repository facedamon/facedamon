package pkg

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func GetPage(c *gin.Context) int {
	res := 0
	p, err := strconv.ParseInt(c.Query("pageSize"), 10, 32)
	if nil != err {
		log.Fatalf("Fail to stronv pageSize to int")
	}
	if p > 0 {
		res = (int(p) - 1) * PageSize
	}
	return res
}
