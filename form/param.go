package form

type PageParam struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
