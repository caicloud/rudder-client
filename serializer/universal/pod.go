package universal

import (
	corev1 "k8s.io/api/core/v1"
)

type Pod struct {
	Restart         corev1.RestartPolicy      `json:"restart"`
	DNS             corev1.DNSPolicy          `json:"dns"`
	Hostname        string                    `json:"hostname"`
	Subdomain       string                    `json:"subdomain"`
	Termination     *int64                    `json:"termination,omitempty"`
	Host            PodHost                   `json:"host"`
	HostAliases     []corev1.HostAlias        `json:"hostAliases,omitempty"`
	SecurityContext corev1.PodSecurityContext `json:"securityContext"`
	Annotations     []*PodAnnotation          `json:"annotations,omitempty"`
	ConsleIsMonitor *bool                     `json:"__isMonitor,omitempty"`
}

type PodAnnotation struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type PodHost struct {
	Network bool `json:"network"`
	PID     bool `json:"pid"`
	IPC     bool `json:"ipc"`
}

func GetPod(tmpl corev1.PodTemplateSpec) *Pod {
	return &Pod{
		Restart:     tmpl.Spec.RestartPolicy,
		DNS:         tmpl.Spec.DNSPolicy,
		Hostname:    tmpl.Spec.Hostname,
		Subdomain:   tmpl.Spec.Subdomain,
		Termination: tmpl.Spec.TerminationGracePeriodSeconds,
		Host: PodHost{
			Network: tmpl.Spec.HostNetwork,
			PID:     tmpl.Spec.HostPID,
			IPC:     tmpl.Spec.HostIPC,
		},
		HostAliases:     tmpl.Spec.HostAliases,
		Annotations:     getPodAnnotaitions(tmpl),
		ConsleIsMonitor: getConsleIsMonitor(tmpl),
	}
}

func getConsleIsMonitor(tmpl corev1.PodTemplateSpec) *bool {
	if tmpl.Annotations != nil {
		if _, ok := tmpl.Annotations["prometheus.io/scrape"]; ok {
			return convertBoolToPointer(true)
		}
	}
	return convertBoolToPointer(false)
}

func getPodAnnotaitions(tmpl corev1.PodTemplateSpec) []*PodAnnotation {
	ret := make([]*PodAnnotation, 0)
	for k, v := range tmpl.Annotations {
		ret = append(ret, &PodAnnotation{Key: k, Value: v})
	}
	return ret
}
