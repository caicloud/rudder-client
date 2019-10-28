package v1

import (
	"encoding/json"
	"fmt"

	"github.com/caicloud/rudder-client/serializer/universal"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type daemonSetSerializer struct{}

func (d *daemonSetSerializer) Encode(obj runtime.Object, chart string, cur int, fn func(runtime.Object) (runtime.Object, error)) (string, string, error) {
	chconfig, err := universal.PrepareChartConfig(chart, cur)
	if err != nil {
		glog.Error(err)
		return "", "", err
	}
	daemon, err := convertObjectToDaemonSet(obj, fn)
	if err != nil {
		glog.Errorf("convertObjectToDaemonSet error: %v", err)
		return "", "", err
	}
	glog.V(4).Infof("apps.v1beta1.DaemonSet: %s", spew.Sdump(daemon))
	controller, err := convertDaemonSetToController(daemon)
	if err != nil {
		glog.Errorf("convertDaemonSetToController error: %v", err)
		return "", "", err
	}
	glog.V(4).Infof("DaemonSet Controller Config: %s", spew.Sdump(controller))
	if chconfig.Config.Controllers[cur] == nil {
		chconfig.Config.Controllers[cur] = new(universal.Controller)
	}
	err = universal.MergeTwoControllers(controller, chconfig.Config.Controllers[cur])
	if err != nil {
		return "", "", err
	}
	chconfigBytes, err := json.Marshal(chconfig)
	if err != nil {
		return "", "", err
	}
	glog.V(4).Infof("chart config: %s", string(chconfigBytes))
	return string(chconfigBytes), daemon.Name, nil
}

func convertObjectToDaemonSet(obj runtime.Object, fn func(runtime.Object) (runtime.Object, error)) (*appsv1.DaemonSet, error) {
	if fn != nil {
		o, err := fn(obj)
		if err != nil {
			return nil, err
		}
		daemonset, ok := o.(*appsv1.DaemonSet)
		if !ok {
			return nil, fmt.Errorf("unknown runtime object type")
		}
		return daemonset, nil
	}
	daemon := new(appsv1.DaemonSet)
	un, ok := obj.(*unstructured.Unstructured)
	if !ok {
		err := fmt.Errorf("assert object as unstructured.Unstructured, object info: %s", spew.Sdump(obj))
		glog.Error(err)
		return nil, err
	}
	ungvk := un.GetObjectKind().GroupVersionKind()
	glog.V(4).Infof("unstructured object gvk: %s", ungvk)
	data, err := un.MarshalJSON()
	if err != nil {
		glog.Errorf("unstructured object: %s MarshalJSON error: %v", un.GetName(), err)
		return nil, err
	}
	err = json.Unmarshal(data, daemon)
	if err != nil {
		glog.Errorf("unstructured object: %s Unmarshal to Deployment error: %v", un.GetName(), err)
		return nil, err
	}
	return daemon, nil
}

func convertDaemonSetToController(daemon *appsv1.DaemonSet) (*universal.Controller, error) {
	tmpl := daemon.Spec.Template

	controller := &universal.DaemonSet{
		Strategy: convertDaemonSetStrategy(daemon.Spec.UpdateStrategy),
		Name:     daemon.Name,
	}

	pod := universal.GetPod(tmpl)
	volumes, err := universal.GetVolumes(tmpl.Spec.Volumes)
	if err != nil {
		glog.Errorf("universal.GetVolumes error: %v", err)
		return nil, err
	}
	initContainers := universal.GetContainers(pod, tmpl.Spec.InitContainers, volumes)
	containers := universal.GetContainers(pod, tmpl.Spec.Containers, volumes)
	schedule, err := universal.GetSchedule(tmpl.Spec)
	if err != nil {
		glog.Errorf("universal.GetSchedule error: %v", err)
		return nil, err
	}

	return &universal.Controller{
		Type:           controller.GetControllerType(),
		Controller:     controller,
		Pod:            pod,
		Schedule:       schedule,
		InitContainers: initContainers,
		Containers:     containers,
		Volumes:        volumes,
	}, nil
}

func convertDaemonSetStrategy(dpStrategy appsv1.DaemonSetUpdateStrategy) universal.Strategy {
	ret := universal.Strategy{
		Type: string(dpStrategy.Type),
	}
	if dpStrategy.Type == "" {
		ret.Type = string(appsv1.RollingUpdateDaemonSetStrategyType)
	}

	if dpStrategy.RollingUpdate != nil {
		ret.Unavailable = dpStrategy.RollingUpdate.MaxUnavailable
	}
	return ret
}
