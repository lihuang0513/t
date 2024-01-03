package router

import (
	"comment-chuli/controller"
	"comment-chuli/setup"
	"comment-chuli/tool"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io/ioutil"
	"time"
)

func SetRouters() *gin.Engine {
	var r *gin.Engine

	if setup.Config.Server.Debug {
		r = gin.Default()
	} else {
		r = ReleaseRouter()
	}

	//pprof.Register(r)

	recommendGroup := r.Group("/live")
	{
		var liveController controller.LiveStartController
		recommendGroup.GET("/start", liveController.Start)
	}

	return r
}

// ReleaseRouter 生产模式使用官方建议设置为 release 模式
func ReleaseRouter() *gin.Engine {
	// 切换到生产模式
	gin.SetMode(gin.ReleaseMode)
	// 禁用 gin 输出接口访问日志
	gin.DefaultWriter = ioutil.Discard

	tool.MakeDir(setup.Config.Server.LogDir)
	// 记录到文件。
	f, _ := rotatelogs.New(
		setup.Config.Server.LogDir+"%Y-%m-%d-error_log.log",
		//设置一天产生一个日志文件
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	engine := gin.New()
	// 应用崩溃恢复中间件
	engine.Use(gin.Recovery())
	// 应用崩溃写入日志
	engine.Use(gin.RecoveryWithWriter(f))

	return engine
}
