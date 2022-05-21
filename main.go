package main

import (
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	initRouter(r)

	controller.InitDb()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
