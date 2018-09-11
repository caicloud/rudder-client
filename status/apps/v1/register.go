package v1

import (
	"github.com/caicloud/rudder-client/status/internal"
	appsv1 "k8s.io/api/apps/v1"
)

var (
	gvkDeployment  = appsv1.SchemeGroupVersion.WithKind("Deployment")
	gvkReplicaSet  = appsv1.SchemeGroupVersion.WithKind("ReplicaSet")
	gvkStatefulSet = appsv1.SchemeGroupVersion.WithKind("StatefulSet")
	gvkDaemonSet   = appsv1.SchemeGroupVersion.WithKind("DaemonSet")
)

func Assist(u internal.Umpire) {
	u.Employ(gvkDeployment, JudgeDeployment)
	u.Employ(gvkReplicaSet, JudgeReplicaSet)
	u.Employ(gvkStatefulSet, JudgeStatefulSet)
	u.Employ(gvkDaemonSet, JudgeDaemonSet)
}
