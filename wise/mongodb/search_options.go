package mongodb

import "go.mongodb.org/mongo-driver/mongo/options"

type Method string

const (
	CURSOR  Method = "cursor"
	DEFAULT Method = "default"
)

func NewMethod(method string) Method {
	switch method {
	case "cursor":
		return CURSOR
	default:
		return DEFAULT
	}
}

type SearchOptions func(*searchOptions)

type searchOptions struct {
	Method   Method `json:"method"`
	PageSize int    `json:"pageSize"`
	Page     int    `json:"page"`
	Sort     map[string]int
}

func NewSearchOptions(fnopts ...SearchOptions) *searchOptions {
	opts := &searchOptions{}

	for _, fn := range fnopts {
		fn(opts)
	}

	return opts
}

func (w searchOptions) ToFindOptions(maxPageSize int) *options.FindOptions {
	opts := options.Find()

	if w.PageSize > 0 {
		if w.PageSize > maxPageSize {
			w.PageSize = maxPageSize
		}

		if w.Page < 0 {
			w.Page = 0
		}

		skip := (w.Page) * w.PageSize

		opts.SetSkip(int64(skip))
		opts.SetLimit(int64(w.PageSize))
	}

	if w.Sort != nil || len(w.Sort) != 0 {
		opts.SetSort(w.Sort)
	}

	return opts
}

func WithPageMethod(method string) SearchOptions {
	return func(opts *searchOptions) {
		opts.Method = NewMethod(method)
	}
}

func WithPageSize(pageSize int) SearchOptions {
	return func(opts *searchOptions) {
		if pageSize > 0 {
			opts.PageSize = pageSize
		}
	}
}

func WithPage(page int) SearchOptions {
	return func(opts *searchOptions) {
		if page >= 0 {
			opts.Page = page
		}
	}
}

func WithSort(sort map[string]int) SearchOptions {
	return func(opts *searchOptions) {
		if sort != nil {
			opts.Sort = sort
		}
	}
}
