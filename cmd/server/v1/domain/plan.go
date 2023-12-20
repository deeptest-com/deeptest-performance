package serverDomain

import (
	"github.com/aaronchen2k/deeptest/pkg/domain"
)

type PlanReqPaginate struct {
	_domain.PaginateReq

	ProjectId  uint   `json:"projectId"`
	CategoryId int64  `json:"categoryId"`
	Status     string `json:"status"`
	AdminId    string `json:"adminId"`
	Keywords   string `json:"keywords"`
	Enabled    string `json:"enabled"`
}
