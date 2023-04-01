# Kubernetes Dynamic Admission Control

Demo of Kubernetes Dynamic Admission Control

## Prerequisites

### Resources

- [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers)
- [A Guide to Kubernetes Admission Controllers](https://kubernetes.io/blog/2019/03/21/a-guide-to-kubernetes-admission-controllers/)


2m35s       Warning   FailedCreate        replicaset/webhook-server-79c79b5877   Error creating: Internal error occurred: failed calling webhook "webhook-server.webhook-demo.svc": failed to call webhook: Post "https://webhook-server.webhook-demo.svc:443/validate?timeout=10s": dial tcp 10.109.147.187:443: connect: connection refused
13s         Warning   FailedCreate        replicaset/webhook-server-79c79b5877   Error creating: Internal error occurred: failed calling webhook "webhook-server.webhook-demo.svc": failed to call webhook: Post "https://webhook-server.webhook-demo.svc:443/validate?timeout=10s": dial tcp 10.109.147.187:443: connect: connection refused

<!-- NOTE: document minikube setup -->
<!-- NOTE: update setup for admission webhook -->
