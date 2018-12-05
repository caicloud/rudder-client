package universal

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/intstr"
)

type Controller struct {
	Type           string           `json:"type"`
	Controller     ControllerConfig `json:"controller"`
	Pod            *Pod             `json:"pod"`
	Schedule       *Schedule        `json:"schedule,omitempty"`
	InitContainers []*Container     `json:"initContainers,omitempty"`
	Containers     []*Container     `json:"containers"`
	Volumes        []*Volume        `json:"volumes,omitempty"`
	Services       []*Service       `json:"services,omitempty"`
}

type ControllerConfig interface {
	GetControllerType() string
}

// =================================================================================================

type Deployment struct {
	Replica  *int32   `json:"replica,omitempty"`
	Strategy Strategy `json:"strategy,omitempty"`
	Ready    int32    `json:"ready"`
}

type Strategy struct {
	Type        string              `json:"type"`
	Unavailable *intstr.IntOrString `json:"unavailable,omitempty"`
	Surge       *intstr.IntOrString `json:"surge,omitempty"`
}

func (d *Deployment) GetControllerType() string {
	return "Deployment"
}

// =================================================================================================

func MergeTwoControllers(src, dst *Controller) error {
	if dst == nil || src == nil {
		return fmt.Errorf("both of input controlllers must not nil")
	}
	dst.Controller = src.Controller
	dst.Pod = src.Pod
	dst.Schedule = src.Schedule
	dst.InitContainers = src.InitContainers
	dst.Containers = src.Containers
	return nil
}
