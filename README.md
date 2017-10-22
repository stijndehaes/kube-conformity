# Kube conformity

[![Build Status](https://travis-ci.org/stijndehaes/kube-conformity.svg?branch=master)](https://travis-ci.org/stijndehaes/kube-conformity) [![Coverage Status](https://coveralls.io/repos/github/stijndehaes/kube-conformity/badge.svg?branch=master)](https://coveralls.io/github/stijndehaes/kube-conformity?branch=master) [![Docker Hub](https://img.shields.io/docker/build/sdehaes/kube-conformity.svg)](https://hub.docker.com/r/sdehaes/kube-conformity/)

This project looks for pods that in your kubernetes cluster that are breaking conformity rules.

At this moment there are three rules defined:

* Labels filled in: Takes a list of labels and check if pods have these labels defined
* Resource requests filled in: Checks all pods if they have resource requests filled in
* Limits requests filled in: Checks all pods if they have limits requests filled in

To run this project you can use the provided docker image. And example deployment is the following:

```bash
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: kube-conformity
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: kube-conformity
    spec:
      containers:
      - name: kube-conformity
        image: sdehaes/kube-conformity:latest
        args:
        - --interval=4h
        - --request-check
        - --limits-check
        - --labels="app"
```