package agentService

import (
	"context"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"log"
	"time"
)

type VuService struct {
}

func (s *VuService) Exec(valCtx context.Context) {
	task := valCtx.Value("task").(domain.Task)
	log.Println(task)

	// 此处执行场景
	s.ExecScenario(valCtx)

	return
}

// 模拟执行场景的耗时过程
func (s *VuService) ExecScenario(valCtx context.Context) {
	task := valCtx.Value("task").(domain.Task)
	processors := []string{"p1", "p2", "p3", "p4", "p5", "p6"}

	for _, processor := range processors {
		// 此处为场景处理器的耗时操作
		log.Println("exec processor", processor)
		time.Sleep(3 * time.Second)

		// 每个场景处理器完成后，检测是否有终止执行的信号
		for {
			select {
			case <-valCtx.Done():
				fmt.Println("stop", task.VuNo)
				return

			default:
			}
		}
	}

	return
}
