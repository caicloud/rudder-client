/*
Copyright 2018 caicloud authors. All rights reserved.
*/

// Code generated by listerfactory-gen. DO NOT EDIT.

package v1

import (
	internalinterfaces "github.com/caicloud/clientset/listerfactory/internalinterfaces"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubernetes "k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/apps/v1"
)

var _ v1.DeploymentLister = &deploymentLister{}

var _ v1.DeploymentNamespaceLister = &deploymentNamespaceLister{}

// deploymentLister implements the DeploymentLister interface.
type deploymentLister struct {
	client           kubernetes.Interface
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewDeploymentLister returns a new DeploymentLister.
func NewDeploymentLister(client kubernetes.Interface) v1.DeploymentLister {
	return NewFilteredDeploymentLister(client, nil)
}

func NewFilteredDeploymentLister(client kubernetes.Interface, tweakListOptions internalinterfaces.TweakListOptionsFunc) v1.DeploymentLister {
	return &deploymentLister{
		client:           client,
		tweakListOptions: tweakListOptions,
	}
}

// List lists all Deployments in the indexer.
func (s *deploymentLister) List(selector labels.Selector) (ret []*appsv1.Deployment, err error) {
	listopt := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	if s.tweakListOptions != nil {
		s.tweakListOptions(&listopt)
	}
	list, err := s.client.AppsV1().Deployments(metav1.NamespaceAll).List(listopt)
	if err != nil {
		return nil, err
	}
	for i := range list.Items {
		ret = append(ret, &list.Items[i])
	}
	return ret, nil
}

func (s *deploymentLister) GetDeploymentsForReplicaSet(*appsv1.ReplicaSet) ([]*appsv1.Deployment, error) {
	return nil, nil
}

// Deployments returns an object that can list and get Deployments.
func (s *deploymentLister) Deployments(namespace string) v1.DeploymentNamespaceLister {
	return deploymentNamespaceLister{client: s.client, tweakListOptions: s.tweakListOptions, namespace: namespace}
}

// deploymentNamespaceLister implements the DeploymentNamespaceLister
// interface.
type deploymentNamespaceLister struct {
	client           kubernetes.Interface
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// List lists all Deployments in the indexer for a given namespace.
func (s deploymentNamespaceLister) List(selector labels.Selector) (ret []*appsv1.Deployment, err error) {
	listopt := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	if s.tweakListOptions != nil {
		s.tweakListOptions(&listopt)
	}
	list, err := s.client.AppsV1().Deployments(s.namespace).List(listopt)
	if err != nil {
		return nil, err
	}
	for i := range list.Items {
		ret = append(ret, &list.Items[i])
	}
	return ret, nil
}

// Get retrieves the Deployment from the indexer for a given namespace and name.
func (s deploymentNamespaceLister) Get(name string) (*appsv1.Deployment, error) {
	return s.client.AppsV1().Deployments(s.namespace).Get(name, metav1.GetOptions{})
}
