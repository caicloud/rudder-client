package universal

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Serializer interface {
	Encode(obj runtime.Object, chart string, cur int, fn func(runtime.Object) (runtime.Object, error)) (string, string, error)
}

type SerializerFactory interface {
	Register(gvk schema.GroupVersionKind, serializer Serializer)
	SerializerFor(gvk schema.GroupVersionKind) (Serializer, error)
}
