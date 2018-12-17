package v1

import (
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/caicloud/rudder-client/serializer/universal"
)

type Service struct {
	Type  string               `json:"type"`
	Name  string               `json:"name"`
	Ports []corev1.ServicePort `json:"ports"`
}

type serviceSerializer struct{}

func (s *serviceSerializer) Encode(obj runtime.Object, chart string, cur int) (string, string, error) {
	chconfig, err := universal.PrepareChartConfig(chart, cur)
	if err != nil {
		glog.Error(err)
		return "", "", err
	}
	svc, err := convertObjectToSerivce(obj)
	if err != nil {
		glog.Errorf("convertObjectToSerivce error: %v", err)
		return "", "", err
	}

	glog.Infof("core.v1.Service: %s", svc.String())
	usvc, err := ConvertServiceToController(svc)
	if err != nil {
		glog.Errorf("ConvertServiceToController error: %v", err)
		return "", "", err
	}
	glog.Infof("Service Config: %s", spew.Sdump(usvc))
	if chconfig.Config.Controllers[cur].Services == nil {
		chconfig.Config.Controllers[cur].Services = make([]*universal.Service, 0)
	}
	chconfig.Config.Controllers[cur].Services = append(chconfig.Config.Controllers[cur].Services, usvc)
	glog.Infof("chart chconfig.Config: %s", spew.Sdump(chconfig.Config))
	chconfigBytes, err := json.Marshal(chconfig)
	if err != nil {
		return "", "", err
	}
	return string(chconfigBytes), svc.Name, nil
}

func convertObjectToSerivce(obj runtime.Object) (*corev1.Service, error) {
	svc := new(corev1.Service)
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
	err = json.Unmarshal(data, svc)
	if err != nil {
		glog.Errorf("unstructured object: %s Unmarshal to Deployment error: %v", un.GetName(), err)
		return nil, err
	}
	return svc, nil
}

func ConvertServiceToController(svc *corev1.Service) (*universal.Service, error) {
	return &universal.Service{
		Type:  svc.Spec.Type,
		Name:  svc.Name,
		Ports: svc.Spec.Ports,
	}, nil
}
