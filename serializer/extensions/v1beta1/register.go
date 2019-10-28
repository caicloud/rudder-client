package v1beta1

import (
	"reflect"

	"github.com/caicloud/rudder-client/serializer/universal"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
)

var (
	gvkDeployment = extensionsv1beta1.SchemeGroupVersion.WithKind(reflect.TypeOf(extensionsv1beta1.Deployment{}).Name())
)

func Register(p universal.SerializerFactory) {
	p.Register(gvkDeployment, &deploymentSerializer{})
}
