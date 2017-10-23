package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/stijndehaes/kube-conformity/kubeconformity"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/stijndehaes/kube-conformity/config"
)

var (
	master         string
	kubeconfig     string
	interval       time.Duration
	debug          bool
	version        string
	configLocation string
)

func init() {
	kingpin.Flag("master", "The address of the Kubernetes cluster to target").StringVar(&master)
	kingpin.Flag("kubeconfig", "Path to a kubeconfig file").StringVar(&kubeconfig)
	kingpin.Flag("interval", "Interval between conformity checks").Default("1h").DurationVar(&interval)
	kingpin.Flag("debug", "Enable debug logging.").BoolVar(&debug)
	kingpin.Flag("config-location", "The location of the config.yaml").Default("config.yaml").StringVar(&configLocation)
}

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}

	kubeConformity := kubeconformity.New(
		client,
		log.StandardLogger(),
		ConstructConfig(),
	)

	for {
		err := kubeConformity.LogNonConformingPods()
		if err != nil {
			log.Fatal(err)
		}

		log.Debugf("Sleeping for %s...", interval)
		time.Sleep(interval)
	}
}

func ConstructConfig() config.KubeConformityConfig {
	kubeConfig := config.KubeConformityConfig{}

	yamlFile, err := ioutil.ReadFile(configLocation)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &kubeConfig)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return kubeConfig
}

func newClient() (*kubernetes.Clientset, error) {
	if kubeconfig == "" {
		if _, err := os.Stat(clientcmd.RecommendedHomeFile); err == nil {
			kubeconfig = clientcmd.RecommendedHomeFile
		}
	}

	kconfig, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		return nil, err
	}

	log.Infof("Targeting cluster at %s", kconfig.Host)

	client, err := kubernetes.NewForConfig(kconfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
