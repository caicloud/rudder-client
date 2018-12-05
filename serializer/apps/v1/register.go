package v1

import (
	appsv1 "k8s.io/api/apps/v1"

	"github.com/caicloud/rudder-client/serializer/universal"
)

var (
	gvkDeployment = appsv1.SchemeGroupVersion.WithKind("Deployment")
)

func Register(p universal.SerializerFactory) {
	p.Register(gvkDeployment, &deploymentSerializer{})
}
