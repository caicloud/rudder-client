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
	Annotations     map[string]string         `json:"annotations,omitempty`
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
		HostAliases: tmpl.Spec.HostAliases,
		Annotations: tmpl.Annotations,
	}
}
