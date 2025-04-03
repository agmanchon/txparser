package memorystorage

import (
	"context"
	"github.com/agmanchon/txparser/pkg/usecases/txparser"
	"sync"
)

var _ txparser.AddressRepository = (*DefaultAddressMemoryStorageAdapter)(nil)

type DefaultAddressMemoryStorageAdapter struct {
	Addresses    []string
	muAddressses sync.RWMutex
}

type DefaultAddressMemoryStorageAdapterOptions struct {
}

func ProvideDefaultAddressMemoryStorageAdapter(options DefaultAddressMemoryStorageAdapterOptions) (*DefaultAddressMemoryStorageAdapter, error) {

	return &DefaultAddressMemoryStorageAdapter{}, nil
}

func (d *DefaultAddressMemoryStorageAdapter) Add(ctx context.Context, data string) (*string, error) {
	d.muAddressses.Lock()
	defer d.muAddressses.Unlock()
	d.Addresses = append(d.Addresses, data)
	return &data, nil
}

func (d *DefaultAddressMemoryStorageAdapter) All(ctx context.Context, filter txparser.AddressListOptions) (*txparser.AddressCollection, error) {
	d.muAddressses.RLock()
	defer d.muAddressses.RUnlock()
	return &txparser.AddressCollection{
		Items: d.Addresses,
	}, nil
}
