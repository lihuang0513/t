package setup

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var PlDb *gorm.DB

func InitPlDb() {

	var err error

	plDbConfig := Config.PlDb
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&timeout=3s",
		plDbConfig.User, plDbConfig.Pwd, plDbConfig.Host, plDbConfig.Port, plDbConfig.Name, plDbConfig.Charset)
	PlDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // 使用 log 包输出日志
			logger.Config{
				LogLevel:                  logger.Silent, // 设置日志级别为 Info，这里你可以根据需要调整级别
				Colorful:                  true,          // 是否启用彩色打印
				IgnoreRecordNotFoundError: true,
			},
		),
	})
	if err != nil {
		log.Fatalf("err,plDb connection failed：%q", err.Error())
	}

	sqlDB, err := PlDb.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	if err != nil {
		log.Fatalf("err,plDb get db failed：%q", err.Error())
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("err,plDb ping failed：%q", err.Error())
	}

	Config.PlDb.Enable = true
}
