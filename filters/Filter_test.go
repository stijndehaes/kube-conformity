package filters

import (
	"testing"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/stretchr/testify/assert"
)

func TestNamespaceFilterFilterPodsEmptyString(t *testing.T) {
	filter := Filter{}
	pod1 := newPod("default", "bar1")
	pod2 := newPod("kube-system", "bar2")
	filteredPods := filter.FilterPods([]v1.Pod{pod1, pod2})
	assert.Len(t, filteredPods, 2)
	assert.Equal(t, "bar1", filteredPods[0].Name)
	assert.Equal(t, "bar2", filteredPods[1].Name)
}


func TestNamespaceFilterFilterPodsInclude(t *testing.T) {
	filter := Filter{
		NamespacesString: "default",
	}
	pod1 := newPod("default", "bar1")
	pod2 := newPod("kube-system", "bar2")
	filteredPods := filter.FilterPods([]v1.Pod{pod1, pod2})
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "bar1", filteredPods[0].Name)
}

func TestNamespaceFilterFilterPodsIncludeMultiples(t *testing.T) {
	filter := Filter{
		NamespacesString: "default,kube-system",
	}
	pod1 := newPod("default", "bar1")
	pod2 := newPod("kube-system", "bar2")
	pod3 := newPod("kube-public", "bar3")
	filteredPods := filter.FilterPods([]v1.Pod{pod1, pod2, pod3})
	assert.Len(t, filteredPods, 2)
	assert.Equal(t, "bar1", filteredPods[0].Name)
	assert.Equal(t, "bar2", filteredPods[1].Name)
}

func TestNamespaceFilterFilterPodsExclude(t *testing.T) {
	filter := Filter{
		NamespacesString: "!default",
	}
	pod1 := newPod("default", "bar1")
	pod2 := newPod("kube-system", "bar2")
	filteredPods := filter.FilterPods([]v1.Pod{pod1, pod2})
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "bar2", filteredPods[0].Name)
}

func TestNamespaceFilterFilterPodsExcludeMultiples(t *testing.T) {
	filter := Filter{
		NamespacesString: "!default,!kube-system",
	}
	pod1 := newPod("default", "bar1")
	pod2 := newPod("kube-system", "bar2")
	pod3 := newPod("kube-public", "bar3")
	filteredPods := filter.FilterPods([]v1.Pod{pod1, pod2, pod3})
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "bar3", filteredPods[0].Name)
}

func newPod(namespace, name string) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}
}