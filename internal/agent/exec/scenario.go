package exec

import (
	"context"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"log"
	"time"
)

func ExecScenario(scenario domain.Scenario, valCtx context.Context, stream *proto.PerformanceService_ExecServer) {
	task := valCtx.Value("task").(domain.Task)
	log.Println(task)

	for _, processorId := range scenario.Processors {
		log.Println("exec processor", processorId)

		// 此处为场景处理器的耗时操作
		time.Sleep(2 * time.Second)

		result := proto.PerformanceExecResult{
			Uuid:   fmt.Sprintf("%s@%s", processorId, task.Uuid),
			Status: "pass",
		}

		//mqData := mq.MqMsg{
		//	Event:  "result",
		//	Result: result,
		//}
		//mq.PubAgentMsg(mqData)
		(*stream).Send(&result)

		// 每个场景处理器完成后，检测是否有终止执行的信号

		select {
		case <-valCtx.Done():
			fmt.Println("stop", task.VuNo)
			return

		default:
		}
	}

	return
}
