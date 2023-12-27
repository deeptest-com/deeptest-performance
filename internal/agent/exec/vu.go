package exec

import (
	"context"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"log"
)

func ExecTaskWithVu(valCtx context.Context, stream *proto.PerformanceService_ExecServer) (err error) {
	task := valCtx.Value("task").(domain.Task)
	log.Println(task)

	var sender MessageSender

	if task.NsqServerAddress != "" {
		sender = NewNsqSender(task.Uuid, task.NsqServerAddress, task.NsqLookupAddress)
	} else {
		sender = NewGrpcSender(stream)
	}

	ExecScenario(valCtx, stream, sender)

	return
}
