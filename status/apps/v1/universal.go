package v1

import (
	"fmt"
	"sort"

	"github.com/caicloud/clientset/informers"
	releaseapi "github.com/caicloud/clientset/pkg/apis/release/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

func getPodsFor(informerFactory informers.SharedInformerFactory, obj runtime.Object) ([]*corev1.Pod, error) {
	var selector labels.Selector
	var namespace string
	var err error
	switch resource := obj.(type) {
	case *appsv1.Deployment:
		namespace = resource.Namespace
		selector, err = metav1.LabelSelectorAsSelector(resource.Spec.Selector)
	case *appsv1.DaemonSet:
		namespace = resource.Namespace
		selector, err = metav1.LabelSelectorAsSelector(resource.Spec.Selector)
	case *appsv1.StatefulSet:
		namespace = resource.Namespace
		selector, err = metav1.LabelSelectorAsSelector(resource.Spec.Selector)
	default:
		return nil, fmt.Errorf("getPodsFor: %v is not supported", obj)
	}

	if err != nil {
		return nil, fmt.Errorf("invalid label selector: %v", err)
	}
	// If a resource with a nil or empty selector creeps in, it should match nothing, not everything.
	if selector.Empty() {
		return nil, nil
	}
	return informerFactory.Native().Core().V1().Pods().Lister().Pods(namespace).List(selector)
}

func getLatestEventFor(kind string, obj metav1.Object, events []*corev1.Event) *corev1.Event {
	if len(events) == 0 {
		return nil
	}
	ret := make([]*corev1.Event, 0)
	for _, e := range events {
		if e.InvolvedObject.Kind == kind &&
			e.InvolvedObject.Name == obj.GetName() &&
			e.InvolvedObject.Namespace == obj.GetNamespace() &&
			e.InvolvedObject.UID == obj.GetUID() {
			ret = append(ret, e)
		}
	}
	if len(ret) == 0 {
		return nil
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].LastTimestamp.After(ret[j].LastTimestamp.Time)
	})
	return ret[0]
}

func getPodStatistics(updated []HyperPod, old []HyperPod) *releaseapi.PodStatistics {
	if len(updated) == 0 && len(old) == 0 {
		return nil
	}

	ret := releaseapi.PodStatistics{
		UpdatedPods: make(releaseapi.PodStatusCounter, len(updated)),
		OldPods:     make(releaseapi.PodStatusCounter, len(old)),
	}
	for _, pod := range updated {
		ret.UpdatedPods[pod.Status.Phase]++
	}

	for _, pod := range old {
		ret.OldPods[pod.Status.Phase]++
	}

	return &ret
}

func getLabel(obj metav1.Object, key string) string {
	labels := obj.GetLabels()
	if labels == nil {
		return ""
	}
	return labels[key]
}
