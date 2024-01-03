package model

import (
	"comment-chuli/setup"
	"time"
)

type LiveCacheQueue struct {
	Id         int
	Filename   string
	Page       int
	CreateTime time.Time `gorm:"column:createtime"`
}

func (lcq *LiveCacheQueue) BatchSaveData(mapLiveQueue []map[string]interface{}) int {
	// 批量写入写源队列
	result := setup.PlDb.Table("live_cache_queue").Create(mapLiveQueue)
	if result.RowsAffected == 0 {
		setup.Logger.Println("err,LiveCacheQueue-BatchSaveData写入失败")
		return 0
	}
	return 1
}

func (lcq *LiveCacheQueue) GetByFPage(filename string, page int) *LiveCacheQueue {
	var row LiveCacheQueue
	result := setup.PlDb.Table("live_cache_queue").Where("filename = ? AND page = ?", filename, page).First(&row)
	if result.RowsAffected == 0 {
		return nil
	}
	return &row
}
