/*
Copyright 2019 caicloud authors. All rights reserved.
*/

// Code generated by listerfactory-gen. DO NOT EDIT.

package v1

import (
	internalinterfaces "github.com/caicloud/clientset/listerfactory/internalinterfaces"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubernetes "k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/core/v1"
)

var _ v1.LimitRangeLister = &limitRangeLister{}

var _ v1.LimitRangeNamespaceLister = &limitRangeNamespaceLister{}

// limitRangeLister implements the LimitRangeLister interface.
type limitRangeLister struct {
	client           kubernetes.Interface
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewLimitRangeLister returns a new LimitRangeLister.
func NewLimitRangeLister(client kubernetes.Interface) v1.LimitRangeLister {
	return NewFilteredLimitRangeLister(client, nil)
}

func NewFilteredLimitRangeLister(client kubernetes.Interface, tweakListOptions internalinterfaces.TweakListOptionsFunc) v1.LimitRangeLister {
	return &limitRangeLister{
		client:           client,
		tweakListOptions: tweakListOptions,
	}
}

// List lists all LimitRanges in the indexer.
func (s *limitRangeLister) List(selector labels.Selector) (ret []*corev1.LimitRange, err error) {
	listopt := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	if s.tweakListOptions != nil {
		s.tweakListOptions(&listopt)
	}
	list, err := s.client.CoreV1().LimitRanges(metav1.NamespaceAll).List(listopt)
	if err != nil {
		return nil, err
	}
	for i := range list.Items {
		ret = append(ret, &list.Items[i])
	}
	return ret, nil
}

// LimitRanges returns an object that can list and get LimitRanges.
func (s *limitRangeLister) LimitRanges(namespace string) v1.LimitRangeNamespaceLister {
	return limitRangeNamespaceLister{client: s.client, tweakListOptions: s.tweakListOptions, namespace: namespace}
}

// limitRangeNamespaceLister implements the LimitRangeNamespaceLister
// interface.
type limitRangeNamespaceLister struct {
	client           kubernetes.Interface
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// List lists all LimitRanges in the indexer for a given namespace.
func (s limitRangeNamespaceLister) List(selector labels.Selector) (ret []*corev1.LimitRange, err error) {
	listopt := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	if s.tweakListOptions != nil {
		s.tweakListOptions(&listopt)
	}
	list, err := s.client.CoreV1().LimitRanges(s.namespace).List(listopt)
	if err != nil {
		return nil, err
	}
	for i := range list.Items {
		ret = append(ret, &list.Items[i])
	}
	return ret, nil
}

// Get retrieves the LimitRange from the indexer for a given namespace and name.
func (s limitRangeNamespaceLister) Get(name string) (*corev1.LimitRange, error) {
	return s.client.CoreV1().LimitRanges(s.namespace).Get(name, metav1.GetOptions{})
}
