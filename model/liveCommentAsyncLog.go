package model

import (
	"comment-chuli/setup"
	"time"
)

type LiveCommentAsyncLog struct {
	Id         int
	Type       string
	Content    string
	Ms         int64
	Status     string
	Cids       string
	CreateTime time.Time `gorm:"column:createtime"`
}

const asyncLogTableName = "live_comment_async_log"

func (lca *LiveCommentAsyncLog) SaveData(typeValue string, content string, ms int64, statusValue string, cids string) int {
	newLiveCommentAsyncLog := LiveCommentAsyncLog{
		Id:         0,
		Type:       typeValue,
		Content:    content,
		Ms:         ms,
		Status:     statusValue,
		Cids:       cids,
		CreateTime: time.Now(),
	}
	result := setup.PlDb.Table(asyncLogTableName).Create(&newLiveCommentAsyncLog)
	if result.RowsAffected == 0 {
		setup.Logger.Printf("err,LiveCommentAsyncLog-SaveData%s", result.Error)
		return 0
	} else {
		return newLiveCommentAsyncLog.Id
	}
}
