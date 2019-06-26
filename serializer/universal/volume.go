package universal

import (
	"fmt"

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
)

type Volume struct {
	ConsoleKind string      `json:"__kind,omitempty"`
	Name        string      `json:"name"`
	VolumeType  string      `json:"type"`
	Source      interface{} `json:"source,omitempty"`
}

// =================================================================================================

type VolumeSource interface {
	GetType() string
	GetConsoleKind() string
}

// =================================================================================================

type DedicatedVolume struct {
	Class  string   `json:"class"`
	Models []string `json:"models"`
}

func (v *DedicatedVolume) GetType() string {
	return "Dedicated"
}

func (v *DedicatedVolume) GetConsoleKind() string {
	return ""
}

// =================================================================================================

type StaticVolume struct {
	Target   string `json:"target"`
	ReadOnly bool   `json:"readonly"`
}

func (v *StaticVolume) GetType() string {
	return "Static"
}

func (v *StaticVolume) GetConsoleKind() string {
	return ""
}

// =================================================================================================

type ScratchVolume struct {
	Medium corev1.StorageMedium `json:"medium"`
}

func (v *ScratchVolume) GetType() string {
	return "Scratch"
}

func (v *ScratchVolume) GetConsoleKind() string {
	return ""
}

// =================================================================================================

type ConfigMapVolume struct {
	Target      string             `json:"target"`
	Items       []corev1.KeyToPath `json:"items"`
	DefaultMode *int32             `json:"default"`
	Optional    *bool              `json:"optional"`
}

func (v *ConfigMapVolume) GetType() string {
	return "Config"
}

func (v *ConfigMapVolume) GetConsoleKind() string {
	return "config"
}

// =================================================================================================

type SecretVolume ConfigMapVolume

func (v *SecretVolume) GetType() string {
	return "Secret"
}

func (v *SecretVolume) GetConsoleKind() string {
	return ""
}

// =================================================================================================

type HostPathVolume struct {
	Path string `json:"path"`
}

func (v *HostPathVolume) GetType() string {
	return "HostPath"
}

func (v *HostPathVolume) GetConsoleKind() string {
	return ""
}

// =================================================================================================

type GlusterfsVolume struct {
	Endpoints string `json:"endpoints"`
	Path      string `json:"path"`
	Readonly  bool   `json:"readonly"`
}

func (v *GlusterfsVolume) GetType() string {
	return "Glusterfs"
}

func (v *GlusterfsVolume) GetConsoleKind() string {
	return ""
}

// =================================================================================================

// type NFSVolume struct{}

// type CephVolume struct{}

// =================================================================================================

func GetVolumes(volumes []corev1.Volume) ([]*Volume, error) {
	ret := make([]*Volume, 0, len(volumes))
	for _, vol := range volumes {
		vs, err := getVolumeSource(vol.VolumeSource)
		if err != nil {
			glog.Errorf("get %s 's volume source error: %s", vol.Name, err)
			return nil, err
		}
		typ := vs.GetType()
		ckind := vs.GetConsoleKind()
		ret = append(ret, &Volume{
			Name:        vol.Name,
			VolumeType:  typ,
			ConsoleKind: ckind,
			Source:      vs,
		})
	}
	return ret, nil
}

// =================================================================================================

func getVolumeSource(vs corev1.VolumeSource) (VolumeSource, error) {
	switch {
	case vs.HostPath != nil:
		//if vs.HostPath.Type != nil {
		//	return nil, fmt.Errorf("%s", "not support specify hostPath type")
		//}
		return &HostPathVolume{Path: vs.HostPath.Path}, nil
	case vs.EmptyDir != nil:
		//if vs.EmptyDir.SizeLimit != nil {
		//	return nil, fmt.Errorf("%s", "not support specify emptyDir SizeLimit")
		//}
		return &ScratchVolume{Medium: vs.EmptyDir.Medium}, nil
	case vs.Secret != nil:
		return &SecretVolume{
			Target:      vs.Secret.SecretName,
			Items:       vs.Secret.Items,
			DefaultMode: vs.Secret.DefaultMode,
			Optional:    vs.Secret.Optional,
		}, nil
	// case vs.NFS != nil:
	case vs.Glusterfs != nil:
		return &GlusterfsVolume{
			Endpoints: vs.Glusterfs.EndpointsName,
			Path:      vs.Glusterfs.Path,
			Readonly:  vs.Glusterfs.ReadOnly,
		}, nil
	case vs.PersistentVolumeClaim != nil:
		return &StaticVolume{
			Target:   vs.PersistentVolumeClaim.ClaimName,
			ReadOnly: vs.PersistentVolumeClaim.ReadOnly,
		}, nil
	// case vs.CephFS != nil:
	case vs.ConfigMap != nil:
		return &ConfigMapVolume{
			Target:      vs.ConfigMap.Name,
			Items:       vs.ConfigMap.Items,
			DefaultMode: vs.ConfigMap.DefaultMode,
			Optional:    vs.ConfigMap.Optional,
		}, nil
	default:
		glog.Infof("volume source: %s", vs.String())
		return nil, fmt.Errorf("not support the volume %s", vs.String())
	}
}
