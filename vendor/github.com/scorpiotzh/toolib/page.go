package toolib

type Pagination struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func (p Pagination) GetLimit() int {
	if p.Size < 1 || p.Size > 200 {
		return 200
	}
	return p.Size
}

func (p Pagination) GetOffset() int {
	page := p.Page
	if p.Page < 1 {
		page = 1
	}
	size := p.GetLimit()
	return (page - 1) * size
}
