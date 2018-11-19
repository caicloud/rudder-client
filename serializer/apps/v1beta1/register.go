package v1beta1

import (
	appsv1beta1 "k8s.io/api/apps/v1beta1"

	"github.com/caicloud/rudder-client/serializer/universal"
)

var (
	gvkDeployment = appsv1beta1.SchemeGroupVersion.WithKind("Deployment")
)

func Register(p universal.SerializerFactory) {
	p.Register(gvkDeployment, &deploymentSerializer{})
}
