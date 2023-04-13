# Kubernetes Dynamic Admission Control

Demo of Kubernetes Dynamic Admission Control exploring various concepts, capabilities and constraints ( focusing on validation admission webhooks ).

[Guide](./minikube-setup.md) for running this demo locally on minikube.

## Observations and notes

### Failure Policy

If you don't want to put your webhook on the critical path, meaning if you webhook server is down you won't be able to schedule any resource, you
can configure [Failure Policy](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#failure-policy).

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

## Resources

- [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers)
- [A Guide to Kubernetes Admission Controllers](https://kubernetes.io/blog/2019/03/21/a-guide-to-kubernetes-admission-controllers/)
