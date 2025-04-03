package memorystorage

import (
	"context"
	"github.com/agmanchon/txparser/pkg/usecases/txparser"
	"math/big"
	"sync"
)

var _ txparser.CurrentBlockRepository = (*DefaultCurrentBlockMemoryStorageAdapter)(nil)

type DefaultCurrentBlockMemoryStorageAdapter struct {
	currentBlockNumber big.Int
	muCurrentBlock     sync.RWMutex
}

type DefaultCurrentBlockMemoryStorageAdapterOptions struct {
}

func ProvideDefaultCurrentBlockMemoryStorageAdapter(options DefaultCurrentBlockMemoryStorageAdapterOptions) (*DefaultCurrentBlockMemoryStorageAdapter, error) {

	return &DefaultCurrentBlockMemoryStorageAdapter{}, nil
}

func (d *DefaultCurrentBlockMemoryStorageAdapter) Add(ctx context.Context, data big.Int) (*big.Int, error) {
	d.muCurrentBlock.Lock()
	defer d.muCurrentBlock.Unlock()
	d.currentBlockNumber = data
	return &data, nil
}

func (d *DefaultCurrentBlockMemoryStorageAdapter) Get(ctx context.Context) (*big.Int, error) {
	d.muCurrentBlock.RLock()
	defer d.muCurrentBlock.RUnlock()
	return &d.currentBlockNumber, nil
}

func (d *DefaultCurrentBlockMemoryStorageAdapter) Update(ctx context.Context, data big.Int) (*big.Int, error) {
	d.muCurrentBlock.Lock()
	defer d.muCurrentBlock.Unlock()
	d.currentBlockNumber = data
	return &data, nil
}
