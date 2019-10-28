package v1

import (
	"github.com/caicloud/rudder-client/serializer/universal"

	appsv1 "k8s.io/api/apps/v1"
)

var (
	gvkDeployment = appsv1.SchemeGroupVersion.WithKind("Deployment")
	gvkDaemonSet  = appsv1.SchemeGroupVersion.WithKind("DaemonSet")
)

// Register register workloads
func Register(p universal.SerializerFactory) {
	p.Register(gvkDeployment, &deploymentSerializer{})
	p.Register(gvkDaemonSet, &daemonSetSerializer{})
}
