package services

import (
	"account-test/internal/core/domain"
	"account-test/internal/core/ports"
	"account-test/internal/core/utils"
	"account-test/static"
	"context"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type AccountSvcImpl struct {
	accountRepo ports.AccountRepository
}

func NewAccountSvc(accountRepo ports.AccountRepository) *AccountSvcImpl {
	return &AccountSvcImpl{
		accountRepo: accountRepo,
	}
}

// PostAccount will accept a HTTP body containing a domain.PostAccount object
// The function will check if the inputs from domain.PostAccount object are valid inputs
// The function will check if the id from domain.PostAccount belongs to an existing account
// The function will fix the balance value to a floating point precision of 5
// The function will create the account with the payload from domain.PostAccount in the account table if all checks are valid
// The function will return HTTP status OK and no body if the creation is successful
func (srv *AccountSvcImpl) PostAccount(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	postAccountBody := domain.PostAccount{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, static.ErrUnableToReadBody, http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &postAccountBody)
	if err != nil {
		http.Error(w, static.ErrUnableToReadBody, http.StatusBadRequest)
		return
	}
	if len(postAccountBody.ID) == 0 {
		http.Error(w, static.ErrIDLengthCannotBeZero, http.StatusBadRequest)
		return
	}
	if len(postAccountBody.ID) > 32 {
		http.Error(w, static.ErrIDLengthTooLong, http.StatusBadRequest)
		return
	}
	accountAlreadyExists := srv.accountRepo.CheckAccountExists(ctx, postAccountBody.ID)
	if accountAlreadyExists {
		http.Error(w, static.ErrAccountAlreadyExist, http.StatusBadRequest)
		return
	}
	accountBalance, err := strconv.ParseFloat(postAccountBody.Balance, 64)
	if err != nil {
		http.Error(w, static.ErrBalanceNotValidNumber, http.StatusBadRequest)
		return
	}
	if accountBalance < 0 {
		http.Error(w, static.ErrBalanceCannotBeNegative, http.StatusBadRequest)
		return
	}
	if accountBalance > math.MaxFloat64 {
		http.Error(w, static.ErrBalanceTooLarge, http.StatusBadRequest)
		return
	}

	err = srv.accountRepo.InsertAccount(ctx, postAccountBody.ID, utils.ToFixed(accountBalance, 5))
	if err != nil {
		log.Println("InsertAccount error - ", err.Error())
		http.Error(w, static.ErrCreatingAccount, http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, http.StatusOK, nil)
}

// GetAccount will accept a HTTP path parameter of account_id
// the function will check if account_id is a valid input
// the function will check if the account_id belongs to an existing account in the system
// the function will then retrieve all the account details associated with the account_id, returned as a domain.Account object
func (srv *AccountSvcImpl) GetAccount(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	accountId := chi.URLParam(r, "account_id")
	if len(accountId) == 0 {
		http.Error(w, static.ErrIDLengthCannotBeZero, http.StatusBadRequest)
		return
	}
	if len(accountId) > 32 {
		http.Error(w, static.ErrIDLengthTooLong, http.StatusBadRequest)
		return
	}
	accountAlreadyExists := srv.accountRepo.CheckAccountExists(ctx, accountId)
	if !accountAlreadyExists {
		http.Error(w, static.ErrAccountDoesNotExist, http.StatusBadRequest)
		return
	}
	account, err := srv.accountRepo.GetAccount(ctx, accountId)
	if err != nil {
		log.Println("GetAccount error - ", err.Error())
		http.Error(w, static.ErrUnableToRetrieveAccount, http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, http.StatusOK, account)

}
