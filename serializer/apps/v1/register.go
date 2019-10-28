package v1

import (
	"reflect"

	"github.com/caicloud/rudder-client/serializer/universal"

	appsv1 "k8s.io/api/apps/v1"
)

var (
	gvkDeployment = appsv1.SchemeGroupVersion.WithKind(reflect.TypeOf(appsv1.Deployment{}).Name())
	gvkDaemonSet  = appsv1.SchemeGroupVersion.WithKind(reflect.TypeOf(appsv1.DaemonSet{}).Name())
)

// Register register workloads
func Register(p universal.SerializerFactory) {
	p.Register(gvkDeployment, &deploymentSerializer{})
	p.Register(gvkDaemonSet, &daemonSetSerializer{})
}
