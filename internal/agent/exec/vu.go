package exec

import (
	"context"
	"github.com/aaronchen2k/deeptest/internal/agent/logs"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"log"
	"time"
)

func ExecTaskWithVu(valCtx context.Context, stream *proto.PerformanceService_ExecServer) (err error) {
	task := valCtx.Value("task").(domain.Task)
	log.Println(task)

	var sender logs.MessageSender

	if task.NsqServerAddress != "" {
		sender = logs.GetNsqSenderInstant(task.Uuid, task.NsqServerAddress, task.NsqLookupAddress)
	} else {
		sender = logs.NewGrpcSender(stream)
	}

	taskTimeoutCtx, cancel := context.WithTimeout(valCtx, time.Second*time.Duration(task.Dur))
	defer cancel()

	for true {
		ExecScenario(taskTimeoutCtx, stream, sender)

		select {
		case <-taskTimeoutCtx.Done():
			goto Label_END_TASK

		default:
		}
	}

Label_END_TASK:

	return
}
