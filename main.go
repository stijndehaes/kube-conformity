package main

import (
	"fmt"
	"github.com/stijndehaes/kube-conformity/config"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/stijndehaes/kube-conformity/kubeconformity"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	master            = *kingpin.Flag("master", "The address of the Kubernetes cluster to target").String()
	kubeConfig        = *kingpin.Flag("kube-config", "Path to a kubeConfig file").String()
	debug             = *kingpin.Flag("debug", "Enable debug logging.").Bool()
	configLocation    = *kingpin.Flag("config-location", "The location of the config.yaml").Default("config.yaml").String()
	jsonLogging       = *kingpin.Flag("json-logging", "Enable json logging.").Bool()
	prometheusEnabled = *kingpin.Flag("prometheus-enabled", "Enable prometheus metrics").Default("true").Bool()
	PrometheusAddr    = *kingpin.Flag("prometheus-addr", "Prometheus metrics addr").Default(":8000").String()
)

func defaultPageHandler(config config.Config) func(w http.ResponseWriter, r *http.Request) {
	configByte, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatal(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
					<head><title>Kube conformity</title></head>
					<body>
					<h1>Kube conformity</h1>
					<p><a href="/metrics">Metrics</a></p>
					<p><a href="/healthz">Health Check</a></p>
					<h2>Configuration</h2>
					<p style='white-space: pre-wrap;'>`))
		w.Write(configByte)
		w.Write([]byte(`</p>
						</body>
					    </html>`))
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func configurePrometheus(config config.Config) {
	log.Info("Prometheus enabled will run it on addr: ", PrometheusAddr)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/", defaultPageHandler(config))
	go func() {
		if err := http.ListenAndServe(PrometheusAddr, nil); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("failed to start HTTP server")
		}
	}()
}

func main() {
	kingpin.Parse()
	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}
	ConfigureLogging()
	config, err := ConstructConfig()
	if err != nil {
		log.Fatal(err)
	}

	if prometheusEnabled {
		configurePrometheus(config)
	}

	kubeConformity := kubeconformity.New(
		client,
		log.StandardLogger(),
		config,
	)

	for {
		err := kubeConformity.LogNonConforming()
		if err != nil {
			log.Fatal(err)
		}

		log.Debugf("Sleeping for %s...", kubeConformity.KubeConformityConfig.Interval)
		time.Sleep(kubeConformity.KubeConformityConfig.Interval)
	}
}

func ConfigureLogging() {
	if jsonLogging {
		log.SetFormatter(&log.JSONFormatter{})
		log.Info("Json logging enabled")
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
	if debug {
		log.Info("Debug level enabled")
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func ConstructConfig() (config.Config, error) {
	kubeConformityConfig := config.Config{}
	yamlFile, err := ioutil.ReadFile(configLocation)
	if err != nil {
		return kubeConformityConfig, err
	}
	err = yaml.Unmarshal(yamlFile, &kubeConformityConfig)
	if err != nil {
		return kubeConformityConfig, err
	}
	return kubeConformityConfig, nil
}

func newClient() (*kubernetes.Clientset, error) {
	if _, err := os.Stat(clientcmd.RecommendedHomeFile); kubeConfig == "" && err == nil {
		kubeConfig = clientcmd.RecommendedHomeFile
	}
	kConfig, err := clientcmd.BuildConfigFromFlags(master, kubeConfig)
	if err != nil {
		return nil, err
	}
	log.Infof("Targeting cluster at %s", kConfig.Host)
	client, err := kubernetes.NewForConfig(kConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}
