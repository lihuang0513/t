package setup

import (
	"comment-chuli/format"
	"github.com/go-ini/ini"
	"log"
)

var Config format.Config

func InitConfig() {
	// 加载 .ini 配置
	loadIni("./conf/config.ini")
}

// load 加载配置项
func loadIni(configPath string) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		//失败
		log.Fatalf("配置文件加载失败：%q", err.Error())
	}
	err = cfg.MapTo(&Config)
	if err != nil {
		//赋值失败
		log.Fatalf("配置文件赋值失败：%q", err.Error())
	}

}
