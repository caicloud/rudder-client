package universal

import (
	corev1 "k8s.io/api/core/v1"
)

type Container struct {
	Name            string                      `json:"name"`
	Image           string                      `json:"image"`
	ImagePullPolicy corev1.PullPolicy           `json:"imagePullPolicy"`
	TTY             bool                        `json:"tty"`
	Command         []string                    `json:"command"`
	Args            []string                    `json:"args"`
	WorkingDir      string                      `json:"workingDir"`
	SecurityContext *corev1.SecurityContext     `json:"securityContext,omitempty"`
	Ports           []corev1.ContainerPort      `json:"ports,omitempty"`
	Env             []corev1.EnvVar             `json:"env,omitempty"`
	EnvFrom         []corev1.EnvFromSource      `json:"envFrom,omitempty"`
	Resources       corev1.ResourceRequirements `json:"resources"`
	Mounts          []corev1.VolumeMount        `json:"mounts,omitempty"`
	Probe           *ContainerProbe             `json:"probe,omitempty"`
	Lifecycle       *corev1.Lifecycle           `json:"lifecycle,omitempty"`
}

type ContainerProbe struct {
	Liveness  *corev1.Probe `json:"liveness,omitempty"`
	Readiness *corev1.Probe `json:"readiness,omitempty"`
}

func GetContainers(containers []corev1.Container) []*Container {
	ret := make([]*Container, 0, len(containers))
	for _, c := range containers {
		ret = append(ret, &Container{
			Name:            c.Name,
			Image:           c.Image,
			ImagePullPolicy: c.ImagePullPolicy,
			TTY:             c.TTY,
			Command:         c.Command,
			Args:            c.Args,
			WorkingDir:      c.WorkingDir,
			SecurityContext: c.SecurityContext,
			Ports:           c.Ports,
			EnvFrom:         c.EnvFrom,
			Env:             c.Env,
			Resources:       c.Resources,
			Mounts:          c.VolumeMounts,
			Probe: &ContainerProbe{
				Liveness:  c.LivenessProbe,
				Readiness: c.ReadinessProbe,
			},
			Lifecycle: c.Lifecycle,
		})
	}
	return ret
}
