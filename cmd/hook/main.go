package main

import (
	"flag"
	"fmt"

	"github.com/chenlujjj/hook/gitlab"
	"github.com/chenlujjj/hook/weixin"
	"github.com/gin-gonic/gin"
)

func main() {
	key := flag.String("key", "", "wechat bot key")
	flag.Parse()
	fmt.Println(key)
	wc := weixin.NewWechatClient(*key)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/gitlab/mr", gitlab.NewMRHandler(wc))
	r.Run()
}
