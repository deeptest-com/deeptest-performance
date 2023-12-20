package service

import (
	"errors"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	repo "github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type PlanService struct {
	PlanRepo        *repo.PlanRepo        `inject:""`
	PlanReportRepo  *repo.PlanReportRepo  `inject:""`
	UserRepo        *repo.UserRepo        `inject:""`
	EnvironmentRepo *repo.EnvironmentRepo `inject:""`
}

func (s *PlanService) Paginate(req v1.PlanReqPaginate, projectId int) (ret _domain.PageData, err error) {
	ret, err = s.PlanRepo.Paginate(req, projectId)

	if err != nil {
		return
	}

	return
}

func (s *PlanService) GetById(id uint, detail bool) (ret model.Plan, err error) {
	userIds := make([]uint, 0)
	ret, err = s.PlanRepo.Get(id)
	if err != nil {
		return
	}
	userIds = append(userIds, ret.AdminId)
	userIds = append(userIds, ret.UpdateUserId)
	userIds = append(userIds, ret.CreateUserId)

	lastPlanReport, err := s.PlanReportRepo.GetLastByPlanId(id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	planExecTimes, err := s.PlanReportRepo.GetPlanExecNumber(id)
	if err != nil {
		return
	}

	testPassRate := ""
	if lastPlanReport.ID == 0 {
		testPassRate = "n/a"
	} else {
		userIds = append(userIds, lastPlanReport.CreateUserId)

		if lastPlanReport.TotalScenarioNum == 0 {
			testPassRate = "0%"
		} else {
			testPassRate = strconv.Itoa(lastPlanReport.PassScenarioNum*100/lastPlanReport.TotalScenarioNum) + "%"
		}
	}

	userIdNameMap := s.UserRepo.GetUserIdNameMap(userIds)

	if adminName, ok := userIdNameMap[ret.AdminId]; ok {
		ret.AdminName = adminName
	}
	if updateUserName, ok := userIdNameMap[ret.UpdateUserId]; ok {
		ret.UpdateUserName = updateUserName
	}
	if createUserName, ok := userIdNameMap[ret.CreateUserId]; ok {
		ret.CreateUserName = createUserName
	}
	ret.TestPassRate = testPassRate
	ret.ExecTimes = planExecTimes
	if lastPlanReport.ID != 0 {
		if executorName, ok := userIdNameMap[lastPlanReport.CreateUserId]; ok {
			ret.ExecutorName = executorName
		}
		ret.ExecTime = lastPlanReport.CreatedAt

		environment, _ := s.EnvironmentRepo.Get(lastPlanReport.ExecEnvId)
		if environment.ID != 0 {
			ret.ExecEnv = environment.Name
		}
	}

	//if detail {
	//	ret.Scenarios, _ = s.PlanRepo.ListScenario(id)
	//}

	return
}

func (s *PlanService) Create(req model.Plan) (po model.Plan, bizErr *_domain.BizErr) {
	po, bizErr = s.PlanRepo.Create(req)

	return
}

func (s *PlanService) Update(req model.Plan) error {
	return s.PlanRepo.Update(req)
}

func (s *PlanService) DeleteById(id uint) error {
	return s.PlanRepo.DeleteById(id)
}

func (s *PlanService) AddScenarios(planId uint, scenarioIds []uint) (err error) {
	err = s.PlanRepo.AddScenarios(planId, scenarioIds)
	return
}

func (s *PlanService) RemoveScenario(planId int, scenarioId int) (err error) {
	err = s.PlanRepo.RemoveScenario(planId, scenarioId)
	return
}

func (s *PlanService) RemoveScenarios(planId int, scenarioIds []uint) (err error) {
	err = s.PlanRepo.RemoveScenarios(planId, scenarioIds)
	return
}

func (s *PlanService) StatusDropDownOptions() map[consts.TestStatus]string {
	return s.PlanRepo.StatusDropDownOptions()
}

func (s *PlanService) TestStageDropDownOptions() map[consts.TestStage]string {
	return s.PlanRepo.TestStageDropDownOptions()
}

func (s *PlanService) Clone(id, userId uint) (ret model.Plan, err error) {
	plan, err := s.PlanRepo.Get(id)
	if err != nil {
		return
	}

	planScenarioRelations, err := s.PlanRepo.ListScenarioRelation(id)
	if err != nil {
		return
	}
	scenarioIds := make([]uint, 0)
	for _, v := range planScenarioRelations {
		scenarioIds = append(scenarioIds, v.ScenarioId)
	}

	plan.ID = 0
	plan.Name = plan.Name + "-COPY"
	plan.CreateUserId = userId
	plan.UpdateUserId = 0
	now := time.Now()
	plan.CreatedAt = &now
	plan, bizErr := s.PlanRepo.Create(plan)
	if bizErr != nil {
		err = errors.New(bizErr.Msg)
		return
	}

	if len(scenarioIds) > 0 {
		err = s.AddScenarios(plan.ID, scenarioIds)
		if err != nil {
			return ret, err
		}
	}

	ret = plan

	return
}

func (s *PlanService) PlanScenariosPaginate(req v1.PlanScenariosReqPaginate, planId uint) (ret _domain.PageData, err error) {
	return s.PlanRepo.PlanScenariosPaginate(req, planId)
}

func (s *PlanService) NotRelationScenarioList(req v1.NotRelationScenarioReqPaginate, projectId int) (ret _domain.PageData, err error) {
	return s.PlanRepo.NotRelationScenarioList(req, projectId)
}
