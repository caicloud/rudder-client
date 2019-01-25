/*
Copyright 2019 caicloud authors. All rights reserved.
*/

// Code generated by listerfactory-gen. DO NOT EDIT.

package apps

import (
	v1 "github.com/caicloud/clientset/listerfactory/apps/v1"
	v1beta1 "github.com/caicloud/clientset/listerfactory/apps/v1beta1"
	v1beta2 "github.com/caicloud/clientset/listerfactory/apps/v1beta2"
	internalinterfaces "github.com/caicloud/clientset/listerfactory/internalinterfaces"
	informers "k8s.io/client-go/informers"
	kubernetes "k8s.io/client-go/kubernetes"
)

// Interface provides access to each of this group's versions.
type Interface interface {
	// V1 provides access to listers for resources in V1.
	V1() v1.Interface
	// V1beta2 provides access to listers for resources in V1beta2.
	V1beta2() v1beta2.Interface
	// V1beta1 provides access to listers for resources in V1beta1.
	V1beta1() v1beta1.Interface
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

// V1beta2 returns a new v1beta2.Interface.
func (g *group) V1beta2() v1beta2.Interface {
	return v1beta2.New(g.client, g.tweakListOptions)
}

func (g *informerGroup) V1beta2() v1beta2.Interface {
	return v1beta2.NewFrom(g.factory)
}

// V1beta1 returns a new v1beta1.Interface.
func (g *group) V1beta1() v1beta1.Interface {
	return v1beta1.New(g.client, g.tweakListOptions)
}

func (g *informerGroup) V1beta1() v1beta1.Interface {
	return v1beta1.NewFrom(g.factory)
}
