package ethereumrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/agmanchon/blockchainmonitor/libs/backend/pkg/utils"
	"github.com/agmanchon/txparser/libs/rpc"
	"github.com/agmanchon/txparser/pkg/usecases/txparser"
	"log"
	"math/big"
)

var _ txparser.BlockchainAdapter = (*DefaultEthereumRpcAdapter)(nil)

type DefaultEthereumRpcAdapter struct {
	rpcClient *rpc.RPCClient
	host      string
}

type DefaultEthereumRpcAdapterOptions struct {
	Host string
}

func ProvideDefaultEthereumRpcAdapter(options DefaultEthereumRpcAdapterOptions) (*DefaultEthereumRpcAdapter, error) {
	if options.Host == "nil" {
		return nil, errors.New("option host is mandatory")
	}

	rpcClient := rpc.NewRPCClient(options.Host)
	return &DefaultEthereumRpcAdapter{rpcClient: rpcClient, host: options.Host}, nil
}

func (d DefaultEthereumRpcAdapter) GetMostRecentBlockNumberAdapter(ctx context.Context, input txparser.GetMostRecentBlockNumberAdapterInput) (*txparser.GetMostRecentBlockNumberAdapterOutput, error) {
	response, err := d.rpcClient.Call(ctx, "eth_blockNumber", nil)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("response is nil")
	}
	var hexBlockNumber string
	if err := json.Unmarshal(response, &hexBlockNumber); err != nil {
		return nil, err
	}

	blockNumber, err := utils.ConvertHexToBigInt(hexBlockNumber)
	if err != nil {
		return nil, err
	}
	return &txparser.GetMostRecentBlockNumberAdapterOutput{
		BlockNumber: *blockNumber,
	}, nil
}

func (d DefaultEthereumRpcAdapter) GetBlockByNumberAdapter(ctx context.Context, input txparser.GetBlockByNumberAdapterInput) (*txparser.GetBlockByNumberAdapterOutput, error) {
	if &input.BlockNumber == nil {
		return nil, errors.New("input block number is nil")
	}
	blockNumberHex := utils.ConvertBigIntToHex(&input.BlockNumber)
	params := []interface{}{blockNumberHex, false}
	response, err := d.rpcClient.Call(ctx, "eth_getBlockByNumber", params)
	if err != nil {
		return nil, err
	}

	if response == nil {
		errorMsg := fmt.Sprintf("block not found %d", input.BlockNumber)
		return nil, errors.New(errorMsg)
	}

	var blockTransactions BlockTransactions
	if err := json.Unmarshal(response, &blockTransactions); err != nil {
		return nil, err
	}

	return &txparser.GetBlockByNumberAdapterOutput{TransactionHashes: blockTransactions.TransactionsHashes}, nil
}

func (d DefaultEthereumRpcAdapter) GetTransactionAdapter(ctx context.Context, input txparser.GetTransactionAdapterInput) (*txparser.GetTransactionAdapterOutput, error) {
	params := []interface{}{input.ID}
	response, err := d.rpcClient.Call(ctx, "eth_getTransactionByHash", params)
	if err != nil {
		return nil, err
	}

	if response == nil {
		errorMsg := fmt.Sprintf("transaction %s not found or cleaned up by the node", input.ID)
		log.Println(errorMsg)
		return &txparser.GetTransactionAdapterOutput{}, nil
	}

	var tx Transaction
	if err := json.Unmarshal(response, &tx); err != nil {
		return nil, err
	}

	source := tx.From
	destination := tx.To
	amountWei, _ := new(big.Int).SetString(tx.Value[2:], 16)
	amountEther := new(big.Float).Quo(new(big.Float).SetInt(amountWei), big.NewFloat(1e18))

	gas := new(big.Int)
	gasPrice := new(big.Int)

	gas.SetString(tx.Gas[2:], 16)           // Strip "0x" and convert to decimal
	gasPrice.SetString(tx.GasPrice[2:], 16) // Strip "0x" and convert to decimal

	feeInWei := new(big.Int).Mul(gas, gasPrice)
	feeWeiInEther := new(big.Float).SetInt(feeInWei)
	feeEtherValue := new(big.Float).Quo(feeWeiInEther, big.NewFloat(1e18))

	//fees := fmt.Sprintf("%.18f", feeEtherValue)

	output := &txparser.GetTransactionAdapterOutput{}
	output.Source = source
	output.Destination = destination
	output.AmountInEther = *amountEther
	output.FeesInEther = *feeEtherValue

	return output, nil
}
