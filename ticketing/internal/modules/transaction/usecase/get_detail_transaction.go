package usecase

import (
	"context"

	"ticketing/internal/modules/transaction/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *transactionUsecaseImpl) GetDetailTransaction(ctx context.Context, id int64) (result domain.ResponseTransaction, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TransactionUsecase:GetDetailTransaction")
	defer trace.Finish()

	repoFilter := domain.FilterTransaction{ID: &id}
	data, err := uc.repoSQL.TransactionRepo().Find(ctx, &repoFilter)
	if err != nil {
		return result, err
	}

	result.Serialize(&data)
	return
}
