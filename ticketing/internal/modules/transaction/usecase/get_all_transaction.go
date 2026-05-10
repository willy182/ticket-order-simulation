package usecase

import (
	"context"

	"ticketing/internal/modules/transaction/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"
)

func (uc *transactionUsecaseImpl) GetAllTransaction(ctx context.Context, filter *domain.FilterTransaction) (result domain.ResponseTransactionList, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TransactionUsecase:GetAllTransaction")
	defer trace.Finish()

	data, err := uc.repoSQL.TransactionRepo().FetchAll(ctx, filter)
	if err != nil {
		return result, err
	}
	count := uc.repoSQL.TransactionRepo().Count(ctx, filter)
	result.Meta = candishared.NewMeta(filter.Page, filter.Limit, count)

	result.Data = make([]domain.ResponseTransaction, len(data))
	for i, detail := range data {
		result.Data[i].Serialize(&detail)
	}

	return
}
