package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (p PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qr := r.URL.Query()

	limit := qr.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return p, err
		}
		p.Limit = l
	}

	offset := qr.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return p, err
		}
		p.Offset = o
	}

	sort := qr.Get("sort")
	if sort != "" {
		p.Sort = sort
	}

	return p, nil
}
