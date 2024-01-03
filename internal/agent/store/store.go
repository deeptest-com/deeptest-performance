package store

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"sync"
	"time"
)

const (
	keyStartTime = "startTime"
	keyEndTime   = "keyEndTime"
	keyDuration  = "keyAvgDuration"

	keyPassCount  = "pass"
	keyFailCount  = "fail"
	keyErrorCount = "error"

	keyAvgQps      = "qps"
	keyAvgDuration = "duration"

	keyRequests = "requests"
)

var (
	summary sync.Map
)

func Init() {
	UpdateStartTime(time.Now().Unix())
	UpdateEndTime(time.Now().Unix())
	UpdateDuration(0)

	UpdateRequests([]domain.RequestItem{})
	UpdatePass(0)
	UpdateFail(0)
	UpdateError(0)

	UpdateAvgQps(0)
	UpdateAvgDuration(0)
}

func GetData() (ret domain.Stat) {
	ret.Requests = GetRequests()

	ret.Pass = GetPass()
	ret.Pass = GetPass()
	ret.Fail = GetFail()
	ret.Error = GetError()

	ret.AvgQps = GetAvgQps()
	ret.AvgDuration = GetAvgDuration()

	return
}

func GetStartTime() (ret int64) {
	val, ok := summary.Load(keyStartTime)

	if ok {
		ret = val.(int64)
	}

	return
}
func GetEndTime() (ret int64) {
	val, ok := summary.Load(keyEndTime)
	if ok {
		ret = val.(int64)
	}

	return
}
func GetDuration() (ret int) {
	val, ok := summary.Load(keyDuration)
	if ok {
		ret = val.(int)
	}

	return
}

func GetRequests() (ret []domain.RequestItem) {
	val, ok := summary.Load(keyRequests)

	if ok {
		ret = val.([]domain.RequestItem)
	}

	return
}
func GetPass() (ret int) {
	val, ok := summary.Load(keyPassCount)

	if ok {
		ret = val.(int)
	}

	return
}
func GetFail() (ret int) {
	val, ok := summary.Load(keyFailCount)
	if ok {
		ret = val.(int)
	}

	return
}
func GetError() (ret int) {
	val, ok := summary.Load(keyErrorCount)
	if ok {
		ret = val.(int)
	}

	return
}
func GetAvgQps() (ret int) {
	val, ok := summary.Load(keyAvgQps)
	if ok {
		ret = val.(int)
	}

	return
}
func GetAvgDuration() (ret int) {
	val, ok := summary.Load(keyAvgDuration)
	if ok {
		ret = val.(int)
	}

	return
}

func AddRequests(val domain.RequestItem) {
	arr := GetRequests()
	arr = append(arr, val)

	UpdateRequests(arr)
}
func AddPass(count int) {
	old := GetPass()
	UpdatePass(old + count)
}
func AddFail(count int) {
	old := GetPass()
	UpdateFail(old + count)
}
func AddError(count int) {
	old := GetError()
	UpdateError(old + count)
}

func UpdateStartTime(val int64) {
	summary.Store(keyStartTime, val)
}
func UpdateEndTime(val int64) {
	summary.Store(keyEndTime, val)
}
func UpdateDuration(val int64) {
	summary.Store(keyDuration, val)
}

func UpdateRequests(val []domain.RequestItem) {
	summary.Store(keyRequests, val)
}
func UpdatePass(val int) {
	summary.Store(keyPassCount, val)
}
func UpdateFail(val int) {
	summary.Store(keyFailCount, val)
}
func UpdateError(val int) {
	summary.Store(keyErrorCount, val)
}
func UpdateAvgQps(val int) {
	summary.Store(keyAvgQps, val)
}
func UpdateAvgDuration(val int) {
	summary.Store(keyAvgDuration, val)
}
