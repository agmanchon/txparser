// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Blockchain Transaction Parser API
 *
 * API for parsing and tracking blockchain transactions with subscription capabilities.
 *
 * API version: 1.0.0
 */

package httpinfragenerated

import (
	"context"
	"errors"
	"net/http"
)

// TxParserAPIService is a service that implements the logic for the TxParserAPIServicer
// This service should implement the business logic for every endpoint for the TxParserAPI API.
// Include any external packages or services that will be required by this service.
type TxParserAPIService struct {
}

// NewTxParserAPIService creates a default api service
func NewTxParserAPIService() *TxParserAPIService {
	return &TxParserAPIService{}
}

// GetCurrentBlock - Get the last parsed block
func (s *TxParserAPIService) GetCurrentBlock(ctx context.Context) (ImplResponse, error) {
	// TODO - update GetCurrentBlock with the required logic for this service method.
	// Add api_tx_parser_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, BlockNumberResponse{}) or use other options such as http.Ok ...
	// return Response(200, BlockNumberResponse{}), nil

	// TODO: Uncomment the next line to return response Response(500, ErrorResponse{}) or use other options such as http.Ok ...
	// return Response(500, ErrorResponse{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetCurrentBlock method not implemented")
}

// SubscribeAddress - Subscribe an address for transaction tracking
func (s *TxParserAPIService) SubscribeAddress(ctx context.Context, subscriptionRequest SubscriptionRequest) (ImplResponse, error) {
	// TODO - update SubscribeAddress with the required logic for this service method.
	// Add api_tx_parser_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, SuccessResponse{}) or use other options such as http.Ok ...
	// return Response(200, SuccessResponse{}), nil

	// TODO: Uncomment the next line to return response Response(400, ErrorResponse{}) or use other options such as http.Ok ...
	// return Response(400, ErrorResponse{}), nil

	// TODO: Uncomment the next line to return response Response(500, ErrorResponse{}) or use other options such as http.Ok ...
	// return Response(500, ErrorResponse{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("SubscribeAddress method not implemented")
}

// GetAddressTransactions - Get transactions for a specific address
func (s *TxParserAPIService) GetAddressTransactions(ctx context.Context, address string) (ImplResponse, error) {
	// TODO - update GetAddressTransactions with the required logic for this service method.
	// Add api_tx_parser_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, []Transaction{}) or use other options such as http.Ok ...
	// return Response(200, []Transaction{}), nil

	// TODO: Uncomment the next line to return response Response(400, ErrorResponse{}) or use other options such as http.Ok ...
	// return Response(400, ErrorResponse{}), nil

	// TODO: Uncomment the next line to return response Response(404, ErrorResponse{}) or use other options such as http.Ok ...
	// return Response(404, ErrorResponse{}), nil

	// TODO: Uncomment the next line to return response Response(500, ErrorResponse{}) or use other options such as http.Ok ...
	// return Response(500, ErrorResponse{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetAddressTransactions method not implemented")
}
