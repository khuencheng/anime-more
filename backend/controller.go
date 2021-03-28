package backend

import "C"
import (
	"github.com/gin-gonic/gin"
	"log"
)

func Recommend(c *gin.Context) {
	keyword, _ := c.GetQuery("anime")
	log.Println(keyword)
	items := GetRecommendService(keyword)

	c.JSON(200, Response{
		Success: true,
		Data:    items},
	)

}
