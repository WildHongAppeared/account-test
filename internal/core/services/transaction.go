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
)

type TransactionSvcImpl struct {
	accountRepo     ports.AccountRepository
	transactionRepo ports.TransactionRepository
}

func NewTransactionSvc(accountRepo ports.AccountRepository, transactionRepo ports.TransactionRepository) *TransactionSvcImpl {
	return &TransactionSvcImpl{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

// PostTransaction will accept a HTTP body containing a domain.Transaction object
// The function will check if the inputs from domain.Transaction object are valid inputs
// The function will check if the source account and destination account, denoted by SourceID and DestinationID, is a valid account within the system
// The function will calculate the amount to be transferred from the balance of source account to destination account
// The function will process the transaction according to the calculated balance of the source and destination account
// The function will fix all calculated values to a floating point precision of 5
// The function will return HTTP status OK and no body if the creation is successful
func (srv *TransactionSvcImpl) PostTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	postTransactionBody := domain.Transaction{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, static.ErrUnableToReadBody, http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &postTransactionBody)
	if err != nil {
		http.Error(w, static.ErrUnableToReadBody, http.StatusBadRequest)
		return
	}
	if len(postTransactionBody.SourceID) == 0 || len(postTransactionBody.DestinationID) == 0 {
		http.Error(w, static.ErrIDLengthCannotBeZero, http.StatusBadRequest)
		return
	}
	if len(postTransactionBody.SourceID) > 32 || len(postTransactionBody.DestinationID) > 32 {
		http.Error(w, static.ErrIDLengthTooLong, http.StatusBadRequest)
		return
	}
	if postTransactionBody.SourceID == postTransactionBody.DestinationID {
		http.Error(w, static.ErrSourceDestinationSame, http.StatusBadRequest)
		return
	}
	sourceAccountExists := srv.accountRepo.CheckAccountExists(ctx, postTransactionBody.SourceID)
	if !sourceAccountExists {
		http.Error(w, static.ErrSourceAccountDoesNotExist, http.StatusBadRequest)
		return
	}
	destinationAccountExists := srv.accountRepo.CheckAccountExists(ctx, postTransactionBody.SourceID)
	if !destinationAccountExists {
		http.Error(w, static.ErrDestinationAccountDoesNotExist, http.StatusBadRequest)
		return
	}
	transferAmount, err := strconv.ParseFloat(postTransactionBody.Amount, 64)
	if err != nil {
		http.Error(w, static.ErrAmountNotValidNumber, http.StatusBadRequest)
		return
	}
	if transferAmount <= 0 {
		http.Error(w, static.ErrAmountCannotBeNegative, http.StatusBadRequest)
		return
	}
	if transferAmount > math.MaxFloat64 {
		http.Error(w, static.ErrAmountTooLarge, http.StatusBadRequest)
		return
	}
	transferAmount = utils.ToFixed(transferAmount, 5)

	sourceAccount, err := srv.accountRepo.GetAccount(ctx, postTransactionBody.SourceID)
	if err != nil {
		log.Println("GetAccount - Source - error - ", err.Error())
		http.Error(w, static.ErrGetSourceAccount, http.StatusInternalServerError)
		return
	}
	destinationAccount, err := srv.accountRepo.GetAccount(ctx, postTransactionBody.DestinationID)
	if err != nil {
		log.Println("GetAccount - Destination - error - ", err.Error())
		http.Error(w, static.ErrGetDestinationAccount, http.StatusInternalServerError)
		return
	}
	sourceAccountAmount, err := strconv.ParseFloat(sourceAccount.Balance, 64)
	if sourceAccountAmount < transferAmount {
		http.Error(w, static.ErrTransferAmountLargerThanAccount, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, static.ErrAmountNotValidNumber, http.StatusBadRequest)
		return
	}
	destinationAccountAmount, err := strconv.ParseFloat(destinationAccount.Balance, 64)
	if err != nil {
		http.Error(w, static.ErrAmountNotValidNumber, http.StatusBadRequest)
		return
	}
	sourceAccountAmount -= transferAmount
	destinationAccountAmount += transferAmount
	err = srv.transactionRepo.ProcessTransaction(ctx, postTransactionBody, utils.ToFixed(sourceAccountAmount, 5), utils.ToFixed(destinationAccountAmount, 5))
	if err != nil {
		log.Println("UpdateTransaction error - ", err.Error())
		http.Error(w, static.ErrUnableToCompleteTransaction, http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, http.StatusOK, nil)
}
