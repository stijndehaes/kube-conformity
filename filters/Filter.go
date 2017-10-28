package filters

import (
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/apimachinery/pkg/labels"
	"log"
	"k8s.io/apimachinery/pkg/selection"
	"fmt"
	"github.com/emicklei/go-restful"
)

type Filter struct {
	NamespacesString string `yaml:"namespaces"`
}

func (f Filter) FilterPods(pods []v1.Pod) []v1.Pod {
	return f.filterOnNamespaces(pods)
}

func (f Filter) filterOnAnnotations(pods []v1.Pod) []v1.Pod {
	reqIncl, reqExcl, _ := ParseFilterString(f.NamespacesString)

	if len(reqIncl) == 0 && len(reqExcl) == 0 {
		return pods
	}

	filteredList:= []v1.Pod{}

	for _, pod := range pods {
		// if there aren't any including requirements, we're in by default
		included := len(reqIncl) == 0

		// convert the pod's namespace to an equivalent label selector
		selector := labels.Set{pod.Namespace: ""}

		// include pod if one including requirement matches
		for _, req := range reqIncl {
			if req.Matches(selector) {
				included = true
				break
			}
		}

		// exclude pod if it is filtered out by at least one excluding requirement
		for _, req := range reqExcl {
			if !req.Matches(selector) {
				included = false
				break
			}
		}

		if included {
			filteredList = append(filteredList, pod)
		}
	}

	return filteredList
}

func ParseFilterString(string string) ([]labels.Requirement, []labels.Requirement, error) {
	selectors, err := labels.Parse(string)
	if err != nil {
		log.Fatal(err)
	}
	if selectors.Empty() {
		return []labels.Requirement{}, []labels.Requirement{}, nil
	}

	// split requirements into including and excluding groups
	reqs, _ := selectors.Requirements()
	reqIncl := []labels.Requirement{}
	reqExcl := []labels.Requirement{}

	for _, req := range reqs {
		switch req.Operator() {
		case selection.Exists:
			reqIncl = append(reqIncl, req)
		case selection.DoesNotExist:
			reqExcl = append(reqExcl, req)
		default:
			return []labels.Requirement{}, []labels.Requirement{}, error(fmt.Errorf("unsupported operator: %s", req.Operator()))
		}
	}
	return reqIncl, reqExcl, nil
}

func (f Filter) filterOnNamespaces(pods []v1.Pod) []v1.Pod {
	reqIncl, reqExcl, _ := ParseFilterString(f.NamespacesString)

	if len(reqIncl) == 0 && len(reqExcl) == 0 {
		return pods
	}

	filteredList:= []v1.Pod{}

	for _, pod := range pods {
		// if there aren't any including requirements, we're in by default
		included := len(reqIncl) == 0

		// convert the pod's namespace to an equivalent label selector
		selector := labels.Set{pod.Namespace: ""}

		// include pod if one including requirement matches
		for _, req := range reqIncl {
			if req.Matches(selector) {
				included = true
				break
			}
		}

		// exclude pod if it is filtered out by at least one excluding requirement
		for _, req := range reqExcl {
			if !req.Matches(selector) {
				included = false
				break
			}
		}

		if included {
			filteredList = append(filteredList, pod)
		}
	}

	return filteredList
}
