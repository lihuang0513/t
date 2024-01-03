package model

import (
	"comment-chuli/setup"
	"gorm.io/gorm"
)

type LiveCommentNum struct {
	Id         int
	Filename   string
	Num        int
	RootNormal int
}

func (lcn *LiveCommentNum) OpCommentNum(filename string, normal int, rootNormal int) int {

	if lcn.getByFilename(filename) == nil {

		newCommentNum := LiveCommentNum{
			Id:         0,
			Filename:   filename,
			Num:        normal,
			RootNormal: rootNormal,
		}
		result := setup.PlDb.Table("live_comment_num").Create(&newCommentNum)
		if result.RowsAffected == 0 {
			setup.Logger.Printf("err,LiveCommentNum-SaveData%s", result.Error)
			return 0
		} else {
			return newCommentNum.Id
		}
	} else {
		result := setup.PlDb.Table("live_comment_num").Where("filename = ?", filename).
			Updates(map[string]interface{}{
				"num":         gorm.Expr("num + ?", normal),
				"root_normal": gorm.Expr("root_normal + ?", rootNormal),
			})
		if result.RowsAffected == 0 {
			setup.Logger.Printf("err,LiveCommentNum-UpdateCommentNum,%s更新失败", filename)
			return 0
		}
		return 1
	}

}

func (lcn *LiveCommentNum) getByFilename(filename string) *LiveCommentNum {
	var row LiveCommentNum
	result := setup.PlDb.Table("live_comment_num").Where("filename = ?", filename).First(&row)
	if result.RowsAffected == 0 {
		return nil
	}
	return &row
}
