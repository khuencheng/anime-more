package backend

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.GET("/", Recommend)
	return router

}

func StartServer(debugMode bool, port string) {
	r := NewRouter()
	if debugMode {
		pprof.Register(r)
	}
	r.Run(fmt.Sprintf("127.0.0.1:%s", port))
}
