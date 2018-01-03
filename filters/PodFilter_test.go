package filters

import (
	"testing"
	"k8s.io/client-go/pkg/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestPodFilter_FilterExcludeNamespace(t *testing.T) {
	filter := PodFilter{
		ExcludeNamespaces: []string{"kube-system"},
	}

	pods := []v1.Pod{
		newPod("default", "name1", "uid1"),
		newPod("kube-system", "name2", "uid2"),
	}

	filteredPods := filter.FilterExcludeNamespace(pods)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name1", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
}

func TestPodFilter_FilterExcludeNamespace_Empty(t *testing.T) {
	filter := PodFilter{
		ExcludeNamespaces: []string{},
	}

	pods := []v1.Pod{
		newPod("default", "name1", "uid1"),
	}

	filteredPods := filter.FilterIncludeNamespace(pods)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name1", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
}

func TestPodFilter_FilterExcludeAnnotations(t *testing.T) {
	filter := PodFilter{
		ExcludeAnnotations: map[string]string{"testkey1": "testvalue1"},
	}

	pods := []v1.Pod{
		newPodWithAnnotations("default", "name1", "uid1", map[string]string{"testkey1": "testvalue1"}),
		newPodWithAnnotations("default", "name2", "uid2", map[string]string{"testkey1": "testvalue2"}),
		newPodWithAnnotations("default", "name3", "uid3", map[string]string{"testkey2": "testvalue1"}),
	}

	filteredPods := filter.FilterExcludeAnnotations(pods)
	assert.Len(t, filteredPods, 2)
	assert.Equal(t, "name2", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
	assert.Equal(t, "name3", filteredPods[1].Name)
	assert.Equal(t, "default", filteredPods[1].Namespace)
}

func TestPodFilter_FilterExcludeAnnotations_Empty(t *testing.T) {
	filter := PodFilter{
		ExcludeAnnotations: map[string]string{},
	}

	pods := []v1.Pod{
		newPodWithAnnotations("default", "name1", "uid1", map[string]string{"testkey1": "testvalue1"}),
	}

	filteredPods := filter.FilterExcludeAnnotations(pods)
	assert.Len(t, filteredPods, 1)
}

func TestPodFilter_FilterExcludeLabels(t *testing.T) {
	filter := PodFilter{
		ExcludeLabels: map[string]string{"testkey1": "testvalue1"},
	}

	pods := []v1.Pod{
		newPodWithLabels("default", "name1", "uid1", map[string]string{"testkey1": "testvalue1"}),
		newPodWithLabels("default", "name2", "uid2", map[string]string{"testkey1": "testvalue2"}),
		newPodWithLabels("default", "name3", "uid3", map[string]string{"testkey2": "testvalue1"}),
	}

	filteredPods := filter.FilterExcludeLabels(pods)
	assert.Len(t, filteredPods, 2)
	assert.Equal(t, "name2", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
	assert.Equal(t, "name3", filteredPods[1].Name)
	assert.Equal(t, "default", filteredPods[1].Namespace)
}

func TestPodFilter_FilterExcludeLabels_Empty(t *testing.T) {
	filter := PodFilter{
		ExcludeLabels: map[string]string{},
	}

	pods := []v1.Pod{
		newPodWithLabels("default", "name1", "uid1", map[string]string{"testkey1": "testvalue1"}),
	}

	filteredPods := filter.FilterExcludeLabels(pods)
	assert.Len(t, filteredPods, 1)
}

func TestPodFilter_FilterExcludeJobs_true(t *testing.T) {
	filter := PodFilter{
		ExcludeJobs: true,
	}

	pods := []v1.Pod{
		newPodWithLabels("default", "name1", "uid1", map[string]string{"job-name": "curator-212312"}),
		newPodWithLabels("default", "name2", "uid2", map[string]string{}),
	}

	filteredPods := filter.FilterExcludeJobs(pods)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name2", filteredPods[0].Name)
}

func TestPodFilter_FilterExcludeJobs_false(t *testing.T) {
	filter := PodFilter{
		ExcludeJobs: false,
	}

	pods := []v1.Pod{
		newPodWithLabels("default", "name1", "uid1", map[string]string{"job-name": "curator-212312"}),
		newPodWithLabels("default", "name2", "uid2", map[string]string{}),
	}

	filteredPods := filter.FilterExcludeJobs(pods)
	assert.Len(t, filteredPods, 2)
}

func TestPodFilter_FilterPods2(t *testing.T) {
	filter := PodFilter{}

	pods := []v1.Pod{
		newPod("default", "name1", "uid1"),
	}

	filteredPods := filter.FilterPods(pods)
	assert.Len(t, filteredPods, 1)
}
