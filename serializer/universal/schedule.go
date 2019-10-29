package universal

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Schedule struct {
	SchedulerName string              `json:"scheduler"`
	NodeSelector  map[string]string   `json:"labels,omitempty"`
	Affinity      *Affinity           `json:"affinity,omitempty"`
	AntiAffinity  *AntiAffinity       `json:"antiaffinity,omitempty"`
	Tolerations   []corev1.Toleration `json:"tolerations,omitempty"`
}

type Affinity struct {
	Pod  *PodAffinity  `json:"pod,omitempty"`
	Node *NodeAffinity `json:"node,omitempty"`
}

type PodAffinity struct {
	Type  string        `json:"type"`
	Terms []interface{} `json:"terms,omitempty"`
}

type NodeAffinity struct {
	Type  string        `json:"type"`
	Terms []interface{} `json:"terms,omitempty"`
}

type AntiAffinity struct {
	Pod *PodAffinity `json:"pod,omitempty"`
}

// PreferredSchedulingTerm An empty preferred scheduling term matches all objects with implicit weight 0
// (i.e. it's a no-op). A null preferred scheduling term matches no objects (i.e. is also a no-op).
type PreferredSchedulingTerm struct {
	// Weight associated with matching the corresponding nodeSelectorTerm, in the range 1-100.
	Weight int32 `json:"weight"`
	// A node selector term, associated with the corresponding weight.
	NodeSelectorTerm
}

// NodeSelectorTerm A null or empty node selector term matches no objects. The requirements of
// them are ANDed.
// The TopologySelectorTerm type implements a subset of the NodeSelectorTerm.
type NodeSelectorTerm struct {
	// A list of node selector requirements by node's labels.
	// +optional
	MatchExpressions []NodeSelectorRequirement `json:"expressions,omitempty"`
}

