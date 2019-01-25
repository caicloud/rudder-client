/*
Copyright 2019 caicloud authors. All rights reserved.
*/

// Code generated by listerfactory-gen. DO NOT EDIT.

package v1beta1

import (
	internalinterfaces "github.com/caicloud/clientset/listerfactory/internalinterfaces"
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubernetes "k8s.io/client-go/kubernetes"
	v1beta1 "k8s.io/client-go/listers/policy/v1beta1"
)

var _ v1beta1.PodDisruptionBudgetLister = &podDisruptionBudgetLister{}

var _ v1beta1.PodDisruptionBudgetNamespaceLister = &podDisruptionBudgetNamespaceLister{}

// podDisruptionBudgetLister implements the PodDisruptionBudgetLister interface.
type podDisruptionBudgetLister struct {
	client           kubernetes.Interface
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewPodDisruptionBudgetLister returns a new PodDisruptionBudgetLister.
func NewPodDisruptionBudgetLister(client kubernetes.Interface) v1beta1.PodDisruptionBudgetLister {
	return NewFilteredPodDisruptionBudgetLister(client, nil)
}

func NewFilteredPodDisruptionBudgetLister(client kubernetes.Interface, tweakListOptions internalinterfaces.TweakListOptionsFunc) v1beta1.PodDisruptionBudgetLister {
	return &podDisruptionBudgetLister{
		client:           client,
		tweakListOptions: tweakListOptions,
	}
}

// List lists all PodDisruptionBudgets in the indexer.
func (s *podDisruptionBudgetLister) List(selector labels.Selector) (ret []*policyv1beta1.PodDisruptionBudget, err error) {
	listopt := v1.ListOptions{
		LabelSelector: selector.String(),
	}
	if s.tweakListOptions != nil {
		s.tweakListOptions(&listopt)
	}
	list, err := s.client.PolicyV1beta1().PodDisruptionBudgets(v1.NamespaceAll).List(listopt)
	if err != nil {
		return nil, err
	}
	for i := range list.Items {
		ret = append(ret, &list.Items[i])
	}
	return ret, nil
}

func (s *podDisruptionBudgetLister) GetPodPodDisruptionBudgets(*corev1.Pod) ([]*policyv1beta1.PodDisruptionBudget, error) {
	return nil, nil
}

// PodDisruptionBudgets returns an object that can list and get PodDisruptionBudgets.
func (s *podDisruptionBudgetLister) PodDisruptionBudgets(namespace string) v1beta1.PodDisruptionBudgetNamespaceLister {
	return podDisruptionBudgetNamespaceLister{client: s.client, tweakListOptions: s.tweakListOptions, namespace: namespace}
}

// podDisruptionBudgetNamespaceLister implements the PodDisruptionBudgetNamespaceLister
// interface.
type podDisruptionBudgetNamespaceLister struct {
	client           kubernetes.Interface
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// List lists all PodDisruptionBudgets in the indexer for a given namespace.
func (s podDisruptionBudgetNamespaceLister) List(selector labels.Selector) (ret []*policyv1beta1.PodDisruptionBudget, err error) {
	listopt := v1.ListOptions{
		LabelSelector: selector.String(),
	}
	if s.tweakListOptions != nil {
		s.tweakListOptions(&listopt)
	}
	list, err := s.client.PolicyV1beta1().PodDisruptionBudgets(s.namespace).List(listopt)
	if err != nil {
		return nil, err
	}
	for i := range list.Items {
		ret = append(ret, &list.Items[i])
	}
	return ret, nil
}

// Get retrieves the PodDisruptionBudget from the indexer for a given namespace and name.
func (s podDisruptionBudgetNamespaceLister) Get(name string) (*policyv1beta1.PodDisruptionBudget, error) {
	return s.client.PolicyV1beta1().PodDisruptionBudgets(s.namespace).Get(name, v1.GetOptions{})
}
