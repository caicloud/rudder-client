/*
Copyright 2020 caicloud authors. All rights reserved.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/caicloud/clientset/pkg/apis/clever/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// TemplateLister helps list Templates.
type TemplateLister interface {
	// List lists all Templates in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Template, err error)
	// Get retrieves the Template from the index for a given name.
	Get(name string) (*v1alpha1.Template, error)
	TemplateListerExpansion
}

// templateLister implements the TemplateLister interface.
type templateLister struct {
	indexer cache.Indexer
}

// NewTemplateLister returns a new TemplateLister.
func NewTemplateLister(indexer cache.Indexer) TemplateLister {
	return &templateLister{indexer: indexer}
}

// List lists all Templates in the indexer.
func (s *templateLister) List(selector labels.Selector) (ret []*v1alpha1.Template, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Template))
	})
	return ret, err
}

// Get retrieves the Template from the index for a given name.
func (s *templateLister) Get(name string) (*v1alpha1.Template, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("template"), name)
	}
	return obj.(*v1alpha1.Template), nil
}