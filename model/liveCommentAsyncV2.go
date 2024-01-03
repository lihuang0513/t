package model

import (
	"comment-chuli/setup"
	"time"
)

type LiveCommentAsyncV2 struct {
	Id         int
	Cid        int
	IsHang     int
	CreateTime time.Time `gorm:"column:createtime"`
}

const lcTableName = "live_comment_async"

func (lc *LiveCommentAsyncV2) GetList(limit int) []LiveCommentAsyncV2 {
	var rows []LiveCommentAsyncV2
	result := setup.PlDb.Table(lcTableName).Where("is_hang = ?", 0).Limit(limit).Find(&rows)
	if result.RowsAffected == 0 {
		return nil
	}
	return rows
}

func (lc *LiveCommentAsyncV2) UpdateHangByFilename(idSlice []int) bool {
	result := setup.PlDb.Table(lcTableName).
		Where("cid IN (?)", idSlice).
		Update("is_hang", 1)

	if result.RowsAffected == 0 {
		return false
	}
	return true
}

func (lc *LiveCommentAsyncV2) DeleteHangByFilename(idSlice []int) bool {
	var rows LiveCommentAsyncV2
	result := setup.PlDb.Table(lcTableName).
		Where("cid IN (?)", idSlice).Delete(&rows)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}
