package agentService

import (
	"context"
	"github.com/aaronchen2k/deeptest/internal/agent/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"io"
	"sync"
)

type PerformanceTestServices struct {
}

func (s *PerformanceTestServices) Exec(stream proto.PerformanceService_ExecServer) (err error) {
	plan, err := stream.Recv()
	if err == io.EOF {
		err = nil
		return
	}
	if plan == nil {
		return
	}

	// simulate execution

	ctx, cancel := context.WithCancel(context.Background())

	planCtx := context.WithValue(ctx, "plan", plan)
	go exec.Monitor(&stream, planCtx)

	var wg sync.WaitGroup

	for i := int32(1); i <= plan.Vus; i++ {
		task := domain.Task{
			Uuid:     plan.Uuid,
			Vus:      int(plan.Vus),
			Dur:      int(plan.Vus),
			VuNo:     int(i),
			Scenario: plan.Scenarios,

			NsqServerAddress: plan.NsqServerAddress,
			NsqLookupAddress: plan.NsqLookupAddress,
		}

		timeoutCtx, _ := context.WithTimeout(ctx, consts.ExecTimeout)
		taskCtx := context.WithValue(timeoutCtx, "task", task)

		wg.Add(1)
		go func() {
			defer wg.Done()
			exec.ExecTaskWithVu(taskCtx, &stream)
		}()
	}

	// 等待所有虚拟用户执行结束
	wg.Wait()

	// 模拟结束
	// send stop instruction
	stopMsg := proto.PerformanceExecResult{
		Instruction: consts.Exit.String(),
		Msg:         "exit test",
	}
	sender := exec.NewGrpcSender(&stream)
	sender.Send(stopMsg)

	cancel()

	return
}

func (s *PerformanceTestServices) ForwardResult(result proto.PerformanceExecResult, stream *proto.PerformanceService_ExecServer) (err error) {
	err = (*stream).Send(&result)

	return
}
