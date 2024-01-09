package exec

import (
	"context"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/agent/logs"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	_httpUtils "github.com/aaronchen2k/deeptest/pkg/lib/http"
	_intUtils "github.com/aaronchen2k/deeptest/pkg/lib/int"
	"github.com/aaronchen2k/deeptest/proto"
	"log"
	"time"
)

func ExecScenario(valCtx context.Context, stream *proto.PerformanceService_ExecServer, sender logs.MessageSender) {
	task := valCtx.Value("task").(domain.Task)

	for index, processor := range task.Scenario[0].Processors {
		log.Println("exec processor", processor)

		{
			_, err := _httpUtils.Get("http://111.231.16.35:9000/get")
			if err != nil {
				log.Println("http request failed", err)
			}
		}

		duration := _intUtils.GenRandNum(100, 2000)
		time.Sleep(time.Duration(duration) * time.Millisecond)

		status := "pass"
		if index%3 == 0 {
			status = "fail"
		}

		record := proto.PerformanceExecRecord{
			Uuid:     fmt.Sprintf("%s@%s", processor.Name, task.Uuid),
			Duration: int64(duration), // 毫秒
			Status:   status,
		}
		summary := logs.GetSummary()

		result := proto.PerformanceExecResult{
			Instruction: "",
			Record:      &record,
			Summary:     &summary,
		}

		sender.Send(result)

		// 每个场景处理器完成后，检测是否有终止执行的信号
		select {
		case <-valCtx.Done():
			fmt.Println("exit left scenario processors by signal", task.VuNo)

			// 中止执行该场景后续处理器
			goto Label_END_SCENARIO

		default:
		}
	}

	fmt.Println("complete scenario normally", task.VuNo)

Label_END_SCENARIO:

	summary := logs.GetSummary()

	result := proto.PerformanceExecResult{
		Instruction: "",
		Summary:     &summary,
	}

	sender.Send(result)

	return
}
