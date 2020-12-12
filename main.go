package main

import "github.com/gin-gonic/gin"

func main() {
	eng := gin.New()
	eng.GET("/", func(context *gin.Context) {
		_, _ = context.Writer.WriteString("hello world")
	})
	_ = eng.Run()
}
