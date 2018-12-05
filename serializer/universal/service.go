package universal

import (
	corev1 "k8s.io/api/core/v1"
)

type Service struct {
	Type  corev1.ServiceType   `json:"type"`
	Name  string               `json:"name"`
	Ports []corev1.ServicePort `json:"ports,omitempty"`
}
