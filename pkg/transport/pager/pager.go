package pager

type Page struct {
	PageIndex int `json:"page_index" form:"page_index,default=1" binding:"required" url:"page_index"`
	PageSize  int `json:"page_size" form:"page_size,default=20" binding:"required" url:"page_size"`
}
