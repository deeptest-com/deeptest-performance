package agentService

import (
	"context"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/mq"
	"github.com/aaronchen2k/deeptest/proto"
	"io"
	"time"
)

type PerformanceTestServices struct {
	VuService *VuService `inject:""`
}

func (s *PerformanceTestServices) Exec(stream proto.PerformanceService_ExecServer) (err error) {
	go mq.SubAgentMsgWithStream(s.ForwardResult, &stream)

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

	for i := int32(1); i <= plan.Vus; i++ {
		task := domain.Task{
			Uuid: plan.Uuid,
			Vus:  int(plan.Vus),
			Dur:  int(plan.Vus),
			VuNo: int(i),
		}
		vCtx := context.WithValue(ctx, "task", task)

		go s.VuService.Exec(vCtx)
	}

	time.Sleep(10 * time.Second)
	cancel()

	//i := 0
	//for true {
	//	if i > 2 {
	//		break
	//	}
	//
	//	result := proto.PerformanceExecResult{
	//		Uuid:   res.Uuid,
	//		Status: "pass",
	//	}
	//
	//	//mqData := mq.MqMsg{
	//	//	Event:  "result",
	//	//	Result: result,
	//	//}
	//	//mq.PubAgentMsg(mqData)
	//	err = stream.Send(&result)
	//
	//	i++
	//}

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
