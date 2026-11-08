package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/amirjbr/shared-expense/internal/ledger/core/entity"
)

// TODO later we can use batch inserts for better performance . read about batch insert
type LedgerRepo struct {
	DB *sql.DB
}

func NewLedgerRepo(db *sql.DB) (*LedgerRepo, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &LedgerRepo{DB: db}, nil
}

// TODO later we can use batch inserts for better performance . read about batch insert
func (r *LedgerRepo) CreateExpense(ctx context.Context, expense entity.Expense) error {

	const insertExpenseQuery = `
		INSERT INTO expenses (
			id,
			group_id,
			title,
			created_by_user_id,
			total_amount,
			split_type,
			occurred_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	const insertPayerQuery = `
		INSERT INTO expense_payers (
			expense_id,
			group_id,
			user_id,
			amount
		)
		VALUES ($1, $2, $3, $4)
	`

	const insertShareQuery = `
		INSERT INTO expense_shares (
			expense_id,
			group_id,
			user_id,
			amount
		)
		VALUES ($1, $2, $3, $4)
	`

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin create expense transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	// 1. Insert parent expense
	_, err = tx.ExecContext(
		ctx,
		insertExpenseQuery,
		expense.ID,
		expense.GroupID,
		expense.Title,
		expense.CreatedByUserID,
		expense.TotalAmount,
		expense.SplitType,
		expense.OccurredAt,
		expense.CreatedAt,
		expense.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert expense: %w", err)
	}

	// 2. Insert all payers
	for _, payer := range expense.Payers {
		_, err = tx.ExecContext(
			ctx,
			insertPayerQuery,
			expense.ID,
			expense.GroupID,
			payer.UserID,
			payer.Amount,
		)
		if err != nil {
			return fmt.Errorf(
				"insert expense payer %s: %w",
				payer.UserID,
				err,
			)
		}
	}

	// 3. Insert all shares
	for _, share := range expense.Shares {
		_, err = tx.ExecContext(
			ctx,
			insertShareQuery,
			expense.ID,
			expense.GroupID,
			share.UserID,
			share.Amount,
		)
		if err != nil {
			return fmt.Errorf(
				"insert expense share %s: %w",
				share.UserID,
				err,
			)
		}
	}

	// 4. Commit everything
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit create expense transaction: %w", err)
	}

	return nil
}
