package kubeconformity

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"fmt"
)

type RuleResult struct {
	Pods     []v1.Pod
	Reason   string
	RuleName string
}

type KubeConformity struct {
	Client               kubernetes.Interface
	Logger               log.StdLogger
	KubeConformityConfig KubeConformityConfig
}

func New(client kubernetes.Interface, logger log.StdLogger, config KubeConformityConfig) *KubeConformity {
	return &KubeConformity{
		Client:               client,
		Logger:               logger,
		KubeConformityConfig: config,
	}
}

func (k *KubeConformity) LogNonConformingPods() error {
	conformityResults, err := k.EvaluateRules()
	if err != nil {
		return err
	}

	for _, ruleResult := range conformityResults {
		k.Logger.Printf("rule name: %s", ruleResult.RuleName)
		k.Logger.Printf("rule reason: %s", ruleResult.Reason)
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
	for _, rule := range k.KubeConformityConfig.RequestsFilledInRules {
		result := rule.findNonConformingPods(podList.Items)
		ruleResults = append(ruleResults, result)
	}
	for _, rule := range k.KubeConformityConfig.LimitsFilledInRules {
		result := rule.findNonConformingPods(podList.Items)
		ruleResults = append(ruleResults, result)
	}
	for _, rule := range k.KubeConformityConfig.LabelsFilledInRules {
		result := rule.findNonConformingPods(podList.Items)
		ruleResults = append(ruleResults, result)
	}
	return ruleResults, nil
}
