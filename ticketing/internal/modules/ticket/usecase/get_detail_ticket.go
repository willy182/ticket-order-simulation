package usecase

import (
	"context"

	"ticketing/internal/modules/ticket/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *ticketUsecaseImpl) GetDetailTicket(ctx context.Context, id int) (result domain.ResponseTicket, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TicketUsecase:GetDetailTicket")
	defer trace.Finish()

	repoFilter := domain.FilterTicket{ID: &id}
	data, err := uc.repoSQL.TicketRepo().Find(ctx, &repoFilter)
	if err != nil {
		return
	}

	result.Serialize(&data)
	return
}
