# Kube conformity

[![Build Status](https://travis-ci.org/stijndehaes/kube-conformity.svg?branch=master)](https://travis-ci.org/stijndehaes/kube-conformity)
[![Coverage Status](https://coveralls.io/repos/github/stijndehaes/kube-conformity/badge.svg?branch=master)](https://coveralls.io/github/stijndehaes/kube-conformity?branch=master)
[![Docker Hub](https://img.shields.io/docker/build/sdehaes/kube-conformity.svg)](https://hub.docker.com/r/sdehaes/kube-conformity/)
[![GitHub release](https://img.shields.io/github/release/stijndehaes/kube-conformity.svg)](https://github.com/stijndehaes/kube-conformity/releases)

This project looks for objects that in your kubernetes cluster that are breaking conformity rules.
An example on how to run this project on kubernetes can be found in the examples folder.

# Rules

At the moment rules on 3 resources are supported.

## Pod rules

* Labels filled in: Takes a list of labels and check if pods have these labels defined
* Resource requests filled in: Checks all pods if they have resource requests filled in
* Limits requests filled in: Checks all pods if they have limits requests filled in

## Deployment rules

* Deployment minimum replicas: Checks that every Deployment has a minimum of a certain number of replicas

## StatefulSet Rules 

* StatefulSet minimum replicas: Checks that every StatefulSet has a minimum of a certain number of replicas


The rules are configured using a yaml config. An example of this config is:

```yaml
interval: 1h
pod_rules_labels_filled_in:
- name: Check if label app is active on every pod except in kube-system
  labels:
  - app
  filter:
    exclude_namespaces:
    - kube-system
pod_rules_limits_filled_in:
- name: Checks if pod limits are filled in everywhere except in kube-system
  filter:
    exclude_namespaces:
    - kube-system
pod_rules_requests_filled_in:
- name: Checks if pod requests are filled in everywhere
  filter:
    exclude_namespaces:
    - kube-system
deployment_rules_replicas_minimum:
- name: Checks that al Deployments have a minimum of 2 replicas
  minimum_replicas: 2
  filter:
    exclude_namespaces:
    - kube-system
stateful_set_rules_replicas_minimum:
- name: Checks that al StatefulSets have a minimum of 2 replicas
  minimum_replicas: 2
  filter:
    exclude_namespaces:
    - kube-system
```

# Filtering
Each rule can be filtered on the base of five fields:

* include_namespaces: A list of namespaces to include, if empty defaults to all namespaces
* exclude_namespaces: A list of namespaces to exclude, if empty defaults to none
* exclude_annotations: A map of annotations to exclude
* exclude_labels: A map of labels to exclude
* exclude_jobs (Only available on the three pod rules): Excludes pod created by a job, filters on the labelkey `job-name`

An example of the yaml configuration:

```yaml
- name: Checks if requests are filled in everywhere
  filter:
    include_namespaces:
    - kube-system
    exclude_namespaces:
    - kube-system
    exclude_annotations:
      annotationKey: AnnotationValue
    exclude_labels:
      labelKey: labelValue
    exclude_jobs: true
```


# Email config
Default the non-conforming pods get logged to stdout.
But it is also possible to have these reports send through email.
The auth_password variable can also be set by and environment variable called: `CONFORMITY_EMAIL_AUTH_PASSWORD`.

An example of a full email config would be.

```yaml
interval: 1h
limits_filled_in_rules:
- name: Checks if limits are filled in everywhere
email_config:
  enabled: true
  to: test@gmail.com
  from: no-reply@kube-conformity.com
  host: 127.0.0.1
  port: 24
  subject: kube-conformity
  auth_username: username
  auth_password: password
  auth_identity: identity
  template: mailtemplate.html
```

Not all values have to be filled in. The following table denotes if a value is required and the default value if it has any.

| Value         | default                       | required  |
| ------------- | ----------------------------- | --------- |
| enabled       | false                         | false     |
| to            |                               | true      |
| from          | no-reply@kube-conformity.com  | true      |
| host          |                               | true      |
| port          | 24                            | true      |
| subject       | kube-conformity               | false     |
| auth_username |                               | false     |
| auth_password |                               | false     |
| auth_identity |                               | false     |
| template      | mailtemplate.html             | true      |

# Command line arguments

Some of the setup is done through command line arguments the arguments available are:

* --master=address : The address of the Kubernetes cluster to target
* --kube-config=path : Path to a kubeConfig file, default = kube config in home directory
* --debug : Enable debug logging.
* --json-logging : Enable json logging.
* --config-location=path : The location of the config.yaml, default = config.yaml

When running in the cluster the kube-config file or master address should be picked up automatically.