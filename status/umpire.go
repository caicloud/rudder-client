package status

import (
	"fmt"
	"sync"

	"github.com/caicloud/clientset/informers"
	releaseapi "github.com/caicloud/clientset/pkg/apis/release/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	appsv1 "github.com/caicloud/rudder-client/status/apps/v1"
	batchv1 "github.com/caicloud/rudder-client/status/batch/v1"
	batchv1beta1 "github.com/caicloud/rudder-client/status/batch/v1beta1"
	corev1 "github.com/caicloud/rudder-client/status/core/v1"
	"github.com/caicloud/rudder-client/status/universal"
)

type umpire struct {
	informerFactory informers.SharedInformerFactory
	assistants      map[schema.GroupVersionKind]universal.Assistant
	mux             sync.RWMutex
}

// NewUmpire returns a new status Umpire
func NewUmpire(informerFactory informers.SharedInformerFactory) universal.Umpire {
	u := umpire{
		informerFactory: informerFactory,
		assistants:      make(map[schema.GroupVersionKind]universal.Assistant),
	}
	u.employ()

	return &u
}

func (u *umpire) employ() {
	appsv1.Assist(u)
	batchv1.Assist(u)
	batchv1beta1.Assist(u)
	corev1.Assist(u)
}

// Employ employs an assistant for specified object kind.
func (u *umpire) Employ(gvk schema.GroupVersionKind, assistant universal.Assistant) {
	u.mux.Lock()
	defer u.mux.Unlock()
	u.assistants[gvk] = assistant
}

// Judge judges the object and generates a sentence.
func (u *umpire) Judge(obj runtime.Object) (releaseapi.ResourceStatus, error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	assistant, ok := u.assistants[gvk]
	if !ok {
		return releaseapi.ResourceStatusFrom(""), fmt.Errorf("can't find an assistant for: %s", gvk.String())
	}
	return assistant(u.informerFactory, obj)
}
