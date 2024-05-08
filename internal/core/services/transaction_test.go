package services

import (
	"account-test/internal/core/domain"
	mock_ports "account-test/internal/mocks/ports"
	"account-test/static"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPostTransaction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name            string
		rec             *httptest.ResponseRecorder
		body            map[string]interface{}
		doMockAccRepo   func(repository *mock_ports.MockAccountRepository)
		doMockTransRepo func(repository *mock_ports.MockTransactionRepository)
		err             string
		statusCode      int
	}{
		{
			name: "Test Case Positive",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"source_account_id":      "123",
				"destination_account_id": "1234",
				"amount":                 "19",
			},
			doMockAccRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(
					&domain.Account{ID: "123", Balance: "123"},
					nil,
				)
				repository.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(
					&domain.Account{ID: "1234", Balance: "123"},
					nil,
				)
			},
			doMockTransRepo: func(repository *mock_ports.MockTransactionRepository) {
				repository.EXPECT().ProcessTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			err: "",
		},
		{
			name: "Test Case Negative - Empty account ID",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"source_account_id":      "",
				"destination_account_id": "1234",
				"amount":                 "19",
			},
			doMockAccRepo: func(repository *mock_ports.MockAccountRepository) {
			},
			doMockTransRepo: func(repository *mock_ports.MockTransactionRepository) {
			},
			err:        static.ErrIDLengthCannotBeZero,
			statusCode: 400,
		},
		{
			name: "Test Case Negative - Account ID too long",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"source_account_id":      "123",
				"destination_account_id": "12341239172491274912749124912894129847129471294912748492184",
				"amount":                 "19",
			},
			doMockAccRepo: func(repository *mock_ports.MockAccountRepository) {
			},
			doMockTransRepo: func(repository *mock_ports.MockTransactionRepository) {
			},
			err:        static.ErrIDLengthTooLong,
			statusCode: 400,
		},
		{
			name: "Test Case Negative - Source and Destination account the same",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"source_account_id":      "123",
				"destination_account_id": "123",
				"amount":                 "19",
			},
			doMockAccRepo: func(repository *mock_ports.MockAccountRepository) {
			},
			doMockTransRepo: func(repository *mock_ports.MockTransactionRepository) {
			},
			err:        static.ErrSourceDestinationSame,
			statusCode: 400,
		},
		{
			name: "Test Case Negative - Source account does not exist",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"source_account_id":      "123",
				"destination_account_id": "1234",
				"amount":                 "19",
			},
			doMockAccRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(false)
			},
			doMockTransRepo: func(repository *mock_ports.MockTransactionRepository) {
			},
			err:        static.ErrSourceAccountDoesNotExist,
			statusCode: 400,
		},
		{
			name: "Test Case Negative - Destination account does not exist",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"source_account_id":      "123",
				"destination_account_id": "1234",
				"amount":                 "19",
			},
			doMockAccRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(false)
			},
			doMockTransRepo: func(repository *mock_ports.MockTransactionRepository) {
			},
			err:        static.ErrDestinationAccountDoesNotExist,
			statusCode: 400,
		},
		{
			name: "Test Case Negative - Invalid amount",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"source_account_id":      "123",
				"destination_account_id": "1234",
				"amount":                 "abc",
			},
			doMockAccRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
			},
			doMockTransRepo: func(repository *mock_ports.MockTransactionRepository) {
			},
			err:        static.ErrAmountNotValidNumber,
			statusCode: 400,
		},
		{
			name: "Test Case Negative - Negative amount",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"source_account_id":      "123",
				"destination_account_id": "1234",
				"amount":                 "-10",
			},
			doMockAccRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
			},
			doMockTransRepo: func(repository *mock_ports.MockTransactionRepository) {
			},
			err:        static.ErrAmountCannotBeNegative,
			statusCode: 400,
		},
		{
			name: "Test Case Negative - GetAccount error",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"source_account_id":      "123",
				"destination_account_id": "1234",
				"amount":                 "100",
			},
			doMockAccRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(
					nil,
					errors.New("random error"),
				)
			},
			doMockTransRepo: func(repository *mock_ports.MockTransactionRepository) {
			},
			err:        static.ErrGetSourceAccount,
			statusCode: 500,
		},
		{
			name: "Test Case Negative - Amount greater than source account balance",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"source_account_id":      "123",
				"destination_account_id": "1234",
				"amount":                 "500",
			},
			doMockAccRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(
					&domain.Account{ID: "123", Balance: "123"},
					nil,
				)
				repository.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(
					&domain.Account{ID: "1234", Balance: "123"},
					nil,
				)
			},
			doMockTransRepo: func(repository *mock_ports.MockTransactionRepository) {
			},
			statusCode: 400,
			err:        static.ErrTransferAmountLargerThanAccount,
		},
		{
			name: "Test Case Negative - ProcessTransaction error",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"source_account_id":      "123",
				"destination_account_id": "1234",
				"amount":                 "100",
			},
			doMockAccRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(
					&domain.Account{ID: "123", Balance: "123"},
					nil,
				)
				repository.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(
					&domain.Account{ID: "1234", Balance: "123"},
					nil,
				)
			},
			doMockTransRepo: func(repository *mock_ports.MockTransactionRepository) {
				repository.EXPECT().ProcessTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("random error"))
			},
			statusCode: 500,
			err:        static.ErrUnableToCompleteTransaction,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAccRepo := mock_ports.NewMockAccountRepository(mockCtrl)
			mockTransRepo := mock_ports.NewMockTransactionRepository(mockCtrl)
			tc.doMockAccRepo(mockAccRepo)
			tc.doMockTransRepo(mockTransRepo)
			transSvc := NewTransactionSvc(mockAccRepo, mockTransRepo)
			handler := http.HandlerFunc(transSvc.PostTransaction)
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(body))
			handler.ServeHTTP(tc.rec, req)

			if len(tc.err) > 0 {
				bodyBytes, _ := io.ReadAll(tc.rec.Body)
				assert.Contains(t, string(bodyBytes), tc.err)
				assert.Equal(t, tc.statusCode, tc.rec.Result().StatusCode)
			} else {
				assert.Equal(t, 200, tc.rec.Result().StatusCode)
			}
		})
	}
}
