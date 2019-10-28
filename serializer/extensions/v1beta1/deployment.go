package v1beta1

import (
	"encoding/json"
	"fmt"

	"github.com/caicloud/rudder-client/serializer/universal"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type deploymentSerializer struct{}

func (d *deploymentSerializer) Encode(obj runtime.Object, chart string, cur int, fn func(runtime.Object) (runtime.Object, error)) (string, string, error) {
	chconfig, err := universal.PrepareChartConfig(chart, cur)
	if err != nil {
		glog.Error(err)
		return "", "", err
	}
	dp, err := convertObjectToDeploy(obj)
	if err != nil {
		glog.Errorf("convertObjectToDeploy error: %v", err)
		return "", "", err
	}
	glog.Infof("extensions.v1beta1..Deployment: %s", spew.Sdump(dp))
	controller, err := convertDeployToController(dp)
	if err != nil {
		glog.Errorf("convertDeployToController error: %v", err)
		return "", "", err
	}
	glog.Infof("Deployment Controller Config: %s", spew.Sdump(controller))
	if chconfig.Config.Controllers[cur] == nil {
		chconfig.Config.Controllers[cur] = new(universal.Controller)
	}
	err = universal.MergeTwoControllers(controller, chconfig.Config.Controllers[cur])
	if err != nil {
		return "", "", err
	}
	glog.Infof("chart config: %s", spew.Sdump(chconfig.Config))
	chconfigBytes, err := json.Marshal(chconfig)
	if err != nil {
		return "", "", err
	}
	return string(chconfigBytes), dp.Name, nil
}

func convertObjectToDeploy(obj runtime.Object) (*extensionsv1beta1.Deployment, error) {
	dp := new(extensionsv1beta1.Deployment)
	un, ok := obj.(*unstructured.Unstructured)
	if !ok {
		err := fmt.Errorf("assert object as unstructured.Unstructured, object info: %s", spew.Sdump(obj))
		glog.Error(err)
		return nil, err
	}
	ungvk := un.GetObjectKind().GroupVersionKind()
	glog.Infof("unstructured object gvk: %s", ungvk)
	data, err := un.MarshalJSON()
	if err != nil {
		glog.Errorf("unstructured object: %s MarshalJSON error: %v", un.GetName(), err)
		return nil, err
	}
	err = json.Unmarshal(data, dp)
	if err != nil {
		glog.Errorf("unstructured object: %s Unmarshal to Deployment error: %v", un.GetName(), err)
		return nil, err
	}
	return dp, nil
}

func convertDeployToController(dp *extensionsv1beta1.Deployment) (*universal.Controller, error) {
	spec := dp.Spec
	tmpl := dp.Spec.Template

	controller := &universal.Deployment{
		Replica:  spec.Replicas,
		Strategy: convertDeployStrategy(dp.Spec.Strategy),
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

func convertDeployStrategy(dpStrategy extensionsv1beta1.DeploymentStrategy) universal.Strategy {
	ret := universal.Strategy{
		Type: string(dpStrategy.Type),
	}
	if dpStrategy.Type == "" {
		ret.Type = string(extensionsv1beta1.RollingUpdateDeploymentStrategyType)
	}

	if dpStrategy.RollingUpdate != nil {
		ret.Unavailable = dpStrategy.RollingUpdate.MaxUnavailable
		ret.Surge = dpStrategy.RollingUpdate.MaxSurge
	}

	return ret
}
