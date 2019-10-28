package universal

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const annotationLogFilesKey = "logging.caicloud.io/required-logfiles"

type Container struct {
	Name               string                      `json:"name"`
	Image              string                      `json:"image"`
	ImagePullPolicy    corev1.PullPolicy           `json:"imagePullPolicy"`
	TTY                bool                        `json:"tty"`
	Command            []string                    `json:"command"`
	Args               []string                    `json:"args"`
	WorkingDir         string                      `json:"workingDir,omitempty"`
	SecurityContext    *corev1.SecurityContext     `json:"securityContext,omitempty"`
	Ports              []ContainerPort             `json:"ports,omitempty"`
	EnvFrom            []EnvFrom                   `json:"envFrom,omitempty"`
	Env                []Env                       `json:"env,omitempty"`
	Resources          corev1.ResourceRequirements `json:"resources"`
	Mounts             []VolumeMount               `json:"mounts,omitempty"`
	Probe              *ContainerProbe             `json:"probe,omitempty"`
	Lifecycle          *Lifecycle                  `json:"lifecycle,omitempty"`
	ConsoleIsEnvCustom *bool                       `json:"__isEnvCustom,omitempty"`
	ConsoleIsEnvFrom   *bool                       `json:"__isEnvFrom,omitempty"`
	ConsoleIsCommand   *bool                       `json:"__isCommand,omitempty"`
	ConsoleIsMountFile *bool                       `json:"__isMountFile,omitempty"`
	ConsoleIsLog       *bool                       `json:"__isLog,omitempty"`
	ConsoleLiveness    *bool                       `json:"__liveness,omitempty"`
	ConsoleReadiness   *bool                       `json:"__readiness,omitempty"`
}

// Lifecycle describes actions that the management system should take in response to container lifecycle
// events. For the PostStart and PreStop lifecycle handlers, management of the container blocks
// until the action is complete, unless the container process fails, in which case the handler is aborted.
type Lifecycle struct {
	// PostStart is called immediately after a container is created. If the handler fails,
	// the container is terminated and restarted according to its restart policy.
	PostStart *Handler `json:"postStart,omitempty"`
	// PreStop is called immediately before a container is terminated due to an
	// API request or management event such as liveness probe failure,
	// preemption, resource contention, etc.
	PreStop *Handler `json:"preStop,omitempty"`
}

type ContainerPort struct {
	Protocol      corev1.Protocol `json:"protocol"`
	ContainerPort int32           `json:"port"`
}

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EnvFrom struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type VolumeMount struct {
	Name        string `json:"name"`
	ReadOnly    bool   `json:"readonly,omitempty"`
	MountPath   string `json:"path"`
	SubPath     string `json:"subpath,omitempty"`
	ConsoleKind string `json:"__kind,omitempty"`
}

type ContainerProbe struct {
	Liveness  *Probe `json:"liveness,omitempty"`
	Readiness *Probe `json:"readiness,omitempty"`
}

type Probe struct {
	Handler             Handler    `json:"handler,inline"`
	InitialDelaySeconds int32      `json:"delay,omitempty"`
	TimeoutSeconds      int32      `json:"timeout,omitempty"`
	PeriodSeconds       int32      `json:"period,omitempty"`
	Threshold           *Threshold `json:"threshold,omitempty"`
}

type Threshold struct {
	SuccessThreshold int32 `json:"success,omitempty"`
	FailureThreshold int32 `json:"failure,omitempty"`
}

type Handler struct {
	Type   string      `json:"type"`
	Method interface{} `json:"method"`
}

type HTTPGetAction struct {
	Path        string              `json:"path,omitempty"`
	Port        intstr.IntOrString  `json:"port"`
	Host        string              `json:"host,omitempty"`
	Scheme      corev1.URIScheme    `json:"scheme,omitempty"`
	HTTPHeaders []corev1.HTTPHeader `json:"headers,omitempty"`
}

