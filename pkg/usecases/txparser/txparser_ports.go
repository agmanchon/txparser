package txparser

import (
	"context"
	"math/big"
	"time"
)

type BlockchainAdapter interface {
	GetMostRecentBlockNumberAdapter(context.Context, GetMostRecentBlockNumberAdapterInput) (*GetMostRecentBlockNumberAdapterOutput, error)
	GetBlockByNumberAdapter(context.Context, GetBlockByNumberAdapterInput) (*GetBlockByNumberAdapterOutput, error)
	GetTransactionAdapter(context.Context, GetTransactionAdapterInput) (*GetTransactionAdapterOutput, error)
}

type TransactionRepository interface {
	Add(ctx context.Context, data Transaction) (*Transaction, error)
	All(ctx context.Context, filter TransactionListOptions) (*TransactionCollection, error)
}

type AddressRepository interface {
	Add(ctx context.Context, data string) (*string, error)
	All(ctx context.Context, filter AddressListOptions) (*AddressCollection, error)
}

type CurrentBlockRepository interface {
	Add(ctx context.Context, data big.Int) (*big.Int, error)
	Get(ctx context.Context) (*big.Int, error)
	Update(ctx context.Context, data big.Int) (*big.Int, error)
}

type GetMostRecentBlockNumberAdapterInput struct {
}

type GetMostRecentBlockNumberAdapterOutput struct {
	BlockNumber big.Int
}

type GetBlockByNumberAdapterInput struct {
	BlockNumber big.Int
}

type GetBlockByNumberAdapterOutput struct {
	TransactionHashes []string
}

type GetTransactionAdapterInput struct {
	ID string
}

type GetTransactionAdapterOutput struct {
	Source        string
	Destination   string
	AmountInEther big.Float
	FeesInEther   big.Float
}

type Transaction struct {
	ID            string
	BlockNumber   big.Int
	Source        string
	Destination   string
	AmountInEther big.Float
	FeesInEther   big.Float
	CreationDate  time.Time
}

type TransactionListOptions struct {
	Address string
}

type TransactionCollection struct {
	Items []Transaction
}

type AddressListOptions struct {
}

type AddressCollection struct {
	Items []string
}
