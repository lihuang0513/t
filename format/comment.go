package format

type CommentResult struct {
	Cid      int    `json:"cid"`
	Room     int64  `json:"room"`
	Status   string `json:"status"`
	ParentId int64  `json:"parent_id"`
	Filename string `json:"filename"`
	Duration int64  `json:"duration"`
}