func GetContainers(pod *Pod, containers []corev1.Container, volumes []*Volume) []*Container {
	ret := make([]*Container, 0, len(containers))
	for _, c := range containers {
		vmounts := convertVolumeMounts(c.VolumeMounts, volumes)
		con := &Container{
			Name:               c.Name,
			Image:              c.Image,
			ImagePullPolicy:    c.ImagePullPolicy,
			TTY:                c.TTY,
			Command:            c.Command,
			Args:               c.Args,
			WorkingDir:         c.WorkingDir,
			SecurityContext:    c.SecurityContext,
			EnvFrom:            convertEnvFrom(c.EnvFrom),
			Env:                convertEnv(c.Env),
			Ports:              convertPort(c.Ports),
			Resources:          c.Resources,
			Mounts:             vmounts,
			Probe:              convertContainerProbe(c.LivenessProbe, c.ReadinessProbe),
			Lifecycle:          convertLifecycle(c.Lifecycle),
			ConsoleIsEnvCustom: getConsoleIsEnvCustom(&c),
			ConsoleIsEnvFrom:   getConsoleIsEnvFrom(&c),
			ConsoleIsCommand:   getConsoleIsCommand(&c),
			ConsoleIsMountFile: getConsoleIsMountFile(vmounts),
			ConsoleIsLog:       getConsoleIsLog(pod, c.Name),
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
	for c := range volumes {
		vmap[volumes[c].Name] = volumes[c]
	}
	ret := make([]VolumeMount, 0, len(vmounts))
	for _, vmount := range vmounts {
		if v, ok := vmap[vmount.Name]; ok {
			ret = append(ret, VolumeMount{
				Name:        vmount.Name,
				ReadOnly:    vmount.ReadOnly,
				MountPath:   vmount.MountPath,
				SubPath:     vmount.SubPath,
				ConsoleKind: v.ConsoleKind})
		}
	}
	return ret
}

// =================================================================================================

func convertPort(in []corev1.ContainerPort) []ContainerPort {
	out := make([]ContainerPort, len(in))
	for i, v := range in {
		out[i].ContainerPort = v.ContainerPort
		out[i].Protocol = v.Protocol
	}
	return out
}

func convertEnv(env []corev1.EnvVar) []Env {
	out := make([]Env, 0)
	for _, v := range env {
		if v.ValueFrom != nil {
			continue
		}
		out = append(out, Env{
			Name:  v.Name,
			Value: v.Value,
		})
	}
	return out
}

func convertEnvFrom(envFrom []corev1.EnvFromSource) []EnvFrom {
	if len(envFrom) == 0 {
		return nil
	}
	ret := make([]EnvFrom, 0)
	for _, v := range envFrom {
		switch {
		case v.ConfigMapRef != nil:
			ret = append(ret, EnvFrom{Type: "Config", Name: v.ConfigMapRef.Name})
		case v.SecretRef != nil:
			ret = append(ret, EnvFrom{Type: "Secret", Name: v.SecretRef.Name})
		}
	}
	return ret
}

// =================================================================================================

func convertContainerProbe(liveness, readiness *corev1.Probe) *ContainerProbe {
	ret := new(ContainerProbe)
	if liveness != nil {
		ret.Liveness = convertProbe(liveness)
	}
	if readiness != nil {
		ret.Readiness = convertProbe(readiness)
	}
	return ret
}

func convertProbe(probe *corev1.Probe) *Probe {
	return &Probe{
		Handler:             convertHandler(probe.Handler),
		InitialDelaySeconds: probe.InitialDelaySeconds,
		TimeoutSeconds:      probe.TimeoutSeconds,
		PeriodSeconds:       probe.PeriodSeconds,
		Threshold: &Threshold{
			SuccessThreshold: probe.SuccessThreshold,
			FailureThreshold: probe.FailureThreshold,
		},
	}
}

func convertLifecycle(l *corev1.Lifecycle) *Lifecycle {
	if l == nil {
		return nil
	}
	ret := &Lifecycle{}
	if l.PreStop != nil {
		ps := convertHandler(*l.PreStop)
		ret.PreStop = &ps
	}
	if l.PostStart != nil {
		ps := convertHandler(*l.PostStart)
		ret.PostStart = &ps
	}
	return ret
}

func convertHandler(handler corev1.Handler) Handler {
	ret := Handler{}
	switch {
	case handler.Exec != nil:
		ret.Type = "EXEC"
		ret.Method = handler.Exec
	case handler.HTTPGet != nil:
		ret.Type = "HTTP"
		ret.Method = &HTTPGetAction{
			Path:        handler.HTTPGet.Path,
			Port:        handler.HTTPGet.Port,
			Host:        handler.HTTPGet.Host,
			Scheme:      handler.HTTPGet.Scheme,
			HTTPHeaders: handler.HTTPGet.HTTPHeaders,
		}
	case handler.TCPSocket != nil:
		ret.Type = "TCP"
		ret.Method = handler.TCPSocket
	default:
		glog.Errorf("unsuport handler: %s", handler)
	}

	return ret
}

// =================================================================================================

func getConsoleIsEnvCustom(c *corev1.Container) *bool {
	return convertBoolToPointer(c.Env != nil && len(c.Env) != 0)
}

func getConsoleIsEnvFrom(c *corev1.Container) *bool {
	return convertBoolToPointer(c.EnvFrom != nil && len(c.EnvFrom) != 0)
}

func getConsoleIsCommand(c *corev1.Container) *bool {
	return convertBoolToPointer(c.Command != nil && len(c.Command) != 0)
}

// getConsoleIsLog gets container console flag from pod annotation
func getConsoleIsLog(p *Pod, containerName string) *bool {
	for _, anno := range p.Annotations {
		if anno.Key == annotationLogFilesKey {
			// the log file annotation is like this:
			// logging.caicloud.io/required-logfiles: '{"files":[{"filename":"sda","logDir":"/he","container":"c0"}]}'
			// we can judge the flag by the key format such as below
			value, ok := anno.Value.(string)
			if ok {
				key := fmt.Sprintf(`"container":"%s"`, containerName)
				index := strings.Index(value, key)
				if index != -1 {
					return convertBoolToPointer(true)
				}
			}
			break
		}
	}
	return convertBoolToPointer(false)
}

func getConsoleIsMountFile(vmounts []VolumeMount) *bool {
	for _, vm := range vmounts {
		if vm.ConsoleKind != "" {
			return convertBoolToPointer(true)
		}
	}
	return convertBoolToPointer(false)
}

func getConsoleLiveness(c *corev1.Container) *bool {
	if c == nil {
		return convertBoolToPointer(false)
	}
	if c.LivenessProbe == nil {
		return convertBoolToPointer(false)
	}
	return convertBoolToPointer(true)
}

func getConsoleReadiness(c *corev1.Container) *bool {
	if c == nil {
		return convertBoolToPointer(false)
	}
	if c.ReadinessProbe == nil {
		return convertBoolToPointer(false)
	}
	return convertBoolToPointer(true)
}
