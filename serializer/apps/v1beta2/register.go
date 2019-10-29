package v1beta2

import (
	"reflect"

	"github.com/caicloud/rudder-client/serializer/universal"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
)

var (
	gvkDeployment = appsv1beta2.SchemeGroupVersion.WithKind(reflect.TypeOf(appsv1beta2.Deployment{}).Name())
)

func Register(p universal.SerializerFactory) {
	p.Register(gvkDeployment, &deploymentSerializer{})
}
