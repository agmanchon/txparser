package txparser

import (
	"context"
	"math/big"
)

type TxParserUsecases interface {
	GetCurrentBlock(ctx context.Context) (*GetCurrentBlockOutput, error)
	Subscribe(ctx context.Context, input SubscribeInput) error
	GetTransactions(ctx context.Context, input GetTransactionsInput) (*GetTransactionsOutput, error)
	StartServer(ctx context.Context, input StartServerInput)
}

type GetCurrentBlockOutput struct {
	BlockNumber big.Int
}
type SubscribeInput struct {
	Address string `valid:"required"`
}
type GetTransactionsInput struct {
	Address string `valid:"required"`
}

type GetTransactionsOutput struct {
	TransactionCollection
}

type StartServerInput struct {
}
