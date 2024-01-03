package setup

import (
	"comment-chuli/tool"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"log"
	"os"
	"time"
)

// Logger 日志客户端
var Logger *log.Logger

func InitLogger() {

	if Config.Server.Debug {
		Logger = log.New(os.Stdout, "", log.LstdFlags)
		return
	}

	// 判断文件夹是否存在没有则创建
	tool.MakeDir(Config.Server.LogDir)
	// 获取日志文件句柄
	logFile, err := rotatelogs.New(
		Config.Server.LogDir+"%Y-%m-%d-comment.log",
		//设置一天产生一个日志文件
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		log.Println("创建文件失败" + err.Error())
	}
	// 设置存储位置
	Logger = log.New(logFile, "", log.Lshortfile|log.LstdFlags)
}
