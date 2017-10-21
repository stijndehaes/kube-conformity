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

func New(client kubernetes.Interface, logger log.StdLogger) *KubeConformity {
	return &KubeConformity{
		Client:      client,
		Logger:      logger,
	}
}

func (k *KubeConformity) LogNonConformingPods() error {
	pods, err := k.FindNonConformingPods()
	if err != nil {
		return err
	}
	for _, pod := range pods {
		k.Logger.Print(pod)
	}
	return nil
}

// Candidates returns the list of pods that are available for termination.
// It returns all pods matching the label selector and at least one namespace.
func (k *KubeConformity) FindNonConformingPods() ([]v1.Pod, error) {

	podList, err := k.Client.Core().Pods(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pods, err := filterConformingPods(podList.Items)
	if err != nil {
		return nil, err
	}

	return pods, nil
}

func filterConformingPods(pods []v1.Pod) ([]v1.Pod, error) {

	filteredList := []v1.Pod{}
	for _, pod := range pods {
		var podNonConform = false

		for _, container := range pod.Spec.Containers {
			log.Debug(container.Resources.Limits.Cpu().IsZero())
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