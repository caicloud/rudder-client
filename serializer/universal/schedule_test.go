package universal

import (
	"encoding/json"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetSchedule(t *testing.T) {
	marshalStr := `{"scheduler":"default","affinity":{"pod":{"type":"Prefered","terms":[{"weight":10,"selector":{"labels":{"hello":"world"},"expressions":[{"key":"key","operator":"NotIn","values":["hello","world"]}]}}]},"node":{"type":"Required","terms":[{"expressions":[{"key":"key","operator":"NotIn","values":["hello","world"]}]}]}},"antiaffinity":{"pod":{"type":"Prefered","terms":[{"weight":10,"selector":{"labels":{"hello":"world"},"expressions":[{"key":"key","operator":"NotIn","values":["hello","world"]}]}}]}}}`
	nst := corev1.NodeSelectorTerm{
		MatchExpressions: []corev1.NodeSelectorRequirement{
			{
				Key:      "key",
				Operator: corev1.NodeSelectorOpNotIn,
				Values:   []string{"hello", "world"},
			},
		},
	}
	pat := corev1.PodAffinityTerm{
		LabelSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"hello": "world"},
			MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "key",
					Operator: metav1.LabelSelectorOpNotIn,
					Values:   []string{"hello", "world"},
				},
			},
		},
	}
	ps := corev1.PodSpec{
		SchedulerName: "default",
		Affinity: &corev1.Affinity{
			NodeAffinity: &corev1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
					NodeSelectorTerms: []corev1.NodeSelectorTerm{nst},
				},
				// PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
				// 	{
				// 		Weight:     10,
				// 		Preference: nst,
				// 	},
				// },
			},
			PodAffinity: &corev1.PodAffinity{
				//RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{pat},
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
					{
						Weight:          10,
						PodAffinityTerm: pat,
					},
				},
			},
			PodAntiAffinity: &corev1.PodAntiAffinity{
				//RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{pat},
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
					{
						Weight:          10,
						PodAffinityTerm: pat,
					},
				},
			},
		},
	}
	schedule, err := GetSchedule(ps)
	if err != nil {
		t.Errorf("GetSchedule failed %v", err)
	}
	bytes, err := json.Marshal(schedule)
	if err != nil {
		t.Errorf("GetSchedule failed %v", err)
	}
	if string(bytes) != marshalStr {
		t.Errorf("convert schedule failed  %s", string(bytes))
	}
}
