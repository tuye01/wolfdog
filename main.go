package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"wolfdog/api"
	"wolfdog/internal/cmd"
)

func main() {
	cmd.InitDBEnv()
	cmd.InitRedisEnv()
	gin.SetMode(gin.DebugMode) //开发环境
	//gin.SetMode(gin.ReleaseMode) //线上环境
	r := gin.Default()
	api.RegisterRoutes(r)
	err := r.Run(":8080")
	if err != nil {
		fmt.Printf("启动失败!")
		return
	} // listen and serve on 0.0.0.0:8080
	fmt.Printf("成功启动!")

}
