package usecase

import (
	"context"
	"fmt"

	"ticketing/internal/modules/ticket/domain"
	"ticketing/pkg/helper"
	shareddomain "ticketing/pkg/shared/domain"

	"github.com/golangid/candi/tracer"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func (uc *ticketUsecaseImpl) GetAllTicket(ctx context.Context, filter *domain.FilterTicket) (results domain.ResponseTicketList, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TicketUsecase:GetAllTicket")
	defer func() {
		trace.Log("filter", filter)
		trace.Log("results", results)
		trace.Finish(tracer.FinishWithError(err))
	}()

	var data []shareddomain.Ticket

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() (errGetData error) {
		defer func() {
			if r := recover(); r != nil {
				errGetData = fmt.Errorf("%v", r)
			}
		}()

		data, errGetData = uc.repoSQL.TicketRepo().FetchAll(egCtx, filter)
		return errGetData
	})

	count := uc.repoSQL.TicketRepo().Count(ctx, filter)
	results.Meta = helper.NewMeta(filter.Page, filter.Limit, int64(count), filter.ShowAll)

	if err = eg.Wait(); err != nil {
		logrus.Error(err)
		trace.SetError(err)
		return
	}

	result := make([]domain.ResponseTicket, 0)
	for _, detail := range data {
		var ticket domain.ResponseTicket
		ticket.Serialize(&detail)
		result = append(result, ticket)
	}
	results.Data = result

	return
}
