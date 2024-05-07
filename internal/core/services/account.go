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

type UserSvc interface {
	PostAccount(w http.ResponseWriter, r *http.Request)
	GetAccount(w http.ResponseWriter, r *http.Request)
}

func NewAccountSvc(accountRepo ports.AccountRepository) *AccountSvcImpl {
	return &AccountSvcImpl{
		accountRepo: accountRepo,
	}
}

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

	err = srv.accountRepo.InsertAccount(ctx, postAccountBody.ID, accountBalance)
	if err != nil {
		log.Println("InsertAccount error - ", err.Error())
		http.Error(w, static.ErrCreatingAccount, http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, http.StatusOK, nil)
}

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
		http.Error(w, static.ErrUnableToRetrieveAccount, http.StatusBadRequest)
		return
	}
	utils.JSONResponse(w, http.StatusOK, account)

}
