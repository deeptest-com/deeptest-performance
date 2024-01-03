package domain

import "github.com/aaronchen2k/deeptest/internal/pkg/consts"

type Stat struct {
	Requests []RequestItem `json:"requests"`
	Pass     int           `json:"pass"`
	Fail     int           `json:"fail"`
	Error    int           `json:"error"`

	AvgQps      int `json:"avgQps"`
	AvgDuration int `json:"avgDuration"`
}

type RequestItem struct {
	Status consts.ResultStatus `json:"status"`
	Dur    int                 `json:"dur"`
}
