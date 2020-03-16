package v1

import (
	"github.com/caicloud/clientset/informers"
	releaseapi "github.com/caicloud/clientset/pkg/apis/release/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
)

func JudgeSVC(informerFactory informers.SharedInformerFactory, obj runtime.Object) (releaseapi.ResourceStatus, error) {
	return releaseapi.ResourceStatusFrom(releaseapi.ResourceRunning), nil
}
