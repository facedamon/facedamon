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
		api.GET("/ModelBaseWorldInfo", ModelQueryAll)
	}

	return r
}
