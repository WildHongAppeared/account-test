package repositories

import (
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

func (i *TransactionPortImpl) UpdateTransaction(ctx context.Context, source_id string, destination_id string, source_amount float64, destination_amount float64) error {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
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

	_, err = tx.ExecContext(
		ctx,
		query,
		source_amount,
		source_id,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		query,
		destination_amount,
		destination_id,
	)

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
