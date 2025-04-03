package http

import (
	"context"
	"errors"
	"github.com/agmanchon/txparser/pkg/infra/httpinfragenerated"
	"github.com/agmanchon/txparser/pkg/usecases/txparser"
)

type DefaultTxParserHttpAdapter struct {
	txParserUsecases txparser.TxParserUsecases
}

type DefaultTxParserHttpAdapterOptions struct {
	TxParserUsecases txparser.TxParserUsecases
}

func ProvideTxParserHttpAdapter(options DefaultTxParserHttpAdapterOptions) (*DefaultTxParserHttpAdapter, error) {
	if options.TxParserUsecases == nil {
		return nil, errors.New("mandatory TxParserUsecases not provided")
	}
	return &DefaultTxParserHttpAdapter{
		txParserUsecases: options.TxParserUsecases,
	}, nil
}

var _ httpinfragenerated.TxParserAPIServicer = (*DefaultTxParserHttpAdapter)(nil)

func (d DefaultTxParserHttpAdapter) GetCurrentBlock(ctx context.Context) (httpinfragenerated.ImplResponse, error) {
	output, err := d.txParserUsecases.GetCurrentBlock(ctx)
	if err != nil {
		return httpinfragenerated.Response(500, httpinfragenerated.ErrorResponse{Error: err.Error()}), nil
	}
	return httpinfragenerated.Response(200, httpinfragenerated.BlockNumberResponse{BlockNumber: output.BlockNumber.Int64()}), nil
}

func (d DefaultTxParserHttpAdapter) SubscribeAddress(ctx context.Context, request httpinfragenerated.SubscriptionRequest) (httpinfragenerated.ImplResponse, error) {
	input := txparser.SubscribeInput{}
	input.Address = request.Address
	err := d.txParserUsecases.Subscribe(ctx, input)
	if err != nil {
		return httpinfragenerated.Response(200, httpinfragenerated.ErrorResponse{Error: err.Error()}), nil
	}
	return httpinfragenerated.Response(200, httpinfragenerated.SuccessResponse{Success: true}), nil
}

func (d DefaultTxParserHttpAdapter) GetAddressTransactions(ctx context.Context, s string) (httpinfragenerated.ImplResponse, error) {
	input := txparser.GetTransactionsInput{}
	input.Address = s
	transactionCollection, err := d.txParserUsecases.GetTransactions(ctx, input)
	if err != nil {
		return httpinfragenerated.Response(500, httpinfragenerated.ErrorResponse{Error: err.Error()}), nil
	}
	response := []httpinfragenerated.Transaction{}
	for _, transaction := range transactionCollection.Items {
		outputTransaction := httpinfragenerated.Transaction{}
		outputTransaction.Source = transaction.Source
		outputTransaction.Destination = transaction.Destination
		outputTransaction.BlockNumber = transaction.BlockNumber.Int64()
		outputTransaction.Id = transaction.ID
		outputTransaction.AmountInEther = transaction.AmountInEther.String()
		outputTransaction.FeesInEther = transaction.FeesInEther.String()

		response = append(response, outputTransaction)
	}
	return httpinfragenerated.Response(200, response), nil
}
