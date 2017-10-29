package filters

import (
	"testing"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/stretchr/testify/assert"
)

func newPod(namespace, name string) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}
}

func newPodWithAnnotations(namespace, name string, annotations map[string]string ) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Annotations: annotations,
		},
	}
}

func TestFilter_FilterIncludeNamespace(t *testing.T) {
	filter := Filter{
		IncludeNamespaces: []string{"default"},
	}

	pods := []v1.Pod{
		newPod("default", "name1"),
		newPod("kube-system", "name2"),
	}

	filteredPods := filter.FilterIncludeNamespace(pods)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name1", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
}

func TestFilter_FilterIncludeNamespace_Empty(t *testing.T) {
	filter := Filter{
		IncludeNamespaces: []string{},
	}

	pods := []v1.Pod{newPod("default", "name")}

	filteredPods := filter.FilterIncludeNamespace(pods)
	assert.Len(t, filteredPods, 1)
}

func TestFilter_FilterExcludeNamespace(t *testing.T) {
	filter := Filter{
		ExcludeNamespaces: []string{"kube-system"},
	}

	pods := []v1.Pod{
		newPod("default", "name1"),
		newPod("kube-system", "name2"),
	}

	filteredPods := filter.FilterExcludeNamespace(pods)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name1", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
}

func TestFilter_FilterExcludeNamespace_Empty(t *testing.T) {
	filter := Filter{
		ExcludeNamespaces: []string{},
	}

	pods := []v1.Pod{
		newPod("default", "name1"),
	}

	filteredPods := filter.FilterIncludeNamespace(pods)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name1", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
}

func TestFilter_FilterExcludeAnnotations(t *testing.T) {
	filter := Filter{
		ExcludeAnnotations: map[string]string{ "testkey1": "testvalue1" },
	}

	pods := []v1.Pod{
		newPodWithAnnotations("default", "name1", map[string]string{ "testkey1": "testvalue1" }),
		newPodWithAnnotations("default", "name2", map[string]string{ "testkey1": "testvalue2" }),
		newPodWithAnnotations("default", "name3", map[string]string{ "testkey2": "testvalue1" }),
	}

	filteredPods := filter.FilterExcludeAnnotations(pods)
	assert.Len(t, filteredPods, 2)
	assert.Equal(t, "name2", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
	assert.Equal(t, "name3", filteredPods[1].Name)
	assert.Equal(t, "default", filteredPods[1].Namespace)
}

func TestFilter_FilterExcludeAnnotations_Empty(t *testing.T) {
	filter := Filter{
		ExcludeAnnotations: map[string]string{},
	}

	pods := []v1.Pod{
		newPodWithAnnotations("default", "name1", map[string]string{ "testkey1": "testvalue1" }),
	}

	filteredPods := filter.FilterExcludeAnnotations(pods)
	assert.Len(t, filteredPods, 1)
}

func TestFilter_FilterPods(t *testing.T) {
	filter := Filter{}

	pods := []v1.Pod{
		newPod("default", "name1"),
	}

	filteredPods := filter.FilterPods(pods)
	assert.Len(t, filteredPods, 1)
}