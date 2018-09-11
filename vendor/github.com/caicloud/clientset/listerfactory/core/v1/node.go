/*
Copyright 2018 caicloud authors. All rights reserved.
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

var _ v1.NodeLister = &nodeLister{}

// nodeLister implements the NodeLister interface.
type nodeLister struct {
	client           kubernetes.Interface
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewNodeLister returns a new NodeLister.
func NewNodeLister(client kubernetes.Interface) v1.NodeLister {
	return NewFilteredNodeLister(client, nil)
}

func NewFilteredNodeLister(client kubernetes.Interface, tweakListOptions internalinterfaces.TweakListOptionsFunc) v1.NodeLister {
	return &nodeLister{
		client:           client,
		tweakListOptions: tweakListOptions,
	}
}

// List lists all Nodes in the indexer.
func (s *nodeLister) List(selector labels.Selector) (ret []*corev1.Node, err error) {
	listopt := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	if s.tweakListOptions != nil {
		s.tweakListOptions(&listopt)
	}
	list, err := s.client.CoreV1().Nodes().List(listopt)
	if err != nil {
		return nil, err
	}
	for i := range list.Items {
		ret = append(ret, &list.Items[i])
	}
	return ret, nil
}

func (s *nodeLister) ListWithPredicate(v1.NodeConditionPredicate) ([]*corev1.Node, error) {
	return nil, nil
}

// Get retrieves the Node from the index for a given name.
func (s *nodeLister) Get(name string) (*corev1.Node, error) {
	return s.client.CoreV1().Nodes().Get(name, metav1.GetOptions{})
}
