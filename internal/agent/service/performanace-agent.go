package agentService

import (
	"context"
	"github.com/aaronchen2k/deeptest/internal/agent/generator"
	"github.com/aaronchen2k/deeptest/internal/agent/logs"
	"github.com/aaronchen2k/deeptest/internal/agent/monitor"
	"github.com/aaronchen2k/deeptest/internal/agent/store"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"github.com/jinzhu/copier"
	"io"
	"log"
	"sync"
)

type PerformanceTestServices struct {
}

func (s *PerformanceTestServices) Exec(stream proto.PerformanceService_ExecServer) (err error) {
	store.Init()

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
	go monitor.Monitor(&stream, planCtx)

	var wg sync.WaitGroup

	tmplTask := domain.Task{
		Uuid:     plan.Uuid,
		Stages:   plan.Stages,
		Scenario: plan.Scenarios,

		NsqServerAddress: plan.NsqServerAddress,
		NsqLookupAddress: plan.NsqLookupAddress,
	}

	if plan.GenerateType == consts.GeneratorConstant.String() {
		generator.GenerateConstant(tmplTask, plan.Stages, stream, planCtx, &wg)
	} else {
		generator.GenerateRamp(tmplTask, plan.Stages, stream, planCtx, &wg)
	}

	// 等待所有虚拟用户执行结束
	wg.Wait()

	// 模拟结束
	// send stop instruction
	stopMsg := proto.PerformanceExecResult{
		Instruction: consts.Exit.String(),
		Msg:         "exit test",
	}
	sender := logs.NewGrpcSender(&stream)
	sender.Send(stopMsg)

	cancel()

	// print summary
	data := store.GetData()
	log.Println(data)

	return
}

func (s *PerformanceTestServices) ForwardResult(result proto.PerformanceExecResult, stream *proto.PerformanceService_ExecServer) (err error) {
	err = (*stream).Send(&result)

	return
}

func (s *PerformanceTestServices) GenTask(tmplTask domain.Task, vuNo int) (task domain.Task) {
	copier.CopyWithOption(&task, tmplTask, copier.Option{DeepCopy: true})

	task.VuNo = vuNo

	return
}
