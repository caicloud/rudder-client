package v1beta1

import (
	"github.com/caicloud/rudder-client/status/internal"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
)

var (
	gvkCronJob = batchv1beta1.SchemeGroupVersion.WithKind("CronJob")
)

func Assist(u internal.Umpire) {
	u.Employ(gvkCronJob, JudgeCronJob)
}
