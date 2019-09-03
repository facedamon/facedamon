package routers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Login(c *gin.Context) {
	log.Println("------------Login---------------")
	c.HTML(http.StatusOK, "login.html", "")
}
