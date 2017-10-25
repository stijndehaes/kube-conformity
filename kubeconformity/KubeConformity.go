package kubeconformity

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"fmt"
	"github.com/stijndehaes/kube-conformity/config"
	"github.com/stijndehaes/kube-conformity/rules"
)

type KubeConformity struct {
	Client               kubernetes.Interface
	Logger               log.StdLogger
	KubeConformityConfig config.Config
}

func New(client kubernetes.Interface, logger log.StdLogger, config config.Config) *KubeConformity {
	return &KubeConformity{
		Client:               client,
		Logger:               logger,
		KubeConformityConfig: config,
	}
}

func (k *KubeConformity) LogNonConformingPods() error {
	conformityResults := k.EvaluateRules()
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
func (k *KubeConformity) EvaluateRules() []rules.RuleResult {

	podList, err := k.Client.CoreV1().Pods(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		k.Logger.Fatal(err)
	}

	ruleResults := []rules.RuleResult{}
	for _, rule := range k.KubeConformityConfig.RequestsFilledInRules {
		result := rule.FindNonConformingPods(podList.Items)
		ruleResults = append(ruleResults, result)
	}
	for _, rule := range k.KubeConformityConfig.LimitsFilledInRules {
		result := rule.FindNonConformingPods(podList.Items)
		ruleResults = append(ruleResults, result)
	}
	for _, rule := range k.KubeConformityConfig.LabelsFilledInRules {
		result := rule.FindNonConformingPods(podList.Items)
		ruleResults = append(ruleResults, result)
	}
	return ruleResults
}
