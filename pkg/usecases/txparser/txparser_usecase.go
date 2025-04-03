package txparser

import (
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"
	"log"
	"math/big"
	"time"
)

type DefaultTxParserUsecase struct {
	blockchainAdapter      BlockchainAdapter
	transactionRepository  TransactionRepository
	addressRepository      AddressRepository
	currentBlockRepository CurrentBlockRepository
	logger                 *logrus.Logger
}

type DefaultTxParserUsecaseOptions struct {
	BlockchainAdapter      BlockchainAdapter
	TransactionRepository  TransactionRepository
	CurrentBlockRepository CurrentBlockRepository
	AddressRepository      AddressRepository
	Logger                 *logrus.Logger
}

var _ TxParserUsecases = (*DefaultTxParserUsecase)(nil)

func ProvideDefaultTxParserUsecase(ctx context.Context, options DefaultTxParserUsecaseOptions) (*DefaultTxParserUsecase, error) {
	if options.BlockchainAdapter == nil {
		return nil, errors.New("blockchain adapter is required")
	}
	if options.TransactionRepository == nil {
		return nil, errors.New("transaction repository is required")
	}
	if options.AddressRepository == nil {
		return nil, errors.New("address repository is required")
	}
	if options.CurrentBlockRepository == nil {
		return nil, errors.New("current block repository is required")
	}
	if options.Logger == nil {
		return nil, errors.New("logger is required")
	}
	return &DefaultTxParserUsecase{
		blockchainAdapter:      options.BlockchainAdapter,
		transactionRepository:  options.TransactionRepository,
		addressRepository:      options.AddressRepository,
		currentBlockRepository: options.CurrentBlockRepository,
		logger:                 options.Logger,
	}, nil
}

func (d DefaultTxParserUsecase) GetCurrentBlock(ctx context.Context) (*GetCurrentBlockOutput, error) {
	blockNumber, err := d.currentBlockRepository.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &GetCurrentBlockOutput{
		BlockNumber: *blockNumber,
	}, nil
}

func (d DefaultTxParserUsecase) Subscribe(ctx context.Context, input SubscribeInput) error {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return err
	}

	_, err = d.addressRepository.Add(ctx, input.Address)
	if err != nil {
		return err
	}
	return nil
}

func (d DefaultTxParserUsecase) GetTransactions(ctx context.Context, input GetTransactionsInput) (*GetTransactionsOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, err
	}
	transactionCollection, err := d.transactionRepository.All(ctx, TransactionListOptions{Address: input.Address})
	if err != nil {
		return nil, err
	}
	return &GetTransactionsOutput{TransactionCollection: *transactionCollection}, nil
}

func (d DefaultTxParserUsecase) StartServer(ctx context.Context, input StartServerInput) {
	for {
		select {
		case <-ctx.Done():
			log.Println("tx parser shutdown")
			return
		default:
			d.processTransactions(ctx)
			t := time.Second
			d.sleep(ctx, t)
		}
	}
}

func (d DefaultTxParserUsecase) processTransactions(ctx context.Context) {

	var currentBlockNumber big.Int
	getBlockFromRepositoryOutput, err := d.currentBlockRepository.Get(ctx)
	if err != nil {
		d.logger.Error(err)
		return
	}
	currentBlockNumber = *getBlockFromRepositoryOutput
	if currentBlockNumber.Cmp(big.NewInt(0)) == 0 {
		getMostRecentBlockNumberAdapterOutput, err := d.blockchainAdapter.GetMostRecentBlockNumberAdapter(ctx, GetMostRecentBlockNumberAdapterInput{})
		if err != nil {
			d.logger.Error(err)
			return
		}
		currentBlockNumber = getMostRecentBlockNumberAdapterOutput.BlockNumber
	} else {
		currentBlockNumber.Add(&currentBlockNumber, big.NewInt(1))
	}
	d.logger.WithFields(logrus.Fields{"blockNumber": currentBlockNumber.String()}).Info("processing")

	getBlockByNumberAdapterOutput, err := d.blockchainAdapter.GetBlockByNumberAdapter(ctx, GetBlockByNumberAdapterInput{BlockNumber: currentBlockNumber})
	if err != nil {
		d.logger.Error(err)
		return
	}

	_, err = d.currentBlockRepository.Update(ctx, currentBlockNumber)
	if err != nil {
		d.logger.Error(err)
	}

	for _, transactionHash := range getBlockByNumberAdapterOutput.TransactionHashes {
		getTransactionByHashInput := GetTransactionAdapterInput{}
		getTransactionByHashInput.ID = transactionHash
		getTransactionOutput, err := d.blockchainAdapter.GetTransactionAdapter(ctx, getTransactionByHashInput)
		if err != nil {
			d.logger.Error(err)
			return
		}
		addressSubscribedCollection, err := d.addressRepository.All(ctx, AddressListOptions{})
		if err != nil {
			d.logger.Error(err)
			return
		}
		for _, addressSubscribed := range addressSubscribedCollection.Items {
			if addressSubscribed == getTransactionOutput.Source || addressSubscribed == getTransactionOutput.Destination {
				d.logger.WithFields(logrus.Fields{"address": addressSubscribed, "id": transactionHash}).Info("transaction stored")
				transactionToStored := Transaction{}
				transactionToStored.ID = transactionHash
				transactionToStored.Source = getTransactionOutput.Source
				transactionToStored.Destination = getTransactionOutput.Destination
				transactionToStored.FeesInEther = getTransactionOutput.FeesInEther
				transactionToStored.AmountInEther = getTransactionOutput.AmountInEther
				transactionToStored.CreationDate = time.Now()
				_, err := d.transactionRepository.Add(ctx, transactionToStored)
				if err != nil {
					d.logger.Error(err)
					return
				}
			}
		}
	}

}

func (d DefaultTxParserUsecase) sleep(ctx context.Context, t time.Duration) {
	timeoutchan := make(chan bool)
	go func() {
		<-time.After(t)
		timeoutchan <- true
	}()

	select {
	case <-timeoutchan:
		break
	case <-ctx.Done():
		d.logger.Info("terminated sleep due to a context cancellation")
	}

}
