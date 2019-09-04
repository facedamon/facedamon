package routers

import (
	"github.com/facedamon/meta-model/pkg"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(pkg.RunMode)
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*.html")

	api := r.Group("/")
	{
		api.GET("/", Login)
	}

	m := NewMbwi()
	modelBaseWorldInfo := r.Group("/ModelBaseWordInfo")
	{
		modelBaseWorldInfo.POST("/queryByNum/:num", m.QueryByNum)
		modelBaseWorldInfo.POST("/count", m.Count)
		modelBaseWorldInfo.POST("/list", m.QueryByStruct)
	}

	return r
}
