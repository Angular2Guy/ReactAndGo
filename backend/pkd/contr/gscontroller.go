package contr

import (
	"angular-and-go/pkd/gasstation"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.POST("/posts", postsCreate)
	r.GET("/gasprice/:id", postsIndex)
	r.GET("/gasstation/:id", postsShow)
	r.PUT("/posts/:id", postsUpdate)
	r.DELETE("posts/:id", postsDelete)
	r.Run() // listen and serve on 0.0.0.0:8080
}

func postsCreate(c *gin.Context) {

}

func postsIndex(c *gin.Context) {
	gasstationId := c.Params.ByName("id")
	fmt.Println(gasstationId)
	gsEntity := gasstation.FindPricesByStid(gasstationId)
	c.JSON(200, gsEntity)
}

func postsShow(c *gin.Context) {
	gasstationId := c.Params.ByName("id")
	fmt.Println(gasstationId)
	gsEntity := gasstation.FindById(gasstationId)
	c.JSON(200, gsEntity)
}

func postsUpdate(c *gin.Context) {

}

func postsDelete(c *gin.Context) {

}
