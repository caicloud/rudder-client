package universal

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Schedule struct {
	SchedulerName string              `json:"scheduler"`
	NodeSelector  map[string]string   `json:"labels"`
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

type LabelSelector struct {
	MatchLabels      map[string]string          `json:"labels,omitempty"`
	MatchExpressions []LabelSelectorRequirement `json:"expressions,omitempty"`
}

type LabelSelectorRequirement struct {
	Key      string                       `json:"key"`
	Operator metav1.LabelSelectorOperator `json:"operator"`
	Values   []string                     `json:"values,omitempty"`
}

type NodeAffinity struct {
	Type  string        `json:"type"`
	Terms []interface{} `json:"terms,omitempty"`
}

type NodeAffinityTerm struct {
	Weight        int32          `json:"weight"`
	LabelSelector *LabelSelector `json:"selector,omitempty"`
}

type AntiAffinity struct {
	Pod *PodAffinity `json:"pod,omitempty"`
}

// =================================================================================================

func GetSchedule(pspec corev1.PodSpec) (*Schedule, error) {
	ret := &Schedule{
		SchedulerName: pspec.SchedulerName,
		NodeSelector:  pspec.NodeSelector,
	}
	if pspec.Affinity != nil {
		ret.Affinity = getAffinity(pspec.Affinity)
		ret.AntiAffinity = getAntiAffinity(pspec.Affinity)
	}
	ret.Tolerations = pspec.Tolerations
	return ret, nil
}

func getAffinity(affinity *corev1.Affinity) *Affinity {
	pod := new(PodAffinity)
	if p := affinity.PodAffinity; p != nil {
		switch {
		case p.RequiredDuringSchedulingIgnoredDuringExecution != nil:
			pod.Type = "Required"
			pod.Terms = make([]interface{}, 0, len(p.RequiredDuringSchedulingIgnoredDuringExecution))
			for _, term := range p.RequiredDuringSchedulingIgnoredDuringExecution {
				pod.Terms = append(pod.Terms, term)
			}
		case p.PreferredDuringSchedulingIgnoredDuringExecution != nil:
			pod.Type = "Prefered"
			pod.Terms = make([]interface{}, 0, len(p.PreferredDuringSchedulingIgnoredDuringExecution))
			for _, term := range p.PreferredDuringSchedulingIgnoredDuringExecution {
				pod.Terms = append(pod.Terms, term)
			}
		}
	}
	node := new(NodeAffinity)
	if n := affinity.NodeAffinity; n != nil {
		switch {
		case n.RequiredDuringSchedulingIgnoredDuringExecution != nil:
			node.Type = "Required"
			node.Terms = make([]interface{}, 1)
			node.Terms[0] = n.RequiredDuringSchedulingIgnoredDuringExecution
		case n.PreferredDuringSchedulingIgnoredDuringExecution != nil:
			node.Type = "Prefered"
			node.Terms = make([]interface{}, 0, len(n.PreferredDuringSchedulingIgnoredDuringExecution))
			for _, term := range n.PreferredDuringSchedulingIgnoredDuringExecution {
				node.Terms = append(node.Terms, term)
			}
		}
	}
	return &Affinity{Pod: pod, Node: node}
}

func getAntiAffinity(affinity *corev1.Affinity) *AntiAffinity {
	pod := new(PodAffinity)
	if p := affinity.PodAntiAffinity; p != nil {
		switch {
		case p.RequiredDuringSchedulingIgnoredDuringExecution != nil:
			pod.Type = "Required"
			pod.Terms = make([]interface{}, 0, len(p.RequiredDuringSchedulingIgnoredDuringExecution))
			for _, term := range p.RequiredDuringSchedulingIgnoredDuringExecution {
				pod.Terms = append(pod.Terms, term)
			}
		case p.PreferredDuringSchedulingIgnoredDuringExecution != nil:
			pod.Type = "Prefered"
			pod.Terms = make([]interface{}, 0, len(p.PreferredDuringSchedulingIgnoredDuringExecution))
			for _, term := range p.PreferredDuringSchedulingIgnoredDuringExecution {
				pod.Terms = append(pod.Terms, term)
			}
		}
	}
	return &AntiAffinity{Pod: pod}
}
