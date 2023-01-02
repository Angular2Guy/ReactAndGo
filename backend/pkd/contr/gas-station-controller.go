package contr

import (
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.POST("/posts", postsCreate)
	r.GET("/posts", postsIndex)
	r.GET("/posts/:id", postsShow)
	r.PUT("/posts/:id", postsUpdate)
	r.DELETE("posts/:id", postsDelete)
	r.Run() // listen and serve on 0.0.0.0:8080
}

func postsCreate(c *gin.Context) {

}

func postsIndex(c *gin.Context) {

}

func postsShow(c *gin.Context) {

}

func postsUpdate(c *gin.Context) {

}

func postsDelete(c *gin.Context) {

}
