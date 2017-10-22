package kubeconformity

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"fmt"
)

type Rule interface {
	findNonConformingPods(pods []v1.Pod) RuleResult
}

type RuleResult struct {
	Pods   []v1.Pod
	Reason string
}

type KubeConformity struct {
	Client kubernetes.Interface
	Logger log.StdLogger
	Rules  []Rule
}

func New(client kubernetes.Interface, logger log.StdLogger, rules []Rule) *KubeConformity {
	return &KubeConformity{
		Client: client,
		Logger: logger,
		Rules: rules,
	}
}

func (k *KubeConformity) LogNonConformingPods() error {
	conformityResults, err := k.EvaluateRules()
	if err != nil {
		return err
	}

	for _, ruleResult := range conformityResults {
		k.Logger.Print(ruleResult.Reason)
		for _, pod := range ruleResult.Pods {
			k.Logger.Print(fmt.Sprintf("%s_%s(%s)", pod.Name, pod.Namespace, pod.UID))
		}
	}
	return nil
}

// Candidates returns the list of pods that are available for termination.
// It returns all pods matching the label selector and at least one namespace.
func (k *KubeConformity) EvaluateRules() ([]RuleResult, error) {

	podList, err := k.Client.CoreV1().Pods(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return []RuleResult{}, err
	}

	ruleResults := []RuleResult{}
	for _, rule := range k.Rules {
		result := rule.findNonConformingPods(podList.Items)
		ruleResults = append(ruleResults, result)
	}
	return ruleResults, nil
}
