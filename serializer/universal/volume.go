package universal

import (
	"fmt"

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
)

func GetVolumes(volumes []corev1.Volume) ([]*Volume, error) {
	ret := make([]*Volume, 0, len(volumes))
	for _, vol := range volumes {
		vs, err := getVolumeSource(vol.VolumeSource)
		if err != nil {
			glog.Errorf("get %s 's volume source error: %s", vol.Name, err)
			return nil, err
		}
		ret = append(ret, &Volume{
			Name:   vol.Name,
			Source: vs,
		})
	}
	return ret, nil
}

// =================================================================================================

func getVolumeSource(vs corev1.VolumeSource) (VolumeSource, error) {
	switch {
	case vs.HostPath != nil:
		if vs.HostPath.Type != nil {
			return nil, fmt.Errorf("%s", "not support specify hostPath type")
		}
		return &HostPathVolume{Path: vs.HostPath.Path}, nil
	case vs.EmptyDir != nil:
		if vs.EmptyDir.SizeLimit != nil {
			return nil, fmt.Errorf("%s", "not support specify emptyDir SizeLimit")
		}
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

// =================================================================================================

type Volume struct {
	Name   string       `json:"name"`
	Source VolumeSource `json:"source"`
}

type VolumeSource interface {
	GetType() string
}

type DedicatedVolume struct {
	Class  string   `json:"class"`
	Models []string `json:"models"`
}

func (v *DedicatedVolume) GetType() string {
	return "Dedicated"
}

type StaticVolume struct {
	Target   string `json:"target"`
	ReadOnly bool   `json:"readonly"`
}

func (v *StaticVolume) GetType() string {
	return "Static"
}

type ScratchVolume struct {
	Medium corev1.StorageMedium `json:"medium"`
}

func (v *ScratchVolume) GetType() string {
	return "Scratch"
}

type ConfigMapVolume struct {
	Target      string             `json:"target"`
	Items       []corev1.KeyToPath `json:"items"`
	DefaultMode *int32             `json:"default"`
	Optional    *bool              `json:"optional"`
}

func (v *ConfigMapVolume) GetType() string {
	return "ConfigMap"
}

type SecretVolume ConfigMapVolume

func (v *SecretVolume) GetType() string {
	return "Secret"
}

type HostPathVolume struct {
	Path string `json:"path"`
}

func (v *HostPathVolume) GetType() string {
	return "HostPath"
}

type GlusterfsVolume struct {
	Endpoints string `json:"endpoints"`
	Path      string `json:"path"`
	Readonly  bool   `json:"readonly"`
}

func (v *GlusterfsVolume) GetType() string {
	return "Glusterfs"
}

// type NFSVolume struct{}

// type CephVolume struct{}
