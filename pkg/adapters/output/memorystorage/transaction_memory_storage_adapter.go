package memorystorage

import (
	"context"
	"github.com/agmanchon/txparser/pkg/usecases/txparser"
	"sync"
)

var _ txparser.TransactionRepository = (*DefaultTransactionMemoryStorageAdapter)(nil)

type DefaultTransactionMemoryStorageAdapter struct {
	transactions   map[string][]txparser.Transaction
	muTransactions sync.RWMutex
}

type DefaultTransactionMemoryStorageAdapterOptions struct {
}

func ProvideDefaultTransactionMemoryStorageAdapter(options DefaultTransactionMemoryStorageAdapterOptions) (*DefaultTransactionMemoryStorageAdapter, error) {

	transactions := make(map[string][]txparser.Transaction, 0)
	return &DefaultTransactionMemoryStorageAdapter{transactions: transactions}, nil
}

func (d *DefaultTransactionMemoryStorageAdapter) Add(ctx context.Context, data txparser.Transaction) (*txparser.Transaction, error) {
	d.muTransactions.Lock()
	defer d.muTransactions.Unlock()
	if _, exists := d.transactions[data.Source]; !exists {
		d.transactions[data.Source] = make([]txparser.Transaction, 0)
	}

	if _, exists := d.transactions[data.Destination]; !exists {
		d.transactions[data.Destination] = make([]txparser.Transaction, 0)
	}

	d.transactions[data.Source] = append(d.transactions[data.Source], data)
	d.transactions[data.Destination] = append(d.transactions[data.Destination], data)

	return &data, nil
}

func (d *DefaultTransactionMemoryStorageAdapter) All(ctx context.Context, filter txparser.TransactionListOptions) (*txparser.TransactionCollection, error) {
	d.muTransactions.RLock()
	defer d.muTransactions.RUnlock()
	if _, exists := d.transactions[filter.Address]; !exists {
		return &txparser.TransactionCollection{}, nil
	}
	return &txparser.TransactionCollection{Items: d.transactions[filter.Address]}, nil
}
