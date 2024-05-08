package repositories

import (
	"account-test/internal/core/domain"
	"account-test/postgres"
	"account-test/static"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TransactionPortImpl struct {
	db       *sqlx.DB
	dbConfig *postgres.DBConfig
}

func NewTransactionPort(db *sqlx.DB, dbConfig *postgres.DBConfig) *TransactionPortImpl {
	return &TransactionPortImpl{
		db:       db,
		dbConfig: dbConfig,
	}
}

// ProcessTransaction accepts a Transaction object, source_amount and destination_amount to update the amount for the account with the transaction.sourceID with source_amount and account with transaction.destinationID with destination_amount
// The function will also call insertTransaction to create a new transaction in the DB for logging of the transactions details
// The function will also call updateTransactionWithErrorMessage to update the created transaction with error message in the event of error happening
// The function will return nil if there is no error and an error object of there is error
func (i *TransactionPortImpl) ProcessTransaction(ctx context.Context, transaction domain.Transaction, source_amount float64, destination_amount float64) error {
	transactionId, err := i.insertTransaction(ctx, transaction) //Insert transaction for logging purpose
	if err != nil {
		return err
	}
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		i.updateTransactionWithErrorMessage(ctx, err.Error(), transactionId) // Update transaction with error message
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := fmt.Sprintf(`
		UPDATE %s.%s SET 
			balance = $1,
			updated_at = NOW()
		WHERE id = $2`,
		i.dbConfig.Schema, static.TableAccount,
	)

	_, err = tx.ExecContext( //Update source account with new amount
		ctx,
		query,
		source_amount,
		transaction.SourceID,
	)
	if err != nil {
		i.updateTransactionWithErrorMessage(ctx, err.Error(), transactionId) // Update transaction with error message
		return err
	}

	_, err = tx.ExecContext( //Update destination account with new amount
		ctx,
		query,
		destination_amount,
		transaction.DestinationID,
	)
	if err != nil {
		i.updateTransactionWithErrorMessage(ctx, err.Error(), transactionId) // Update transaction with error message
		return err
	}

	err = tx.Commit()
	if err != nil {
		i.updateTransactionWithErrorMessage(ctx, err.Error(), transactionId) // Update transaction with error message
		return err
	}
	return nil
}

// updateTransactionWithErrorMessage will accept a error message and the ID of a transaction to update the transaction row in DB with the error message for logging purpose
// The function will return nil if there is no error and an error object of there is error
func (i *TransactionPortImpl) updateTransactionWithErrorMessage(ctx context.Context, message string, id int) error {

	query := fmt.Sprintf(`
		UPDATE %s.%s SET 
			error_message = $1,
		WHERE id = $2`,
		i.dbConfig.Schema, static.TableTransaction,
	)

	_, err := i.db.ExecContext(
		ctx,
		query,
		message,
		id,
	)

	if err != nil {
		return err
	}
	return nil
}

// insertTransaction will accept a domain.Transaction object to create a new row in the transaction table to log the transaction details
// The function will return the id of the created transaction object and an error object of there is error
func (i *TransactionPortImpl) insertTransaction(ctx context.Context, transaction domain.Transaction) (int, error) {

	query := fmt.Sprintf(`
	INSERT INTO %s.%s( 
		source_account_id, destination_account_id, amount 
	)
	VALUES (
		$1, $2, $3
	) RETURNING id
	`,
		i.dbConfig.Schema, static.TableTransaction,
	)

	row := i.db.QueryRowContext(
		ctx,
		query,
		transaction.SourceID,
		transaction.DestinationID,
		transaction.Amount,
	)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
