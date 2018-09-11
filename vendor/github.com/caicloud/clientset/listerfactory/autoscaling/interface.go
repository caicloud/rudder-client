/*
Copyright 2018 caicloud authors. All rights reserved.
*/

// Code generated by listerfactory-gen. DO NOT EDIT.

package autoscaling

import (
	v1 "github.com/caicloud/clientset/listerfactory/autoscaling/v1"
	v2beta1 "github.com/caicloud/clientset/listerfactory/autoscaling/v2beta1"
	internalinterfaces "github.com/caicloud/clientset/listerfactory/internalinterfaces"
	informers "k8s.io/client-go/informers"
	kubernetes "k8s.io/client-go/kubernetes"
)

// Interface provides access to each of this group's versions.
type Interface interface {
	// V1 provides access to listers for resources in V1.
	V1() v1.Interface
	// V2beta1 provides access to listers for resources in V2beta1.
	V2beta1() v2beta1.Interface
}

type group struct {
	client           kubernetes.Interface
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

type informerGroup struct {
	factory informers.SharedInformerFactory
}

// New returns a new Interface.
func New(client kubernetes.Interface, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &group{client: client, tweakListOptions: tweakListOptions}
}

// NewFrom returns a new Interface
func NewFrom(factory informers.SharedInformerFactory) Interface {
	return &informerGroup{factory: factory}
}

// V1 returns a new v1.Interface.
func (g *group) V1() v1.Interface {
	return v1.New(g.client, g.tweakListOptions)
}

func (g *informerGroup) V1() v1.Interface {
	return v1.NewFrom(g.factory)
}

// V2beta1 returns a new v2beta1.Interface.
func (g *group) V2beta1() v2beta1.Interface {
	return v2beta1.New(g.client, g.tweakListOptions)
}

func (g *informerGroup) V2beta1() v2beta1.Interface {
	return v2beta1.NewFrom(g.factory)
}
