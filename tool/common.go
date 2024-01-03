package tool

import (
	"log"
	"os"
)

func InArray[T int | string](value T, array []T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func MakeDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatal("创建文件失败" + err.Error())
		}
	}
}
