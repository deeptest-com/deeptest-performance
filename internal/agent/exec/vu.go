package exec

import (
	"context"
	"github.com/aaronchen2k/deeptest/internal/agent/logs"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"log"
)

func ExecTaskWithVu(valCtx context.Context, stream *proto.PerformanceService_ExecServer) (err error) {
	task := valCtx.Value("task").(domain.Task)
	log.Println(task)

	var sender logs.MessageSender

	if task.NsqServerAddress != "" {
		sender = logs.NewNsqSender(task.Uuid, task.NsqServerAddress, task.NsqLookupAddress)
	} else {
		sender = logs.NewGrpcSender(stream)
	}

	ExecScenario(valCtx, stream, sender)

	return
}
