package usecase

import (
	"context"
	"time"

	"ticketing/internal/modules/transaction/domain"
	"ticketing/pkg/helper"

	"github.com/golangid/candi/candihelper"
	"github.com/golangid/candi/candishared"
	taskqueueworker "github.com/golangid/candi/codebase/app/task_queue_worker"
	"github.com/golangid/candi/tracer"
	"github.com/sirupsen/logrus"
)

func (uc *transactionUsecaseImpl) UpdateTransaction(ctx context.Context, data *domain.RequestTransaction) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TransactionUsecase:UpdateTransaction")
	defer trace.Finish()

	repoFilter := domain.FilterTransaction{ID: &data.ID}
	existing, err := uc.repoSQL.TransactionRepo().Find(ctx, &repoFilter)
	if err != nil {
		logrus.Error(err)
		trace.SetError(err)
		return
	}
	existing.Status = data.Status
	err = uc.repoSQL.TransactionRepo().Save(ctx, &existing, candishared.DBUpdateSetUpdatedFields("Status"))
	return
}

func (uc *transactionUsecaseImpl) SendEmail(ctx context.Context, req domain.ReqSendEmail) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TransactionUsecase:SendEmail")
	defer trace.Finish()

	logrus.Infof("Dear %s, your transaction in purchasing ticket %s is %s", req.CustomerName, req.TicketTitle, req.Status)

	return
}

func (uc *transactionUsecaseImpl) GenerateTicketCode(ctx context.Context, id int64) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TransactionUsecase:GenerateTicketCode")
	defer trace.Finish()

	existing, err := uc.repoSQL.TransactionRepo().Find(ctx, &domain.FilterTransaction{ID: &id, Preloads: []string{"TicketData"}})
	if err != nil {
		logrus.Error(err)
		trace.SetError(err)
		return
	}

	existing.TicketCode = candihelper.ToStringPtr(helper.GenerateTicketCode(8))
	err = uc.repoSQL.TransactionRepo().Save(ctx, &existing, candishared.DBUpdateSetUpdatedFields("TicketCode"))
	if err != nil {
		logrus.Error(err)
		trace.SetError(err)
		return
	}

	durInterval, _ := time.ParseDuration("3s")
	reqSendEmail := domain.ReqSendEmail{
		CustomerName:  existing.CustomerName,
		CustomerEmail: existing.CustomerEmail,
		CustomerPhone: existing.CustomerPhone,
		Status:        existing.Status,
		TicketTitle:   existing.TicketData.Title, // Assuming TicketCode contains the ticket title
	}

	// send job email
	_, err = taskqueueworker.AddJob(ctx, &taskqueueworker.AddJobRequest{
		TaskName:      helper.SendEmail,
		MaxRetry:      3,
		RetryInterval: durInterval,
		Args:          candihelper.ToBytes(reqSendEmail),
	})
	return
}
