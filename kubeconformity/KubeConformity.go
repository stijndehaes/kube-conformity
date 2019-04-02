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
	statefulSetRuleResults := k.EvaluateStatefulSetRules()
	if len(podRuleResults) > 0 {
		k.Logger.Println(fmt.Sprint("Presenting Pod rule results"))
	}
	for _, ruleResult := range podRuleResults {
		k.Logger.Println(fmt.Sprintf("rule name: %s", ruleResult.RuleName))
		k.Logger.Println(fmt.Sprintf("rule reason: %s", ruleResult.Reason))
		for _, pod := range ruleResult.Pods {
			k.Logger.Println(fmt.Sprintf("%s_%s", pod.Name, pod.Namespace))
		}
	}
	if len(deploymentRuleResults) > 0 {
		k.Logger.Println(fmt.Sprint("Presenting Deployment rule results"))
	}
	for _, ruleResult := range deploymentRuleResults {
		k.Logger.Println(fmt.Sprintf("rule name: %s", ruleResult.RuleName))
		k.Logger.Println(fmt.Sprintf("rule reason: %s", ruleResult.Reason))
		for _, deployment := range ruleResult.Deployments {
			k.Logger.Println(fmt.Sprintf("%s_%s", deployment.Name, deployment.Namespace))
		}
	}
	if len(statefulSetRuleResults) > 0 {
		k.Logger.Println(fmt.Sprint("Presenting StatefulSet rule results"))
	}
	for _, ruleResult := range statefulSetRuleResults {
		k.Logger.Println(fmt.Sprintf("rule name: %s", ruleResult.RuleName))
		k.Logger.Println(fmt.Sprintf("rule reason: %s", ruleResult.Reason))
		for _, statefulSet := range ruleResult.StatefulSets {
			k.Logger.Println(fmt.Sprintf("%s_%s", statefulSet.Name, statefulSet.Namespace))
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

func (k *KubeConformity) EvaluateStatefulSetRules() []rules.StatefulSetRuleResult {
	statefulSetList, err := k.Client.AppsV1().StatefulSets(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		k.Logger.Fatal(err)
	}
	var ruleResults []rules.StatefulSetRuleResult
	for _, rule := range k.KubeConformityConfig.StatefulSetRuleReplicasMinimum {
		result := rule.FindNonConformingStatefulSet(statefulSetList.Items)
		ruleResults = append(ruleResults, result)
	}
	return ruleResults
}