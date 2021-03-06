package v1beta1

import (
	"fmt"
	"sort"

	statusbatchv1 "github.com/caicloud/rudder-client/status/batch/v1"

	"github.com/caicloud/clientset/informers"
	releaseapi "github.com/caicloud/clientset/pkg/apis/release/v1alpha1"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

func JudgeCronJob(informerFactory informers.SharedInformerFactory, obj runtime.Object) (releaseapi.ResourceStatus, error) {
	cronjob, ok := obj.(*batchv1beta1.CronJob)
	if !ok {
		return releaseapi.ResourceStatusFrom(""), fmt.Errorf("unknown type for CronJob: %s", obj.GetObjectKind().GroupVersionKind().String())
	}
	if cronjob.Spec.Suspend != nil && *cronjob.Spec.Suspend {
		return releaseapi.ResourceStatusFrom(releaseapi.ResourceSuspended), nil
	}
	if len(cronjob.Status.Active) > 0 {
		return releaseapi.ResourceStatus{
			Phase:   releaseapi.ResourceProgressing,
			Reason:  "JobRunning",
			Message: fmt.Sprintf("there are %v jobs are running", len(cronjob.Status.Active)),
		}, nil
	}
	jobList, err := getJobForCronJob(informerFactory, cronjob)
	if err != nil {
		return releaseapi.ResourceStatusFrom(""), err
	}

	if len(jobList) == 0 {
		return releaseapi.ResourceStatusFrom(releaseapi.ResourcePending), nil
	}
	// sorts job by creation time
	sort.Slice(jobList, func(i, j int) bool {
		return jobList[i].CreationTimestamp.After(jobList[j].CreationTimestamp.Time)
	})

	return statusbatchv1.JudgeJob(informerFactory, jobList[0])
}

func getJobForCronJob(informerFactory informers.SharedInformerFactory, cronjob *batchv1beta1.CronJob) ([]*batchv1.Job, error) {
	var ret []*batchv1.Job
	js, err := informerFactory.Native().Batch().V1().Jobs().Lister().Jobs(cronjob.Namespace).List(labels.NewSelector())
	if err != nil {
		return nil, err
	}
	for _, job := range js {
		for _, or := range job.GetOwnerReferences() {
			if or.Kind == "CronJob" && or.Name == cronjob.Name {
				ret = append(ret, job)
				break
			}
		}
	}
	return ret, nil
}
