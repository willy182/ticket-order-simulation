package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"ticketing/internal/modules/ticket/domain"

	"github.com/golangid/candi/tracer"
)

func (uc *ticketUsecaseImpl) CreateTicket(ctx context.Context, req *domain.RequestTicket) (result domain.ResponseTicket, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TicketUsecase:CreateTicket")
	defer trace.Finish()

	data := req.Deserialize()
	err = uc.repoSQL.TicketRepo().Save(ctx, &data)
	if err != nil {
		return
	}

	result.Serialize(&data)

	byteData, _ := json.Marshal(result)
	keyDataTicket := fmt.Sprintf("data_ticket_%d", result.ID)
	keyQuotaTicket := fmt.Sprintf("kuota_ticket_%d", result.ID)
	uc.cache.Set(ctx, keyDataTicket, string(byteData), -1)
	uc.cache.Set(ctx, keyQuotaTicket, strconv.Itoa(result.Quota), -1)

	return
}
