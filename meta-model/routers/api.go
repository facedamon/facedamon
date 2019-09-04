package routers

import (
	"log"
	"net/http"

	"github.com/facedamon/meta-model/sql"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	log.Println("------------Login---------------")
	c.HTML(http.StatusOK, "login.html", "")
}

func ModelQueryAll(c *gin.Context) {
	m := sql.NewModelBaseWordInfo()
	m.QueryByNum("0000000001")
	c.JSON(http.StatusOK, m)
}
