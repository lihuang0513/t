package model

import "comment-chuli/setup"

type LiveFloorComment struct {
	Filename string
	FloorNum int
}

func (lfc *LiveFloorComment) GetByCidSlice(cidSlice []int) []LiveFloorComment {
	var rows []LiveFloorComment

	result := setup.PlDb.Table("live_floor_comment").Where("comment_id in ?", cidSlice).Find(&rows)
	if result.RowsAffected == 0 {
		return nil
	}
	return rows
}
