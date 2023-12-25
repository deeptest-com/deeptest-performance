package agentService

import (
	"context"
	"github.com/aaronchen2k/deeptest/internal/agent/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/mq"
	"github.com/aaronchen2k/deeptest/proto"
	"io"
	"time"
)

type PerformanceTestServices struct {
}

func (s *PerformanceTestServices) Exec(stream proto.PerformanceService_ExecServer) (err error) {
	//go mq.SubAgentMsgWithStream(s.ForwardResult, &stream)

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

	for i := int32(1); i <= plan.Vus; i++ {
		task := domain.Task{
			Uuid: plan.Uuid,
			Vus:  int(plan.Vus),
			Dur:  int(plan.Vus),
			VuNo: int(i),
		}
		taskCtx := context.WithValue(ctx, "task", task)
		timeoutCtx, _ := context.WithTimeout(taskCtx, consts.ExecTimeout)

		go exec.ExecTask(timeoutCtx, &stream)
	}

	time.Sleep(10 * time.Second)
	cancel()

	mqData := mq.MqMsg{
		Event: "exit",
	}
	mq.PubAgentMsg(mqData)

	return
}

func (s *PerformanceTestServices) ForwardResult(mqMsg mq.MqMsg, stream *proto.PerformanceService_ExecServer) (err error) {
	result := mqMsg.Result

	err = (*stream).Send(&result)

	return
}
