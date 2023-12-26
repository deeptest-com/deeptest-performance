package exec

import (
	"context"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"log"
	"time"
)

func ExecScenario(valCtx context.Context, stream *proto.PerformanceService_ExecServer, sender MessageSender) {
	task := valCtx.Value("task").(domain.Task)
	log.Println(task)

	for _, processorId := range task.Scenario.Processors {
		log.Println("exec processor", processorId)

		// 此处为场景处理器的耗时操作
		time.Sleep(2 * time.Second)

		result := proto.PerformanceExecResult{
			Uuid:   fmt.Sprintf("%s@%s", processorId, task.Uuid),
			Status: "pass",
		}

		sender.Send(result)

		// 每个场景处理器完成后，检测是否有终止执行的信号
		select {
		case <-valCtx.Done():
			fmt.Println("stop", task.VuNo)

			// 中止执行该场景后续处理器
			return

		default:
		}
	}

	return
}
