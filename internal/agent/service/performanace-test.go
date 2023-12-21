package agentService

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/mq"
	"github.com/aaronchen2k/deeptest/proto"
	"io"
)

type PerformanceTestServices struct{}

func (s *PerformanceTestServices) Exec(stream proto.PerformanceService_ExecServer) (err error) {
	go mq.SubAgentMsgWithStream(s.ForwardResult, &stream)

	res, err := stream.Recv()
	if err == io.EOF {
		err = nil
		return
	}
	if res == nil {
		return
	}

	// simulate execution
	i := 0
	for true {
		if i > 2 {
			break
		}

		result := proto.PerformanceExecResult{
			Uuid:   res.Uuid,
			Status: "pass",
		}

		//mqData := mq.MqMsg{
		//	Event:  "result",
		//	Result: result,
		//}
		//mq.PubAgentMsg(mqData)
		err = stream.Send(&result)

		i++
	}

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
