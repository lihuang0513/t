package main

import (
	"comment-chuli/router"
	"comment-chuli/setup"
	"fmt"
	"github.com/fvbock/endless"
	"log"
	"time"
)

func init() {
	// 时区初始化
	loc := time.FixedZone("UTC", 8*3600)
	time.Local = loc

	// 初始化配置文件、数据库链接、redis链接
	setup.Init()

}

func main() {

	// 设置路由
	r := router.SetRouters()

	// 4，配置(使用endless进行重启服务)
	addr := fmt.Sprintf("%s:%d", setup.Config.Server.Host, setup.Config.Server.Port)
	s := endless.NewServer(addr, r)

	err := s.ListenAndServe()
	if err != nil {
		log.Printf("server err: %v", err)
	}
}
