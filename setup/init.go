package setup

import (
	"sync"
)

var once sync.Once

func Init() {
	once.Do(func() {
		InitConfig()
		InitLogger()
		InitPlDb()
	})
}
