package v1

import (
	"reflect"

	"github.com/caicloud/rudder-client/serializer/universal"

	corev1 "k8s.io/api/core/v1"
)

var (
	gvkService = corev1.SchemeGroupVersion.WithKind(reflect.TypeOf(corev1.Service{}).Name())
)

func Register(p universal.SerializerFactory) {
	p.Register(gvkService, &serviceSerializer{})
}
