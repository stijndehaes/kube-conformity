package kubeconformity

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/api/core/v1"
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

func (k *KubeConformity) LogNonConforming() error {
	podRuleResults := k.EvaluatePodRules()
	deploymentRuleResults := k.EvaluateDeploymentRules()
	for _, ruleResult := range podRuleResults {
		ruleName := fmt.Sprintf("rule name: %s", ruleResult.RuleName)
		k.Logger.Println(ruleName)
		ruleReason := fmt.Sprintf("rule reason: %s", ruleResult.Reason)
		k.Logger.Println(ruleReason)
		for _, pod := range ruleResult.Pods {
			podName := fmt.Sprintf("%s_%s", pod.Name, pod.Namespace)
			k.Logger.Println(podName)
		}
	}
	for _, ruleResult := range deploymentRuleResults {
		ruleName := fmt.Sprintf("rule name: %s", ruleResult.RuleName)
		k.Logger.Println(ruleName)
		ruleReason := fmt.Sprintf("rule reason: %s", ruleResult.Reason)
		k.Logger.Println(ruleReason)
		for _, deployment := range ruleResult.Deployments {
			podName := fmt.Sprintf("%s_%s", deployment.Name, deployment.Namespace)
			k.Logger.Println(podName)
		}
	}
	if k.KubeConformityConfig.EmailConfig.Enabled {
		k.Logger.Println("Sending mail with conformity results")
		return k.KubeConformityConfig.EmailConfig.SendMail(podRuleResults, deploymentRuleResults)
	}
	return nil
}

func (k *KubeConformity) EvaluatePodRules() []rules.PodRuleResult {
	podList, err := k.Client.CoreV1().Pods(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		k.Logger.Fatal(err)
	}

	var ruleResults []rules.PodRuleResult
	for _, rule := range k.KubeConformityConfig.PodRulesRequestsFilledIn {
		result := rule.FindNonConformingPods(podList.Items)
		ruleResults = append(ruleResults, result)
	}
	for _, rule := range k.KubeConformityConfig.PodRulesLimitsFilledIn {
		result := rule.FindNonConformingPods(podList.Items)
		ruleResults = append(ruleResults, result)
	}
	for _, rule := range k.KubeConformityConfig.PodRulesLabelsFilledIn {
		result := rule.FindNonConformingPods(podList.Items)
		ruleResults = append(ruleResults, result)
	}
	return ruleResults
}

func (k *KubeConformity) EvaluateDeploymentRules() []rules.DeploymentRuleResult {
	deploymentList, err := k.Client.AppsV1().Deployments(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		k.Logger.Fatal(err)
	}
	var ruleResults []rules.DeploymentRuleResult
	for _, rule := range k.KubeConformityConfig.DeploymentRuleReplicasMinimum {
		result := rule.FindNonConformingDeployment(deploymentList.Items)
		ruleResults = append(ruleResults, result)
	}
	return ruleResults
}