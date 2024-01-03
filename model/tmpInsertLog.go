package model

import (
	"comment-chuli/setup"
	"time"
)

type TmpInsertLog struct {
	Id         int
	Position   string
	Message    string
	Type       int
	MUid       int
	OpName     string
	CreateTime time.Time `gorm:"column:createtime"`
}

func (til *TmpInsertLog) SaveData(position string, message string, uid int, opName string, typeValue int) int {
	newTmpInsertLog := TmpInsertLog{
		Id:         0,
		Position:   position,
		Message:    message,
		Type:       typeValue,
		MUid:       uid,
		OpName:     opName,
		CreateTime: time.Now(),
	}
	result := setup.PlDb.Table("tmp_insert_log").Create(&newTmpInsertLog)
	if result.Error != nil {
		setup.Logger.Printf("err,TmpInsertLog-SaveData%s", result.Error)
		return 0
	} else {
		return newTmpInsertLog.Id
	}
}
