/*
Copyright 2020 caicloud authors. All rights reserved.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/caicloud/clientset/pkg/apis/devops/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CargoLister helps list Cargos.
type CargoLister interface {
	// List lists all Cargos in the indexer.
	List(selector labels.Selector) (ret []*v1.Cargo, err error)
	// Get retrieves the Cargo from the index for a given name.
	Get(name string) (*v1.Cargo, error)
	CargoListerExpansion
}

// cargoLister implements the CargoLister interface.
type cargoLister struct {
	indexer cache.Indexer
}

// NewCargoLister returns a new CargoLister.
func NewCargoLister(indexer cache.Indexer) CargoLister {
	return &cargoLister{indexer: indexer}
}

// List lists all Cargos in the indexer.
func (s *cargoLister) List(selector labels.Selector) (ret []*v1.Cargo, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Cargo))
	})
	return ret, err
}

// Get retrieves the Cargo from the index for a given name.
func (s *cargoLister) Get(name string) (*v1.Cargo, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("cargo"), name)
	}
	return obj.(*v1.Cargo), nil
}