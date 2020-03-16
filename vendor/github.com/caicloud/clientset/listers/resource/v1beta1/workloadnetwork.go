/*
Copyright 2020 caicloud authors. All rights reserved.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "github.com/caicloud/clientset/pkg/apis/resource/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// WorkloadNetworkLister helps list WorkloadNetworks.
type WorkloadNetworkLister interface {
	// List lists all WorkloadNetworks in the indexer.
	List(selector labels.Selector) (ret []*v1beta1.WorkloadNetwork, err error)
	// WorkloadNetworks returns an object that can list and get WorkloadNetworks.
	WorkloadNetworks(namespace string) WorkloadNetworkNamespaceLister
	WorkloadNetworkListerExpansion
}

// workloadNetworkLister implements the WorkloadNetworkLister interface.
type workloadNetworkLister struct {
	indexer cache.Indexer
}

// NewWorkloadNetworkLister returns a new WorkloadNetworkLister.
func NewWorkloadNetworkLister(indexer cache.Indexer) WorkloadNetworkLister {
	return &workloadNetworkLister{indexer: indexer}
}

// List lists all WorkloadNetworks in the indexer.
func (s *workloadNetworkLister) List(selector labels.Selector) (ret []*v1beta1.WorkloadNetwork, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.WorkloadNetwork))
	})
	return ret, err
}

// WorkloadNetworks returns an object that can list and get WorkloadNetworks.
func (s *workloadNetworkLister) WorkloadNetworks(namespace string) WorkloadNetworkNamespaceLister {
	return workloadNetworkNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// WorkloadNetworkNamespaceLister helps list and get WorkloadNetworks.
type WorkloadNetworkNamespaceLister interface {
	// List lists all WorkloadNetworks in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1beta1.WorkloadNetwork, err error)
	// Get retrieves the WorkloadNetwork from the indexer for a given namespace and name.
	Get(name string) (*v1beta1.WorkloadNetwork, error)
	WorkloadNetworkNamespaceListerExpansion
}

// workloadNetworkNamespaceLister implements the WorkloadNetworkNamespaceLister
// interface.
type workloadNetworkNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all WorkloadNetworks in the indexer for a given namespace.
func (s workloadNetworkNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.WorkloadNetwork, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.WorkloadNetwork))
	})
	return ret, err
}

// Get retrieves the WorkloadNetwork from the indexer for a given namespace and name.
func (s workloadNetworkNamespaceLister) Get(name string) (*v1beta1.WorkloadNetwork, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("workloadnetwork"), name)
	}
	return obj.(*v1beta1.WorkloadNetwork), nil
}