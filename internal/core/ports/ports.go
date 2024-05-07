package ports

import (
	"account-test/internal/core/domain"
	"context"
)

type AccountRepository interface {
	InsertAccount(ctx context.Context, id string, balance float64) error
	GetAccount(ctx context.Context, id string) (*domain.Account, error)
	CheckAccountExists(ctx context.Context, id string) bool
}
