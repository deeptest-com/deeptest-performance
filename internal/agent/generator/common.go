package generator

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/jinzhu/copier"
)

func genTask(tmplTask domain.Task, vuNo int) (task domain.Task) {
	copier.CopyWithOption(&task, tmplTask, copier.Option{DeepCopy: true})

	task.VuNo = vuNo

	return
}
