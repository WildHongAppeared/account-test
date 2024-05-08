package services

import (
	"account-test/internal/core/domain"
	mock_ports "account-test/internal/mocks/ports"
	"account-test/static"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAccount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name       string
		rec        *httptest.ResponseRecorder
		req        *http.Request
		account_id string
		doMockRepo func(repository *mock_ports.MockAccountRepository)
		want       domain.Account
		err        string
		statusCode int
	}{
		{
			name:       "Test Case Positive",
			rec:        httptest.NewRecorder(),
			req:        httptest.NewRequest("GET", "/accounts/{account_id}", nil),
			account_id: "123",
			doMockRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(
					&domain.Account{ID: "123", Balance: "123"},
					nil,
				)
			},
			want: domain.Account{ID: "123", Balance: "123"},
			err:  "",
		},
		{
			name:       "Test Case Negative - Empty account passed as parameter",
			rec:        httptest.NewRecorder(),
			req:        httptest.NewRequest("GET", "/accounts/{account_id}", nil),
			account_id: "",
			doMockRepo: func(repository *mock_ports.MockAccountRepository) {
			},
			want:       domain.Account{},
			err:        static.ErrIDLengthCannotBeZero,
			statusCode: 400,
		},
		{
			name:       "Test Case Negative - Account parameter longer than 32 char",
			rec:        httptest.NewRecorder(),
			req:        httptest.NewRequest("GET", "/accounts/{account_id}", nil),
			account_id: "99999999999999999999999999999999999999999999999999999999999999999999",
			doMockRepo: func(repository *mock_ports.MockAccountRepository) {
			},
			want:       domain.Account{},
			err:        static.ErrIDLengthTooLong,
			statusCode: 400,
		},
		{
			name:       "Test Case Negative - Account does not exist",
			rec:        httptest.NewRecorder(),
			req:        httptest.NewRequest("GET", "/accounts/{account_id}", nil),
			account_id: "123",
			doMockRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(false)
			},
			want:       domain.Account{},
			err:        static.ErrAccountDoesNotExist,
			statusCode: 400,
		},
		{
			name:       "Test Case Negative - Repository error",
			rec:        httptest.NewRecorder(),
			req:        httptest.NewRequest("GET", "/accounts/{account_id}", nil),
			account_id: "123",
			doMockRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
				repository.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(
					nil,
					errors.New("random error"),
				)
			},
			want:       domain.Account{},
			err:        static.ErrUnableToRetrieveAccount,
			statusCode: 500,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAccRepo := mock_ports.NewMockAccountRepository(mockCtrl)
			tc.doMockRepo(mockAccRepo)
			accSvc := NewAccountSvc(mockAccRepo)
			handler := http.HandlerFunc(accSvc.GetAccount)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("account_id", tc.account_id)

			r := tc.req.WithContext(context.WithValue(tc.req.Context(), chi.RouteCtxKey, rctx))
			handler.ServeHTTP(tc.rec, r)

			if len(tc.err) > 0 {
				bodyBytes, _ := io.ReadAll(tc.rec.Body)
				assert.Contains(t, string(bodyBytes), tc.err)
				assert.Equal(t, tc.statusCode, tc.rec.Result().StatusCode)
			} else {
				var response domain.Account
				_ = json.NewDecoder(tc.rec.Body).Decode(&response)
				assert.Equal(t, tc.want, response)
				assert.Equal(t, 200, tc.rec.Result().StatusCode)
			}
		})
	}
}

func TestPostAccount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name       string
		rec        *httptest.ResponseRecorder
		body       map[string]interface{}
		doMockRepo func(repository *mock_ports.MockAccountRepository)
		err        string
		statusCode int
	}{
		{
			name: "Test Case Positive",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"account_id":      "123",
				"initial_balance": "123",
			},
			doMockRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(false)
				repository.EXPECT().InsertAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil,
				)
			},
			err: "",
		},
		{
			name: "Test Case Negative - Empty account passed as parameter",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"account_id":      "",
				"initial_balance": "",
			},
			doMockRepo: func(repository *mock_ports.MockAccountRepository) {
			},
			err:        static.ErrIDLengthCannotBeZero,
			statusCode: 400,
		},
		{
			name: "Test Case Negative - Account parameter longer than 32 char",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"account_id":      "99999999999999999999999999999999999999999999999999999999999999999999",
				"initial_balance": "123",
			},
			doMockRepo: func(repository *mock_ports.MockAccountRepository) {
			},
			err:        static.ErrIDLengthTooLong,
			statusCode: 400,
		},
		{
			name: "Test Case Negative - Account already exist",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"account_id":      "123",
				"initial_balance": "123",
			},
			doMockRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(true)
			},
			err:        static.ErrAccountAlreadyExist,
			statusCode: 400,
		},
		{
			name: "Test Case Negative - initial_balance not a number",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"account_id":      "123",
				"initial_balance": "abc",
			},
			doMockRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(false)
			},
			err:        static.ErrBalanceNotValidNumber,
			statusCode: 400,
		},
		{
			name: "Test Case Negative - negative initial_balance",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"account_id":      "123",
				"initial_balance": "-123",
			},
			doMockRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(false)
			},
			err:        static.ErrBalanceCannotBeNegative,
			statusCode: 400,
		},
		{
			name: "Test Case Negative - repository error",
			rec:  httptest.NewRecorder(),
			body: map[string]interface{}{
				"account_id":      "123",
				"initial_balance": "123",
			},
			doMockRepo: func(repository *mock_ports.MockAccountRepository) {
				repository.EXPECT().CheckAccountExists(gomock.Any(), gomock.Any()).Return(false)
				repository.EXPECT().InsertAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					errors.New("random error"),
				)
			},
			err:        static.ErrCreatingAccount,
			statusCode: 500,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAccRepo := mock_ports.NewMockAccountRepository(mockCtrl)
			tc.doMockRepo(mockAccRepo)
			accSvc := NewAccountSvc(mockAccRepo)
			handler := http.HandlerFunc(accSvc.PostAccount)
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest("POST", "/accounts", bytes.NewReader(body))
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
