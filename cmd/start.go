package cmd

import (
	"context"
	"github.com/agmanchon/txparser/config"
	httptxparser "github.com/agmanchon/txparser/pkg/adapters/input/http"
	"github.com/agmanchon/txparser/pkg/adapters/output/ethereumrpc"
	"github.com/agmanchon/txparser/pkg/adapters/output/memorystorage"
	"github.com/agmanchon/txparser/pkg/infra/httpinfragenerated"
	"github.com/agmanchon/txparser/pkg/usecases/txparser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the application " + appName,
	Run:   startRun,
}

func startRun(_ *cobra.Command, _ []string) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.ReadInstallConfig()
	checkError(err)

	log := logrus.New()
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	// Ethereum Rpc Adapter
	ethereumRpcAdapter, err := ethereumrpc.ProvideDefaultEthereumRpcAdapter(ethereumrpc.DefaultEthereumRpcAdapterOptions{Host: cfg.EthereumUrl})
	checkError(err)

	//transactionMemoryStorage Adapter
	transactionMemoryStorageAdapter, err := memorystorage.ProvideDefaultTransactionMemoryStorageAdapter(memorystorage.DefaultTransactionMemoryStorageAdapterOptions{})
	checkError(err)

	//currentBlockNumberMemoryStorage Adapter
	currentBlockNumberMemoryStorageAdapter, err := memorystorage.ProvideDefaultCurrentBlockMemoryStorageAdapter(memorystorage.DefaultCurrentBlockMemoryStorageAdapterOptions{})
	checkError(err)

	//addressesMemoryStorage Adapter
	addressesMemoryStorageAdapter, err := memorystorage.ProvideDefaultAddressMemoryStorageAdapter(memorystorage.DefaultAddressMemoryStorageAdapterOptions{})
	checkError(err)

	//txParserUsecase
	txParserOptions := txparser.DefaultTxParserUsecaseOptions{}
	txParserOptions.CurrentBlockRepository = currentBlockNumberMemoryStorageAdapter
	txParserOptions.AddressRepository = addressesMemoryStorageAdapter
	txParserOptions.TransactionRepository = transactionMemoryStorageAdapter
	txParserOptions.BlockchainAdapter = ethereumRpcAdapter
	txParserOptions.Logger = log
	txParserUsecases, err := txparser.ProvideDefaultTxParserUsecase(ctx, txParserOptions)
	checkError(err)

	//httpInfra
	txParserhttpInfra, err := httptxparser.ProvideTxParserHttpAdapter(httptxparser.DefaultTxParserHttpAdapterOptions{TxParserUsecases: txParserUsecases})
	checkError(err)

	TxParserAPIController := httpinfragenerated.NewTxParserAPIController(txParserhttpInfra)

	router := httpinfragenerated.NewRouter(TxParserAPIController)

	go txParserUsecases.StartServer(ctx, txparser.StartServerInput{})
	fs := http.FileServer(http.Dir("./swagger-ui"))
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", fs))

	// Serve the OpenAPI specification
	router.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api-spec/openapi.yaml")
	})

	go http.ListenAndServe(":8080", router)

	for {
		select {
		case <-signalChan:
			cancel()
			log.Println("Shutdown complete.")
			os.Exit(0)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
