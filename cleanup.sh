#!/usr/bin/env bash

kubectl delete namespace webhook-demo
kubectl delete namespace dev
kubectl delete validatingwebhookconfigurations.admissionregistration.k8s.io webhook-demo
