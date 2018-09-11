package v1

import (
	"github.com/caicloud/rudder-client/status/internal"
	batchv1 "k8s.io/api/batch/v1"
)

var (
	gvkJob = batchv1.SchemeGroupVersion.WithKind("Job")
)

func Assist(u internal.Umpire) {
	u.Employ(gvkJob, JudgeJob)
}
