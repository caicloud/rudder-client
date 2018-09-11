/*
Copyright 2018 caicloud authors. All rights reserved.
*/

// Code generated by listerfactory-gen. DO NOT EDIT.

package v1beta1

import (
	internalinterfaces "github.com/caicloud/clientset/listerfactory/internalinterfaces"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubernetes "k8s.io/client-go/kubernetes"
	v1beta1 "k8s.io/client-go/listers/extensions/v1beta1"
)

var _ v1beta1.IngressLister = &ingressLister{}

var _ v1beta1.IngressNamespaceLister = &ingressNamespaceLister{}

// ingressLister implements the IngressLister interface.
type ingressLister struct {
	client           kubernetes.Interface
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewIngressLister returns a new IngressLister.
func NewIngressLister(client kubernetes.Interface) v1beta1.IngressLister {
	return NewFilteredIngressLister(client, nil)
}

func NewFilteredIngressLister(client kubernetes.Interface, tweakListOptions internalinterfaces.TweakListOptionsFunc) v1beta1.IngressLister {
	return &ingressLister{
		client:           client,
		tweakListOptions: tweakListOptions,
	}
}

// List lists all Ingresses in the indexer.
func (s *ingressLister) List(selector labels.Selector) (ret []*extensionsv1beta1.Ingress, err error) {
	listopt := v1.ListOptions{
		LabelSelector: selector.String(),
	}
	if s.tweakListOptions != nil {
		s.tweakListOptions(&listopt)
	}
	list, err := s.client.ExtensionsV1beta1().Ingresses(v1.NamespaceAll).List(listopt)
	if err != nil {
		return nil, err
	}
	for i := range list.Items {
		ret = append(ret, &list.Items[i])
	}
	return ret, nil
}

// Ingresses returns an object that can list and get Ingresses.
func (s *ingressLister) Ingresses(namespace string) v1beta1.IngressNamespaceLister {
	return ingressNamespaceLister{client: s.client, tweakListOptions: s.tweakListOptions, namespace: namespace}
}

// ingressNamespaceLister implements the IngressNamespaceLister
// interface.
type ingressNamespaceLister struct {
	client           kubernetes.Interface
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// List lists all Ingresses in the indexer for a given namespace.
func (s ingressNamespaceLister) List(selector labels.Selector) (ret []*extensionsv1beta1.Ingress, err error) {
	listopt := v1.ListOptions{
		LabelSelector: selector.String(),
	}
	if s.tweakListOptions != nil {
		s.tweakListOptions(&listopt)
	}
	list, err := s.client.ExtensionsV1beta1().Ingresses(s.namespace).List(listopt)
	if err != nil {
		return nil, err
	}
	for i := range list.Items {
		ret = append(ret, &list.Items[i])
	}
	return ret, nil
}

// Get retrieves the Ingress from the indexer for a given namespace and name.
func (s ingressNamespaceLister) Get(name string) (*extensionsv1beta1.Ingress, error) {
	return s.client.ExtensionsV1beta1().Ingresses(s.namespace).Get(name, v1.GetOptions{})
}
