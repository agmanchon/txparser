package test

import (
	"context"
	"github.com/agmanchon/txparser/pkg/adapters/output/memorystorage"
	"github.com/agmanchon/txparser/pkg/usecases/txparser"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

const (
	testTransactionId1      = "testTransactionId1"
	testTransactionId2      = "testTransactionId2"
	testTransactionId3      = "testTransactionId3"
	testTransactionAddress1 = "testAddress1"
	testTransactionAddress2 = "testAddrres2"
)

func TestTransactionMemoryStorageAdapterOk(t *testing.T) {
	ctx := context.TODO()

	testTransaction1 := txparser.Transaction{
		ID:            testTransactionId1,
		BlockNumber:   big.Int{},
		Source:        testTransactionAddress1,
		Destination:   "",
		AmountInEther: big.Float{},
		FeesInEther:   big.Float{},
		CreationDate:  time.Time{},
	}

	testTransaction2 := txparser.Transaction{
		ID:            testTransactionId2,
		BlockNumber:   big.Int{},
		Source:        "",
		Destination:   testTransactionAddress2,
		AmountInEther: big.Float{},
		FeesInEther:   big.Float{},
		CreationDate:  time.Time{},
	}

	testTransaction3 := txparser.Transaction{
		ID:            testTransactionId3,
		BlockNumber:   big.Int{},
		Source:        testTransactionAddress1,
		Destination:   testTransactionAddress2,
		AmountInEther: big.Float{},
		FeesInEther:   big.Float{},
		CreationDate:  time.Time{},
	}

	options := memorystorage.DefaultTransactionMemoryStorageAdapterOptions{}
	transactionMemoryStorageAdapter, err := memorystorage.ProvideDefaultTransactionMemoryStorageAdapter(options)
	require.Nil(t, err)
	transactionReturned, err := transactionMemoryStorageAdapter.Add(ctx, testTransaction1)
	require.Nil(t, err)
	require.EqualValues(t, testTransaction1, *transactionReturned)
	transactionReturned, err = transactionMemoryStorageAdapter.Add(ctx, testTransaction2)
	require.Nil(t, err)
	require.EqualValues(t, testTransaction2, *transactionReturned)
	transactionReturned, err = transactionMemoryStorageAdapter.Add(ctx, testTransaction3)
	require.Nil(t, err)
	require.EqualValues(t, testTransaction3, *transactionReturned)
	transactionCollection, err := transactionMemoryStorageAdapter.All(ctx, txparser.TransactionListOptions{Address: testTransactionAddress1})
	require.Nil(t, err)
	require.Len(t, transactionCollection.Items, 2)
	require.Contains(t, transactionCollection.Items, testTransaction1)
	require.Contains(t, transactionCollection.Items, testTransaction3)

	transactionCollection, err = transactionMemoryStorageAdapter.All(ctx, txparser.TransactionListOptions{Address: testTransactionAddress2})
	require.Nil(t, err)
	require.Len(t, transactionCollection.Items, 2)
	require.Contains(t, transactionCollection.Items, testTransaction2)
	require.Contains(t, transactionCollection.Items, testTransaction3)
}
