package v1beta2

import (
	"github.com/caicloud/rudder-client/serializer/universal"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
)

var (
	gvkDeployment = appsv1beta2.SchemeGroupVersion.WithKind("Deployment")
)

func Register(p universal.SerializerFactory) {
	p.Register(gvkDeployment, &deploymentSerializer{})
}
