package util

func PageAndPageSize(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 15
	}
	return page, pageSize
}

func SkipAndLimit(page, pageSize int) (int, int) {
	page, pageSize = PageAndPageSize(page, pageSize)
	return (page - 1) * pageSize, pageSize
}
