package generator

import (
	"context"
	"github.com/aaronchen2k/deeptest/internal/agent/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"sync"
	"time"
)

func GenerateRamp(tmplTask domain.Task, stages []*proto.Stage, stream proto.PerformanceService_ExecServer,
	planCtx context.Context, wg *sync.WaitGroup) (err error) {

	if len(stages) != 1 {
		return
	}

	index := 0
	for i := 1; i <= len(stages); i++ {
		stage := stages[i]
		target := stage.Target
		dur := stage.Dur
		startTime := time.Now().Unix()

		for j := int32(1); j <= stage.Target; j++ {
			task := genTask(tmplTask, index)

			timeoutCtx, _ := context.WithTimeout(planCtx, consts.ExecTimeout)
			taskCtx := context.WithValue(timeoutCtx, "task", task)

			wg.Add(1)
			go func() {
				defer wg.Done()
				exec.ExecTaskWithVu(taskCtx, &stream)
			}()

			index++

			leftVus := target - j - 1
			leftTime := getLeftTime(startTime, dur)
			waitTime(int64(leftVus), leftTime)
		}
	}

	return
}

func getLeftTime(startTime int64, dur int32) (leftTime int64) {
	currTime := time.Now().Unix()
	leftTime = int64(dur) - (currTime - startTime)

	if leftTime < 0 {
		leftTime = 0
	}

	return
}

func waitTime(leftVus, leftTime int64) (err error) {
	time.Sleep(time.Duration(leftTime/leftVus) * time.Second)

	return
}
