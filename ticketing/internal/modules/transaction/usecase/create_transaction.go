package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	domainTicket "ticketing/internal/modules/ticket/domain"
	"ticketing/internal/modules/transaction/domain"
	"ticketing/pkg/helper"

	"github.com/golangid/candi/candihelper"
	taskqueueworker "github.com/golangid/candi/codebase/app/task_queue_worker"
	"github.com/golangid/candi/tracer"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

func (uc *transactionUsecaseImpl) CreateTransaction(ctx context.Context, req *domain.RequestTransaction) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TransactionUsecase:CreateTransaction")
	defer trace.Finish()

	durInterval, _ := time.ParseDuration("3s")

	_, err = taskqueueworker.AddJob(ctx, &taskqueueworker.AddJobRequest{
		TaskName:      helper.ReserveTicket,
		MaxRetry:      3,
		RetryInterval: durInterval,
		Args:          candihelper.ToBytes(&req),
	})
	if err != nil {
		logrus.Error(err)
		trace.SetError(err)
		return
	}

	return
}

func (uc *transactionUsecaseImpl) ReserveTicket(ctx context.Context, req domain.RequestTransaction) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "TransactionUsecase:ReserveTicket")
	defer trace.Finish()

	keyDataTicket := fmt.Sprintf("data_ticket_%d", req.TicketID)
	keyQuotaTicket := fmt.Sprintf("kuota_ticket_%d", req.TicketID)

	var dataByteTicket []byte
	dataByteTicket, err = uc.cache.Get(ctx, keyDataTicket)
	if err != nil {
		logrus.Error(err)
		trace.SetError(err)
		return
	}

	var dataCacheTicket domainTicket.ResponseTicket
	err = json.Unmarshal(dataByteTicket, &dataCacheTicket)
	if err != nil {
		logrus.Error(err)
		trace.SetError(err)
		return
	}

	switch req.Status {
	case helper.FAILED:
		data := req.Deserialize()
		err = uc.repoSQL.TransactionRepo().Save(ctx, &data)
		if err != nil {
			logrus.Error(err)
			trace.SetError(err)
			return
		}

		durInterval, _ := time.ParseDuration("3s")
		reqSendEmail := domain.ReqSendEmail{
			CustomerName:  req.CustomerName,
			CustomerEmail: req.CustomerEmail,
			CustomerPhone: req.CustomerPhone,
			Status:        req.Status,
			TicketTitle:   dataCacheTicket.Title,
		}

		// send job email
		_, err = taskqueueworker.AddJob(ctx, &taskqueueworker.AddJobRequest{
			TaskName:      helper.SendEmail,
			MaxRetry:      3,
			RetryInterval: durInterval,
			Args:          candihelper.ToBytes(reqSendEmail),
		})
	case helper.SUCCESS:
		// ======================================================
		// LUA SCRIPT
		// ======================================================
		const reserveTicketLua = `
		local stock = tonumber(redis.call("GET", KEYS[1]))
		local qty = tonumber(ARGV[1])

		-- key tidak ditemukan
		if not stock then
			return -2
		end

		-- stock tidak cukup
		if stock < qty then
			return -1
		end

		-- kurangi stock
		return redis.call("DECRBY", KEYS[1], qty)
		`

		// eksekusi LUA script untuk mengurangi kuota tiket
		var resExec int
		resExec, err = redis.Int(uc.cache.DoCommand(ctx, true, "EVAL", reserveTicketLua, 1, keyQuotaTicket, req.Qty))
		if err != nil {
			logrus.Error(err)
			trace.SetError(err)
			return
		}

		switch resExec {
		case -2:
			err = errors.New("ticket stock not found")
			return

		case -1:
			err = errors.New("ticket sold out")
			return
		}

		req.InvoiceNumber = candihelper.ToStringPtr(helper.GenerateInvoice())
		req.TotalAmount = float64(req.Qty) * dataCacheTicket.Price
		data := req.Deserialize()
		err = uc.repoSQL.TransactionRepo().Save(ctx, &data)
		if err != nil {
			logrus.Error(err)
			trace.SetError(err)
			return
		}

		var result domain.ResponseTransaction
		result.Serialize(&data)

		durInterval, _ := time.ParseDuration("3s")

		// send job generate ticket code
		_, err = taskqueueworker.AddJob(ctx, &taskqueueworker.AddJobRequest{
			TaskName:      helper.TaskGenerateTicketCode,
			MaxRetry:      3,
			RetryInterval: durInterval,
			Args:          candihelper.ToBytes(result.ID),
		})
	case helper.PENDING:
		data := req.Deserialize()
		err = uc.repoSQL.TransactionRepo().Save(ctx, &data)
	}

	return
}
