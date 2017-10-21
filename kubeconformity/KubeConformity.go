package kubeconformity

import (

log "github.com/sirupsen/logrus"

metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
"k8s.io/client-go/kubernetes"
"k8s.io/client-go/pkg/api/v1"
)

type KubeConformity struct {
	Client kubernetes.Interface
	Logger log.StdLogger
}

type KubeConformityResult struct {
	ResourceProblems []v1.Pod
}

func New(client kubernetes.Interface, logger log.StdLogger) *KubeConformity {
	return &KubeConformity{
		Client:      client,
		Logger:      logger,
	}
}

func (k *KubeConformity) LogNonConformingPods() error {
	conformityResult, err := k.FindNonConformingPods()
	if err != nil {
		return err
	}
	k.Logger.Print(conformityResult)
	return nil
}

// Candidates returns the list of pods that are available for termination.
// It returns all pods matching the label selector and at least one namespace.
func (k *KubeConformity) FindNonConformingPods() (KubeConformityResult, error) {

	podList, err := k.Client.Core().Pods(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return KubeConformityResult{}, err
	}

	pods, err := filterOnResources(podList.Items)
	if err != nil {
		return KubeConformityResult{}, err
	}

	return KubeConformityResult{pods}, nil
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