// NodeSelectorRequirement A node selector requirement is a selector that contains values, a key, and an operator
// that relates the key and values.
type NodeSelectorRequirement struct {
	// The label key that the selector applies to.
	Key string `json:"key" protobuf:"bytes,1,opt,name=key"`
	// Represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
	Operator corev1.NodeSelectorOperator `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=NodeSelectorOperator"`
	// An array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. If the operator is Gt or Lt, the values
	// array must have a single element, which will be interpreted as an integer.
	// This array is replaced during a strategic merge patch.
	// +optional
	Values []string `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}

// PodAffinityTerm Defines a set of pods (namely those matching the labelSelector
// relative to the given namespace(s)) that this pod should be
// co-located (affinity) or not co-located (anti-affinity) with,
// where co-located is defined as running on a node whose value of
// the label with key <topologyKey> matches that of any node on which
// a pod of the set of pods is running
type PodAffinityTerm struct {
	// A label query over a set of resources, in this case pods.
	// +optional
	LabelSelector *LabelSelector `json:"selector,omitempty"`
}

// WeightedPodAffinityTerm The weights of all of the matched WeightedPodAffinityTerm fields are added per-node to find the most preferred node(s)
type WeightedPodAffinityTerm struct {
	// weight associated with matching the corresponding podAffinityTerm,
	// in the range 1-100.
	Weight int32 `json:"weight"`
	// Required. A pod affinity term, associated with the corresponding weight.
	PodAffinityTerm
}

// LabelSelector A label selector is a label query over a set of resources. The result of matchLabels and
// matchExpressions are ANDed. An empty label selector matches all objects. A null
// label selector matches no objects.
type LabelSelector struct {
	// matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
	// map is equivalent to an element of matchExpressions, whose key field is "key", the
	// operator is "In", and the values array contains only "value". The requirements are ANDed.
	// +optional
	MatchLabels map[string]string `json:"labels,omitempty"`
	// matchExpressions is a list of label selector requirements. The requirements are ANDed.
	// +optional
	MatchExpressions []LabelSelectorRequirement `json:"expressions,omitempty"`
}

// LabelSelectorRequirement A label selector requirement is a selector that contains values, a key, and an operator that
// relates the key and values.
type LabelSelectorRequirement struct {
	// key is the label key that the selector applies to.
	// +patchMergeKey=key
	// +patchStrategy=merge
	Key string `json:"key" patchStrategy:"merge" patchMergeKey:"key" protobuf:"bytes,1,opt,name=key"`
	// operator represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator metav1.LabelSelectorOperator `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=LabelSelectorOperator"`
	// values is an array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. This array is replaced during a strategic
	// merge patch.
	// +optional
	Values []string `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}

// =================================================================================================

func convertToWeightedPodAffinityTerm(wpat corev1.WeightedPodAffinityTerm) *WeightedPodAffinityTerm {
	return &WeightedPodAffinityTerm{
		Weight:          wpat.Weight,
		PodAffinityTerm: *convertToPodAffinityTerm(wpat.PodAffinityTerm),
	}
}

func convertToPodAffinityTerm(pat corev1.PodAffinityTerm) *PodAffinityTerm {
	if pat.LabelSelector == nil {
		return nil
	}
	ret := &PodAffinityTerm{
		LabelSelector: &LabelSelector{
			MatchLabels:      pat.LabelSelector.MatchLabels,
			MatchExpressions: []LabelSelectorRequirement{},
		},
	}
	for _, me := range pat.LabelSelector.MatchExpressions {
		ret.LabelSelector.MatchExpressions = append(ret.LabelSelector.MatchExpressions, LabelSelectorRequirement{
			Key:      me.Key,
			Operator: me.Operator,
			Values:   me.Values,
		})
	}
	return ret
}

// GetSchedule convert pod affinity info to Schedule info
func GetSchedule(ps corev1.PodSpec) (*Schedule, error) {
	ret := &Schedule{
		SchedulerName: ps.SchedulerName,
		NodeSelector:  ps.NodeSelector,
	}
	if ps.Affinity != nil {
		ret.Affinity = getAffinity(ps.Affinity)
		ret.AntiAffinity = getAntiAffinity(ps.Affinity)
	}
	ret.Tolerations = ps.Tolerations
	return ret, nil
}

// checkPodAffinityPointer make sure point is not empty
func checkPodAffinityPointer(pa **PodAffinity) {
	if *pa == nil {
		*pa = &PodAffinity{}
	}
}

// checkNodeAffinityPointer make sure point is not empty
func checkNodeAffinityPointer(na **NodeAffinity) {
	if *na == nil {
		*na = &NodeAffinity{}
	}
}

func getAffinity(affinity *corev1.Affinity) *Affinity {
	var pa *PodAffinity
	if p := affinity.PodAffinity; p != nil {
		switch {
		case p.RequiredDuringSchedulingIgnoredDuringExecution != nil:
			checkPodAffinityPointer(&pa)
			pa.Type = "Required"
			pa.Terms = make([]interface{}, 0, len(p.RequiredDuringSchedulingIgnoredDuringExecution))
			for _, term := range p.RequiredDuringSchedulingIgnoredDuringExecution {
				pa.Terms = append(pa.Terms, convertToPodAffinityTerm(term))
			}
		case p.PreferredDuringSchedulingIgnoredDuringExecution != nil:
			checkPodAffinityPointer(&pa)
			pa.Type = "Prefered"
			pa.Terms = make([]interface{}, 0, len(p.PreferredDuringSchedulingIgnoredDuringExecution))
			for _, term := range p.PreferredDuringSchedulingIgnoredDuringExecution {
				pa.Terms = append(pa.Terms, convertToWeightedPodAffinityTerm(term))
			}
		}
	}
	var na *NodeAffinity
	if n := affinity.NodeAffinity; n != nil {
		switch {
		case n.RequiredDuringSchedulingIgnoredDuringExecution != nil:
			checkNodeAffinityPointer(&na)
			na.Type = "Required"
			na.Terms = make([]interface{}, 0, len(n.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms))
			for _, term := range n.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms {
				na.Terms = append(na.Terms, convertToNodeSelectorTerm(term))
			}
		case n.PreferredDuringSchedulingIgnoredDuringExecution != nil:
			checkNodeAffinityPointer(&na)
			na.Type = "Prefered"
			na.Terms = make([]interface{}, 0, len(n.PreferredDuringSchedulingIgnoredDuringExecution))
			for _, term := range n.PreferredDuringSchedulingIgnoredDuringExecution {
				na.Terms = append(na.Terms, convertToPreferredSchedulingTerm(term))
			}
		}
	}
	return &Affinity{Pod: pa, Node: na}
}

func convertToNodeSelectorTerm(nst corev1.NodeSelectorTerm) *NodeSelectorTerm {
	ret := &NodeSelectorTerm{
		MatchExpressions: []NodeSelectorRequirement{},
	}
	for _, n := range nst.MatchExpressions {
		ret.MatchExpressions = append(ret.MatchExpressions, NodeSelectorRequirement{
			Key:      n.Key,
			Operator: n.Operator,
			Values:   n.Values,
		})
	}
	return ret
}

func convertToPreferredSchedulingTerm(pst corev1.PreferredSchedulingTerm) *PreferredSchedulingTerm {
	return &PreferredSchedulingTerm{
		Weight:           pst.Weight,
		NodeSelectorTerm: *convertToNodeSelectorTerm(pst.Preference),
	}
}

func getAntiAffinity(affinity *corev1.Affinity) *AntiAffinity {
	var pa *PodAffinity
	if p := affinity.PodAntiAffinity; p != nil {
		switch {
		case p.RequiredDuringSchedulingIgnoredDuringExecution != nil:
			checkPodAffinityPointer(&pa)
			pa.Type = "Required"
			pa.Terms = make([]interface{}, 0, len(p.RequiredDuringSchedulingIgnoredDuringExecution))
			for _, term := range p.RequiredDuringSchedulingIgnoredDuringExecution {
				pa.Terms = append(pa.Terms, convertToPodAffinityTerm(term))
			}
		case p.PreferredDuringSchedulingIgnoredDuringExecution != nil:
			checkPodAffinityPointer(&pa)
			pa.Type = "Prefered"
			pa.Terms = make([]interface{}, 0, len(p.PreferredDuringSchedulingIgnoredDuringExecution))
			for _, term := range p.PreferredDuringSchedulingIgnoredDuringExecution {
				pa.Terms = append(pa.Terms, convertToWeightedPodAffinityTerm(term))
			}
		}
	}
	return &AntiAffinity{Pod: pa}
}
