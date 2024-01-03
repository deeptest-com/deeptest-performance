package logs

import (
	"github.com/aaronchen2k/deeptest/internal/agent/store"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"time"
)

func Count(result proto.PerformanceExecResult) (err error) {
	status := consts.ResultStatus(result.Status)

	// deal with request
	request := domain.RequestItem{
		Status: status,
		Dur:    int(result.Duration),
	}
	store.AddRequests(request)

	// deal with status
	if status == consts.Pass {
		store.AddPass(1)
	} else if status == consts.Fail {
		store.AddFail(1)
	} else if status == consts.Error {
		store.AddError(1)
	}

	currTime := time.Now().Unix()
	store.UpdateEndTime(currTime)
	store.UpdateDuration(currTime - store.GetStartTime())

	// count average duration
	avgDuration := computeAvgDuration(result.Duration)
	store.UpdateAvgDuration(avgDuration)

	// count average qps, must put after all other actions
	avgQps := computeAvgQps()
	store.UpdateAvgQps(avgQps)

	return
}

func computeAvgDuration(requestDur int64) (ret int) {
	requestNum := len(store.GetRequests())
	oldVal := store.GetAvgDuration()

	newVal := (oldVal*requestNum+int(requestDur))/requestNum + 1
	store.UpdateAvgDuration(newVal)

	return
}

func computeAvgQps() (ret int) {
	successNum := store.GetPass()

	duration := store.GetDuration()

	secs := duration / 1000
	qps := successNum / secs

	store.UpdateAvgDuration(qps)

	return
}
