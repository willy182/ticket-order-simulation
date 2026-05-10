package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"ticketing/internal/modules/ticket/domain"

	"github.com/golangid/candi/candishared"
	"github.com/golangid/candi/tracer"
	"github.com/sirupsen/logrus"
)

func (uc *ticketUsecaseImpl) UpdateTicket(ctx context.Context, data *domain.RequestTicket) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TicketUsecase:UpdateTicket")
	defer trace.Finish()

	repoFilter := domain.FilterTicket{ID: &data.ID}
	existing, err := uc.repoSQL.TicketRepo().Find(ctx, &repoFilter)
	if err != nil {
		logrus.Error(err)
		trace.SetError(err)
		return
	}

	existing.Title = data.Title
	existing.Quota = data.Quota
	existing.Price = data.Price
	err = uc.repoSQL.TicketRepo().Save(ctx, &existing, candishared.DBUpdateSetUpdatedFields("Title", "Quota", "Price"))
	if err != nil {
		logrus.Error(err)
		trace.SetError(err)
		return
	}

	var result domain.ResponseTicket
	result.Serialize(&existing)

	byteData, _ := json.Marshal(result)
	keyDataTicket := fmt.Sprintf("data_ticket_%d", result.ID)
	keyQuotaTicket := fmt.Sprintf("kuota_ticket_%d", result.ID)
	uc.cache.Set(ctx, keyDataTicket, string(byteData), -1)
	uc.cache.Set(ctx, keyQuotaTicket, strconv.Itoa(result.Quota), -1)

	return
}
