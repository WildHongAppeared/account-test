package repositories

import (
	"account-test/internal/core/domain"
	"account-test/postgres"
	"account-test/static"
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type AccountPortImpl struct {
	db       *sqlx.DB
	dbConfig *postgres.DBConfig
}

func NewAccountPort(db *sqlx.DB, dbConfig *postgres.DBConfig) *AccountPortImpl {
	return &AccountPortImpl{
		db:       db,
		dbConfig: dbConfig,
	}
}

func (i *AccountPortImpl) InsertAccount(ctx context.Context, id string, balance float64) error {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := fmt.Sprintf(`
			INSERT INTO %s.%s( 
				id, balance 
			)
			VALUES (
				$1, $2
			)
		`,
		i.dbConfig.Schema, static.TableAccount,
	)

	_, err = i.db.ExecContext(
		ctx,
		query,
		id,
		balance,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (i *AccountPortImpl) CheckAccountExists(ctx context.Context, id string) bool {
	var isExist bool
	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s.%s WHERE id = $1)`, i.dbConfig.Schema, static.TableAccount)
	err := i.db.QueryRowContext(ctx, query, id).Scan(&isExist)
	if err != nil {
		log.Println("CheckAccountExists-Repository-Error : " + err.Error())
	}
	return isExist

}

func (i *AccountPortImpl) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	query := fmt.Sprintf(`
	SELECT 
		id, balance
	FROM %s.%s 
	WHERE id = $1`,
		i.dbConfig.Schema, static.TableAccount,
	)

	var response domain.Account
	err := i.db.QueryRowContext(ctx, query, id).Scan(
		&response.ID,
		&response.Balance,
	)
	if err != nil {
		return nil, err
	}

	return &response, nil

}
