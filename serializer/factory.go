package serializer

import (
	"flag"
	"fmt"
	"sync"

	appsv1 "github.com/caicloud/rudder-client/serializer/apps/v1"
	appsv1beta1 "github.com/caicloud/rudder-client/serializer/apps/v1beta1"
	appsv1beta2 "github.com/caicloud/rudder-client/serializer/apps/v1beta2"
	corev1 "github.com/caicloud/rudder-client/serializer/core/v1"
	extensionsv1beta1 "github.com/caicloud/rudder-client/serializer/extensions/v1beta1"
	"github.com/caicloud/rudder-client/serializer/universal"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() { // TODO: use klog
	flag.Set("logtostderr", "true")
	flag.Set("alsologtostderr", "true")
}

type serializerFactory struct {
	mux         sync.RWMutex
	serializers map[schema.GroupVersionKind]universal.Serializer
}

func NewSerializerFactory() universal.SerializerFactory {
	factory := &serializerFactory{
		serializers: make(map[schema.GroupVersionKind]universal.Serializer),
	}
	factory.register()
	return factory
}

func (sf *serializerFactory) register() {
	appsv1.Register(sf)
	appsv1beta1.Register(sf)
	appsv1beta2.Register(sf)
	corev1.Register(sf)
	extensionsv1beta1.Register(sf)
}

func (sf *serializerFactory) Register(gvk schema.GroupVersionKind, serializer universal.Serializer) {
	sf.mux.Lock()
	defer sf.mux.Unlock()
	sf.serializers[gvk] = serializer
}

func (sf *serializerFactory) SerializerFor(gvk schema.GroupVersionKind) (universal.Serializer, error) {
	s, ok := sf.serializers[gvk]
	if !ok {
		return nil, fmt.Errorf("not found the %s object serializer", gvk.String())
	}
	return s, nil
}
