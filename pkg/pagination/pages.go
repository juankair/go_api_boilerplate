package pagination

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	DefaultPageSize = 8
	MaxPageSize     = 1000
	PageVar         = "page"
	PageSizeVar     = "per_page"
)

type Pages struct {
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	PageCount  int         `json:"page_count"`
	TotalCount int         `json:"total_count"`
	Items      interface{} `json:"items"`
}

func New(page, perPage, total int) *Pages {
	if perPage <= 0 {
		perPage = DefaultPageSize
	}
	if perPage > MaxPageSize {
		perPage = MaxPageSize
	}
	pageCount := -1
	if total >= 0 {
		pageCount = (total + perPage - 1) / perPage
		if page > pageCount {
			page = pageCount
		}
	}
	if page < 1 {
		page = 1
	}

	return &Pages{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		PageCount:  pageCount,
	}
}

func NewFromRequest(req *http.Request, count int) *Pages {
	page := parseInt(req.URL.Query().Get(PageVar), 1)
	perPage := parseInt(req.URL.Query().Get(PageSizeVar), DefaultPageSize)
	return New(page, perPage, count)
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}

func (p *Pages) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func (p *Pages) Limit() int {
	return p.PerPage
}

func (p *Pages) BuildLinkHeader(baseURL string, defaultPerPage int) string {
	links := p.BuildLinks(baseURL, defaultPerPage)
	header := ""
	if links[0] != "" {
		header += fmt.Sprintf("<%v>; rel=\"first\", ", links[0])
		header += fmt.Sprintf("<%v>; rel=\"prev\"", links[1])
	}
	if links[2] != "" {
		if header != "" {
			header += ", "
		}
		header += fmt.Sprintf("<%v>; rel=\"next\"", links[2])
		if links[3] != "" {
			header += fmt.Sprintf(", <%v>; rel=\"last\"", links[3])
		}
	}
	return header
}

func (p *Pages) BuildLinks(baseURL string, defaultPerPage int) [4]string {
	var links [4]string
	pageCount := p.PageCount
	page := p.Page
	if pageCount >= 0 && page > pageCount {
		page = pageCount
	}
	if strings.Contains(baseURL, "?") {
		baseURL += "&"
	} else {
		baseURL += "?"
	}
	if page > 1 {
		links[0] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, 1)
		links[1] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page-1)
	}
	if pageCount >= 0 && page < pageCount {
		links[2] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page+1)
		links[3] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, pageCount)
	} else if pageCount < 0 {
		links[2] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page+1)
	}
	if perPage := p.PerPage; perPage != defaultPerPage {
		for i := 0; i < 4; i++ {
			if links[i] != "" {
				links[i] += fmt.Sprintf("&%v=%v", PageSizeVar, perPage)
			}
		}
	}

	return links
}
