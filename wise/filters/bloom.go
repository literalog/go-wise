package filters

import (
	"context"
	"sync"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
)

type Bloom[T any] struct {
	innerBloom *bloom.BloomFilter
	builder    Builder[T]
	lock       sync.Mutex
}

type Builder[T any] interface {
	Build([]T) [][]byte
	Fetch(context.Context) ([]T, error)
}

func NewBloomFilter[T any](ctx context.Context, builder Builder[T]) (*Bloom[T], error) {
	data, err := builder.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	inBF := bloom.NewWithEstimates(uint(len(data)), 0.01)

	bf := &Bloom[T]{
		innerBloom: inBF,
		builder:    builder,
		lock:       sync.Mutex{},
	}

	bf.Add(builder.Build(data))
	bf.subscribe(ctx)

	return bf, nil
}

func (b *Bloom[T]) Add(keys [][]byte) {
	for _, k := range keys {
		b.innerBloom.Add(k)
	}
}

func (b *Bloom[T]) Has(key string) bool {
	return b.innerBloom.Test([]byte(key))
}

func (b *Bloom[T]) Refresh(ctx context.Context) error {
	data, err := b.builder.Fetch(ctx)
	if err != nil {
		return err
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	b.innerBloom.ClearAll()
	b.Add(b.builder.Build(data))

	return nil
}

func (b *Bloom[T]) subscribe(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			b.Refresh(ctx)
		}
	}()
}
