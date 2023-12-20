package domain

import "github.com/aaronchen2k/deeptest/internal/pkg/consts"

type GlobalVar struct {
	VarId       uint                       `gorm:"-" json:"varId"`
	Name        string                     `json:"name"`
	RightValue  string                     `gorm:"type:text" json:"rightValue"`
	LocalValue  string                     `gorm:"type:text" json:"localValue"`
	RemoteValue string                     `gorm:"type:text" json:"remoteValue"`
	ValueType   consts.ExtractorResultType `json:"valueType"`
}
