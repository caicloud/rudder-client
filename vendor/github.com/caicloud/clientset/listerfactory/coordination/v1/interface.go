/*
Copyright 2019 caicloud authors. All rights reserved.
*/

// Code generated by listerfactory-gen. DO NOT EDIT.

package v1

import (
	internalinterfaces "github.com/caicloud/clientset/listerfactory/internalinterfaces"
	informers "k8s.io/client-go/informers"
	kubernetes "k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/coordination/v1"
)

// Interface provides access to all the listers in this group version.
type Interface interface { // Leases returns a LeaseLister
	Leases() v1.LeaseLister
}

type version struct {
	client           kubernetes.Interface
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

type infromerVersion struct {
	factory informers.SharedInformerFactory
}

// New returns a new Interface.
func New(client kubernetes.Interface, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{client: client, tweakListOptions: tweakListOptions}
}

// NewFrom returns a new Interface.
func NewFrom(factory informers.SharedInformerFactory) Interface {
	return &infromerVersion{factory: factory}
}

// Leases returns a LeaseLister.
func (v *version) Leases() v1.LeaseLister {
	return &leaseLister{client: v.client, tweakListOptions: v.tweakListOptions}
}

// Leases returns a LeaseLister.
func (v *infromerVersion) Leases() v1.LeaseLister {
	return v.factory.Coordination().V1().Leases().Lister()
}
