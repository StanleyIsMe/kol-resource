package sqlboiler

import (
	"context"
	"database/sql"
	commonErrors "kolresource/internal/common/errors"

	"github.com/aarondl/sqlboiler/v4/boil"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	CtxTransactionKey contextKey = "CtxTransaction"
)

func (r *EmailRepository) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return commonErrors.TransactionError{Err: err}
	}

	err = fn(context.WithValue(ctx, CtxTransactionKey, tx))
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return &commonErrors.TransactionRollbackError{DbErr: err, TxErr: rollbackErr}
		}

		return commonErrors.TransactionError{Err: err}
	}

	err = tx.Commit()
	if err != nil {
		return commonErrors.TransactionError{Err: err}
	}

	return nil
}

//nolint:ireturn
func (r *EmailRepository) getTx(ctx context.Context) boil.ContextExecutor {
	tx, ok := ctx.Value(CtxTransactionKey).(*sql.Tx)
	if !ok {
		return r.db
	}

	return tx
}
