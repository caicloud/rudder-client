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

// ExtendedResourceLister helps list ExtendedResources.
type ExtendedResourceLister interface {
	// List lists all ExtendedResources in the indexer.
	List(selector labels.Selector) (ret []*v1beta1.ExtendedResource, err error)
	// Get retrieves the ExtendedResource from the index for a given name.
	Get(name string) (*v1beta1.ExtendedResource, error)
	ExtendedResourceListerExpansion
}

// extendedResourceLister implements the ExtendedResourceLister interface.
type extendedResourceLister struct {
	indexer cache.Indexer
}

// NewExtendedResourceLister returns a new ExtendedResourceLister.
func NewExtendedResourceLister(indexer cache.Indexer) ExtendedResourceLister {
	return &extendedResourceLister{indexer: indexer}
}

// List lists all ExtendedResources in the indexer.
func (s *extendedResourceLister) List(selector labels.Selector) (ret []*v1beta1.ExtendedResource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.ExtendedResource))
	})
	return ret, err
}

// Get retrieves the ExtendedResource from the index for a given name.
func (s *extendedResourceLister) Get(name string) (*v1beta1.ExtendedResource, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("extendedresource"), name)
	}
	return obj.(*v1beta1.ExtendedResource), nil
}