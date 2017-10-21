package kubeconformity

import (
	log "github.com/sirupsen/logrus"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"fmt"
)

type KubeConformity struct {
	Client kubernetes.Interface
	Logger log.StdLogger
	Labels []string
}

type KubeConformityResult struct {
	ResourceProblems 	[]v1.Pod
	LabelProblems 		[]v1.Pod
}


type PrintableKubeConformityResult struct {
	ResourceProblems 	[]string
	LabelProblems 		[]string
}

func New(client kubernetes.Interface, logger log.StdLogger, labels []string) *KubeConformity {
	return &KubeConformity{
		Client: client,
		Logger: logger,
		Labels: labels,
	}
}


func (k *KubeConformity) LogNonConformingPods() error {
	conformityResult, err := k.FindNonConformingPods()
	if err != nil {
		return err
	}
	resourceProblemString := []string{}
	for _, pod := range conformityResult.ResourceProblems {
		resourceProblemString = append(resourceProblemString, fmt.Sprintf("%s_%s(%s)", pod.Name, pod.Namespace, pod.UID))
	}
	labelProblemsString := []string{}
	for _, pod := range conformityResult.LabelProblems {
		labelProblemsString = append(labelProblemsString, fmt.Sprintf("%s_%s(%s)", pod.Name, pod.Namespace, pod.UID))
	}

	k.Logger.Print(PrintableKubeConformityResult{resourceProblemString, labelProblemsString})
	return nil
}

// Candidates returns the list of pods that are available for termination.
// It returns all pods matching the label selector and at least one namespace.
func (k *KubeConformity) FindNonConformingPods() (KubeConformityResult, error) {

	podList, err := k.Client.Core().Pods(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return KubeConformityResult{}, err
	}

	resourcePods, err := filterOnResources(podList.Items)
	if err != nil {
		return KubeConformityResult{}, err
	}

	labelPods, err := filterOnLabels(podList.Items, k.Labels)
	if err != nil {
		return KubeConformityResult{}, err
	}

	return KubeConformityResult{resourcePods, labelPods}, nil
}

func filterOnResources(pods []v1.Pod) ([]v1.Pod, error) {

	filteredList := []v1.Pod{}
	for _, pod := range pods {
		var podNonConform = false

		for _, container := range pod.Spec.Containers {
			podNonConform = podNonConform || container.Resources.Limits.Cpu().IsZero()
			podNonConform = podNonConform || container.Resources.Limits.Memory().IsZero()
			podNonConform = podNonConform || container.Resources.Requests.Cpu().IsZero()
			podNonConform = podNonConform || container.Resources.Requests.Memory().IsZero()
		}

		if podNonConform {
			filteredList = append(filteredList, pod)
		}
	}

	return filteredList, nil
}

func filterOnLabels(pods []v1.Pod, labels []string) ([]v1.Pod, error) {

	filteredList := []v1.Pod{}
	if len(labels) == 0 {
		return filteredList, nil
	}
	for _, pod := range pods {
		for _, label := range labels {
			containsLabel := false
			for podLabelKey, _ := range pod.ObjectMeta.Labels {
				if podLabelKey == label {
					containsLabel = true
				}
			}
			if !containsLabel {
				filteredList = append(filteredList, pod)
				break
			}
		}
	}

	return filteredList, nil
}
