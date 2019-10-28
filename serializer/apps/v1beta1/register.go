package v1beta1

import (
	"github.com/caicloud/rudder-client/serializer/universal"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
)

var (
	gvkDeployment = appsv1beta1.SchemeGroupVersion.WithKind("Deployment")
)

func Register(p universal.SerializerFactory) {
	p.Register(gvkDeployment, &deploymentSerializer{})
}
