package routers

import (
	"fmt"
	"net/http"

	"github.com/facedamon/meta-model/sql"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", "")
}

var mbwi = sql.NewMbwi()

type mbwiServices string

func NewMbwi() mbwiServices {
	return "mbsi"
}

func (mbwiServices) QueryByNum(c *gin.Context) {
	n := c.Param("num")
	m, err := mbwi.QueryByNum(n)

	if nil != err {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("sql.ModelBaseWordInfo.QueryByNum is failed with num='%s' '%v'", n, err))
		return
	}
	c.JSON(http.StatusOK, m)
}

func (mbwiServices) Count(c *gin.Context) {
	n, err := mbwi.Count()
	if nil != err {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("sql.ModelBaseWordInfo.Count is failed '%v'", err))
		return
	}
	c.JSON(http.StatusOK, n)
}

func (mbwiServices) QueryByStruct(c *gin.Context) {
	var json sql.ModelBaseWordInfo
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ms, err := mbwi.QueryByStruct(json)
	if nil != err {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("sql.ModelBaseWordInfo.QueryByStruct is failed '%v'", err))
		return
	}
	c.JSON(http.StatusOK, ms)
}

func (mbwiServices) UpdateByNum(c *gin.Context) {
	//n := c.Param("num")
	/*m, _ := mbwi.QueryByNum(n)
	id, err := mbwi.UpdateByNum(m)
	if nil != err {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("sql.ModelBaseWordInfo.UpdateByNum is failed '%v'", err))
		return
	}
	c.JSON(http.StatusOK, id)*/
}

func (mbwiServices) Insert(c *gin.Context) {
	var json sql.ModelBaseWordInfo
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	n, err := mbwi.Insert(json)
	if nil != err {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("sql.ModelBaseWordInfo.Insert is failed '%v'", err))
		return
	}
	c.JSON(http.StatusOK, n)
}
