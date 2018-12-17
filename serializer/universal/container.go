package universal

import (
	corev1 "k8s.io/api/core/v1"
)

type Container struct {
	Name               string                      `json:"name"`
	Image              string                      `json:"image"`
	ImagePullPolicy    corev1.PullPolicy           `json:"imagePullPolicy"`
	TTY                bool                        `json:"tty"`
	Command            []string                    `json:"command"`
	Args               []string                    `json:"args"`
	WorkingDir         string                      `json:"workingDir,omitempty"`
	SecurityContext    *corev1.SecurityContext     `json:"securityContext,omitempty"`
	Ports              []corev1.ContainerPort      `json:"ports,omitempty"`
	Env                []corev1.EnvVar             `json:"env,omitempty"`
	EnvFrom            []corev1.EnvFromSource      `json:"envFrom,omitempty"`
	Resources          corev1.ResourceRequirements `json:"resources"`
	Mounts             []VolumeMount               `json:"mounts,omitempty"`
	Probe              *ContainerProbe             `json:"probe,omitempty"`
	Lifecycle          *corev1.Lifecycle           `json:"lifecycle,omitempty"`
	ConsoleIsEnvCustom *bool                       `json:"__isEnvCustom,omitempty"`
	ConsoleIsEnvFrom   *bool                       `json:"__isEnvFrom,omitempty"`
	ConsoleIsCommand   *bool                       `json:"__isCommand,omitempty"`
	ConsoleIsMountFile *bool                       `json:"__isMountFile,omitempty"`
	ConsoleIsLog       *bool                       `json:"__isLog,omitempty"`
	ConsoleLiveness    *bool                       `json:"__liveness,omitempty"`
	ConsoleReadiness   *bool                       `json:"__readiness,omitempty"`
}

type VolumeMount struct {
	corev1.VolumeMount `json:",inline"`
	ConsoleKind        string `json:"__kind,omitempty"`
}

type ContainerProbe struct {
	Liveness  *corev1.Probe `json:"liveness,omitempty"`
	Readiness *corev1.Probe `json:"readiness,omitempty"`
}

func GetContainers(containers []corev1.Container, volumes []*Volume) []*Container {
	ret := make([]*Container, 0, len(containers))
	for _, c := range containers {
		vmounts := convertVolumeMounts(c.VolumeMounts, volumes)
		con := &Container{
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
			Mounts:          vmounts,
			Probe: &ContainerProbe{
				Liveness:  c.LivenessProbe,
				Readiness: c.ReadinessProbe,
			},
			Lifecycle:          c.Lifecycle,
			ConsoleIsEnvCustom: getConsoleIsEnvCustom(&c),
			ConsoleIsEnvFrom:   getConsoleIsEnvFrom(&c),
			ConsoleIsCommand:   getConsoleIsCommand(&c),
			ConsoleIsMountFile: getConsoleIsMountFile(vmounts),
			ConsoleIsLog:       getConsoleIsLog(&c),
			ConsoleLiveness:    getConsoleLiveness(&c),
			ConsoleReadiness:   getConsoleReadiness(&c),
		}
		ret = append(ret, con)
	}
	return ret
}

// =================================================================================================

func convertVolumeMounts(vmounts []corev1.VolumeMount, volumes []*Volume) []VolumeMount {
	vmap := make(map[string]*Volume)
	for c, _ := range volumes {
		vmap[volumes[c].Name] = volumes[c]
	}
	ret := make([]VolumeMount, 0, len(vmounts))
	for _, vmount := range vmounts {
		if v, ok := vmap[vmount.Name]; ok {
			ret = append(ret, VolumeMount{VolumeMount: vmount, ConsoleKind: v.ConsoleKind})
		}
	}
	return ret
}

// =================================================================================================

func getConsoleLiveness(c *corev1.Container) *bool {
	if c == nil {
		return nil
	}
	if c.LivenessProbe == nil {
		return nil
	}
	ret := true
	return &ret
}

func getConsoleReadiness(c *corev1.Container) *bool {
	if c == nil {
		return nil
	}
	if c.ReadinessProbe == nil {
		return nil
	}
	ret := true
	return &ret
}

func getConsoleIsMountFile(vmounts []VolumeMount) *bool {
	ret := false
	for _, vm := range vmounts {
		if vm.ConsoleKind != "" {
			ret = true
			return &ret
		}
	}
	return &ret
}

func getConsoleIsLog(c *corev1.Container) *bool {
	ret := true
	return &ret
}

func getConsoleIsEnvCustom(c *corev1.Container) *bool {
	ret := true
	return &ret
}

func getConsoleIsEnvFrom(c *corev1.Container) *bool {
	ret := true
	return &ret
}

func getConsoleIsCommand(c *corev1.Container) *bool {
	ret := true
	return &ret
}
