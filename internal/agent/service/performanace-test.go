package agentService

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/mq"
	"github.com/aaronchen2k/deeptest/proto"
	"io"
)

type PerformanceTestServices struct{}

func (services *PerformanceTestServices) Exec(stream proto.PerformanceService_ExecServer) (err error) {
	i := 0

	for {
		err := stream.Send(&proto.PerformanceExecResult{
			Msg: "Hello, I am Agent",
		})
		if err != nil {
			return err
		}

		//
		res, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if res == nil {
			continue
		}

		err = stream.Send(&proto.PerformanceExecResult{
			Uuid:   res.Uuid,
			Status: "pass",
		})
		if err != nil {
			return err
		}

		i++
	}

	mqData := mq.MqMsg{
		Event: "exit",
	}
	mq.PubMsg(mqData)

	return
}
