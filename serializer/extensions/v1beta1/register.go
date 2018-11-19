package v1beta1

import (
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"

	"github.com/caicloud/rudder-client/serializer/universal"
)

var (
	gvkDeployment = extensionsv1beta1.SchemeGroupVersion.WithKind("Deployment")
)

func Register(p universal.SerializerFactory) {
	p.Register(gvkDeployment, &deploymentSerializer{})
}
