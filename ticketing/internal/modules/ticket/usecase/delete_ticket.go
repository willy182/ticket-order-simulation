package usecase

import (
	"context"
	"fmt"

	"ticketing/internal/modules/ticket/domain"

	"github.com/golangid/candi/tracer"
	"github.com/sirupsen/logrus"
)

func (uc *ticketUsecaseImpl) DeleteTicket(ctx context.Context, id int) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TicketUsecase:DeleteTicket")
	defer trace.Finish()

	repoFilter := domain.FilterTicket{ID: &id}
	err = uc.repoSQL.TicketRepo().Delete(ctx, &repoFilter)
	if err != nil {
		logrus.Error(err)
		trace.SetError(err)
		return
	}

	keyDataTicket := fmt.Sprintf("data_ticket_%d", id)
	keyQuotaTicket := fmt.Sprintf("kuota_ticket_%d", id)
	uc.cache.Delete(ctx, keyDataTicket)
	uc.cache.Delete(ctx, keyQuotaTicket)

	return
}
