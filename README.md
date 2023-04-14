# Kubernetes Dynamic Admission Control

Demo of Kubernetes Dynamic Admission Control (DAC) exploring various concepts, capabilities and constraints ( focusing on validation admission webhooks ).

[Guide](./minikube-setup.md) for running this demo locally on minikube.

## Example

This repo contains an example of all moving components for creating and testing Dynamic Admission Control

- [ValidatingWebhookConfiguration](./deployment/webhookConfig.yaml.tpl)
- [Webhook Server](./deployment/deployment.yaml)
- [Test busybox pod](./deployment/busy-box.yaml) for experimenting and verifying DAC works.
- [Service](./deployment/svc.yaml)
- [Makefile](./Makefile) containing some helper commands
  - Script for setting up and deploying all components. `make deploy`
  - Script for tearing down all created resources. `make cleanup`
  - Script for building and pushing docker image in the local registry. `make docker-image`
- [Minikube guide](./minikube-setup.md) containing all steps to run this example locally on Minikube.

### Webhook server

Webhook-server allows some configurable behaviour, but you can tweak it and add your business and experiment with different scenarios.

- You can set env `ALLOW_SCHEDULING` to `"true"` and instead of an error the DAC will emit a warning

```sh
Warning: Team label not set on pod
```

- You can set env `SKIP_NAMESPACE` to a comma separated string for a list of namespace to skip. By default `"kube-public"` and `"kube-system"` are skipped.

## Observations and notes

### Namespace Selector

You can limit which requests are reaching your webhook-server depending on resources' namespaces [namespaceSelector](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#matching-requests-namespaceselector).

I couldn't find a way to whitelist namespaces by their name, except of building this logic inside the [webhook-server](./webhook-server/main.go#L117).

For the namespaceSelector you can test it with adding this to [webhook configuration](./deployment/webhookConfig.yaml.tpl)

```yaml
namespaceSelector:
  matchExpressions:
    - key: webhook-skip
      operator: DoesNotExist
```

and then add the label to the namespaces that you want to skip e.g. `webhook-demo`

```sh
kubectl label ns webhook-demo webhook-skip=true
```

### Reliability

Dynamic Admission Controller can be on the critical path, depending how you your Webhook configured e.g. can block various actions on Kubernetes resources.

There are couple of things someone can do increase the reliability and the availability of the admission controller.

_This is not an exhaustive list_

- The Dynamic Admission Controller is backed by a Server, that means all concepts for making a stateless deployment more reliable in Kubernetes (if it is hosted on the kubernetes)
  - Use load balancing techniques that provides high availability and performance benefits [[Availability](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#availability)].
  - Rolling update strategy and [Pod Disruption Budget](https://kubernetes.io/docs/tasks/run-application/configure-pdb/).
    > Rolling updates allow Deployments' update to take place with zero downtime by incrementally updating Pods instances with new ones
  - [Horizontal Pod Autoscaling](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) for matching increasing demand.
- [Avoiding deadlock in self hosted webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#avoiding-deadlocks-in-self-hosted-webhooks)
- [Avoiding operating on the kube-system namespace](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#avoiding-operating-on-the-kube-system-namespace)
- It is recommended that admission webhooks should evaluate as quickly as possible (typically in milliseconds), since they add to API request latency. It is encouraged to use a small [timeout for webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#timeouts).
- Configure [Failure Policy](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#failure-policy) depending your needs.
- The API server exposes Prometheus metrics from the /metrics endpoint, which can be used for monitoring and diagnosing API server status.

  e.g. `apiserver_admission_webhook_admission_duration_seconds_sum`, `apiserver_admission_webhook_rejection_count`

## Resources

- [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers)
- [A Guide to Kubernetes Admission Controllers](https://kubernetes.io/blog/2019/03/21/a-guide-to-kubernetes-admission-controllers/)
