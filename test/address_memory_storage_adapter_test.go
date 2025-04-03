package test

import (
	"context"
	"github.com/agmanchon/txparser/pkg/adapters/output/memorystorage"
	"github.com/agmanchon/txparser/pkg/usecases/txparser"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testAddress1 = "testAddress1"
	testAddress2 = "testAddrres2"
)

func TestAddressMemoryStorageAdapterOk(t *testing.T) {
	ctx := context.TODO()

	options := memorystorage.DefaultAddressMemoryStorageAdapterOptions{}
	addressMemoryStorageAdapter, err := memorystorage.ProvideDefaultAddressMemoryStorageAdapter(options)
	require.Nil(t, err)
	addressReturned, err := addressMemoryStorageAdapter.Add(ctx, testAddress1)
	require.Nil(t, err)
	require.EqualValues(t, testAddress1, *addressReturned)
	addressReturned, err = addressMemoryStorageAdapter.Add(ctx, testAddress2)
	require.Nil(t, err)
	require.EqualValues(t, testAddress2, *addressReturned)
	addressCollection, err := addressMemoryStorageAdapter.All(ctx, txparser.AddressListOptions{})
	require.Nil(t, err)
	require.Len(t, addressCollection.Items, 2)
	require.Contains(t, addressCollection.Items, testAddress1)
	require.Contains(t, addressCollection.Items, testAddress2)
}
