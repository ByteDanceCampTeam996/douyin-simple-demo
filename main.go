package main

import (
	"github.com/ByteDanceCampTeam996/douyin-simple-demo/controller"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	// 访问日志记录输出
	// 日志记录到文件
	//f, _ := os.Create("./temp/logs/douyin.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	/*gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}*/

	// 定义Restful API接口路由
	initRouter(r)

	// 初始化数据库连接
	controller.ConnectDB()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
