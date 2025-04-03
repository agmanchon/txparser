package test

import (
	"context"
	"github.com/agmanchon/txparser/pkg/adapters/output/ethereumrpc"
	"github.com/agmanchon/txparser/pkg/usecases/txparser"
	"github.com/stretchr/testify/require"
	"testing"
)

var ethEndpoint = "https://ethereum-rpc.publicnode.com"

func TestEthereumRpcAdapterOk(t *testing.T) {
	ctx := context.TODO()

	options := ethereumrpc.DefaultEthereumRpcAdapterOptions{}
	options.Host = ethEndpoint
	ethAdapter, err := ethereumrpc.ProvideDefaultEthereumRpcAdapter(options)
	require.Nil(t, err)

	getMostRecentBlockNumberInput := txparser.GetMostRecentBlockNumberAdapterInput{}
	currentBlockNumber, err := ethAdapter.GetMostRecentBlockNumberAdapter(ctx, getMostRecentBlockNumberInput)
	require.Nil(t, err)
	require.NotNil(t, currentBlockNumber)

	getBlockByNumberInput := txparser.GetBlockByNumberAdapterInput{}
	getBlockByNumberInput.BlockNumber = currentBlockNumber.BlockNumber
	block, err := ethAdapter.GetBlockByNumberAdapter(ctx, getBlockByNumberInput)
	require.Nil(t, err)
	require.NotNil(t, block)

	for _, transactionHash := range block.TransactionHashes {
		getTransactionByHashInput := txparser.GetTransactionAdapterInput{}
		getTransactionByHashInput.ID = transactionHash
		transaction, err := ethAdapter.GetTransactionAdapter(ctx, getTransactionByHashInput)
		require.Nil(t, err)
		require.NotNil(t, transaction)
		//time.Sleep(100 * time.Millisecond)
	}
}
