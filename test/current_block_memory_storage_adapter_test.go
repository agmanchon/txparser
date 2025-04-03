package test

import (
	"context"
	"github.com/agmanchon/txparser/pkg/adapters/output/memorystorage"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

const (
	testCurrentBlock1 = "12345"
	testCurrentBlock2 = "6789"
)

func TestCurrentBlockMemoryStorageAdapterOk(t *testing.T) {
	ctx := context.TODO()

	testCurrentBlock1BigInt := new(big.Int)
	testCurrentBlock1BigInt, success := testCurrentBlock1BigInt.SetString(testCurrentBlock1, 10)
	require.True(t, success)
	require.NotNil(t, testCurrentBlock1BigInt)

	testCurrentBlock2BigInt := new(big.Int)
	testCurrentBlock2BigInt, success = testCurrentBlock2BigInt.SetString(testCurrentBlock2, 10)
	require.True(t, success)
	require.NotNil(t, testCurrentBlock2BigInt)

	options := memorystorage.DefaultCurrentBlockMemoryStorageAdapterOptions{}
	currentBlockMemoryStorageAdapter, err := memorystorage.ProvideDefaultCurrentBlockMemoryStorageAdapter(options)
	require.Nil(t, err)

	currentBlockReturned, err := currentBlockMemoryStorageAdapter.Add(ctx, *testCurrentBlock1BigInt)
	require.Nil(t, err)
	require.EqualValues(t, testCurrentBlock1, currentBlockReturned.String())
	currentBlockReturned, err = currentBlockMemoryStorageAdapter.Update(ctx, *testCurrentBlock2BigInt)
	require.Nil(t, err)
	require.EqualValues(t, testCurrentBlock2, currentBlockReturned.String())
	currentBlockReturned, err = currentBlockMemoryStorageAdapter.Get(ctx)
	require.Nil(t, err)
	require.EqualValues(t, testCurrentBlock2, currentBlockReturned.String())

}
