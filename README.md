# Kube conformity

[![Build Status](https://travis-ci.org/stijndehaes/kube-conformity.svg?branch=master)](https://travis-ci.org/stijndehaes/kube-conformity)
[![Coverage Status](https://coveralls.io/repos/github/stijndehaes/kube-conformity/badge.svg?branch=master)](https://coveralls.io/github/stijndehaes/kube-conformity?branch=master)
[![Docker Hub](https://img.shields.io/docker/build/sdehaes/kube-conformity.svg)](https://hub.docker.com/r/sdehaes/kube-conformity/)
[![GitHub release](https://img.shields.io/github/release/stijndehaes/kube-conformity.svg)](https://github.com/stijndehaes/kube-conformity/releases)

This project looks for pods that in your kubernetes cluster that are breaking conformity rules.
An example on how to run this project on kubernetes can be found in the examples folder.

# Rules

At this moment there are three rules defined:

* Labels filled in: Takes a list of labels and check if pods have these labels defined
* Resource requests filled in: Checks all pods if they have resource requests filled in
* Limits requests filled in: Checks all pods if they have limits requests filled in

The rules are configured using a yaml config. An example of this config is:

```yaml
interval: 1h
labels_filled_in_rules:
- name: Check if label app is active on every pod except in kube-system
  labels:
  - app
  filter:
    exclude_namespaces:
    - kube-system
limits_filled_in_rules:
- name: Checks if limits are filled in everywhere except in kube-system
  filter:
    exclude_namespaces:
    - kube-system
requests_filled_in_rules:
- name: Checks if requests are filled in everywhere
  filter:
    exclude_namespaces:
    - kube-system
```

# Filtering
Each rule can be filtered on the base of four fields:

```yaml
- name: Checks if requests are filled in everywhere
  filter:
    include_namespaces
    - kube-system
    exclude_namespaces:
    - kube-system
    exclude_annotations:
      annotationKey: AnnotationValue
    exlcude_labels:
      labelKey: labelValue
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