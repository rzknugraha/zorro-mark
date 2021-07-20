package helpers

//Page data sturct
type Page struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalData    int `json:"total_data"`
	NextPage     int `json:"next_page"`
	PreviousPage int `json:"previous_page"`
}

//Paginate data sturct
type Paginate struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Page    Page        `json:"page"`
	Data    interface{} `json:"data"`
}

//Pages data sturct
type Pages struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	PageCount  int `json:"page_count"`
	TotalCount int `json:"total_count"`
}

//PageReq data sturct
type PageReq struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Page   int `json:"page"`
}

//WrapPaginate wraping paginate
func WrapPaginate(pages *Pages, data interface{}) *Paginate {
	return &Paginate{
		Code:    2200,
		Message: "Sucess",
		Error:   "",
		Page: Page{
			CurrentPage:  pages.Page,
			TotalPages:   pages.PageCount,
			TotalData:    pages.TotalCount,
			NextPage:     pages.NextPage(),
			PreviousPage: pages.PrevPage(),
		},
		Data: data,
	}
}

//NewPaginate is s
func NewPaginate(page, perPage, total int) *Pages {
	pageCount := -1
	if total >= 0 {
		pageCount = (total + perPage - 1) / perPage
		if page > pageCount {
			page = pageCount
		}
	}

	return &Pages{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		PageCount:  pageCount,
	}
}

//NextPage is ss
func (p *Pages) NextPage() int {
	pageCount := p.PageCount
	page := p.Page

	if pageCount >= 0 && page >= pageCount {
		page = pageCount
	} else {
		page = page + 1
	}

	return page
}

//PrevPage is a
func (p *Pages) PrevPage() int {
	page := p.Page
	if page > 1 {
		page = page - 1
	}

	return page
}
