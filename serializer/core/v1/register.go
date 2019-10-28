package v1

import (
	"github.com/caicloud/rudder-client/serializer/universal"

	corev1 "k8s.io/api/core/v1"
)

var (
	gvkService = corev1.SchemeGroupVersion.WithKind("Service")
)

func Register(p universal.SerializerFactory) {
	p.Register(gvkService, &serviceSerializer{})
}
