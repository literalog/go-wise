package mongodb

import "github.com/literalog/go-wise/wise/filters"

type RepositoryOptions func(*repositoryOptions)

type repositoryOptions struct {
	bloomFilter     filters.Bloom[string]
	maxPageSize     int
	defaultPageSize int
}

func NewRepositoryOptions(fnopts ...RepositoryOptions) *repositoryOptions {
	opts := &repositoryOptions{}

	for _, fn := range fnopts {
		fn(opts)
	}

	return opts
}

func WithBloomFilter(bf filters.Bloom[string]) RepositoryOptions {
	return func(opts *repositoryOptions) {
		opts.bloomFilter = bf
	}
}

func WithMaxPageSize(size int) RepositoryOptions {
	return func(opts *repositoryOptions) {
		opts.maxPageSize = size
	}
}

func WithDefaultPageSize(size int) RepositoryOptions {
	return func(opts *repositoryOptions) {
		opts.defaultPageSize = size
	}
}
