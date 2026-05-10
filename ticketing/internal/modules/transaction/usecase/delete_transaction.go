package usecase

import (
	"context"

	"ticketing/internal/modules/transaction/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *transactionUsecaseImpl) DeleteTransaction(ctx context.Context, id int64) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TransactionUsecase:DeleteTransaction")
	defer trace.Finish()

	repoFilter := domain.FilterTransaction{ID: &id}
	return uc.repoSQL.TransactionRepo().Delete(ctx, &repoFilter)
}
