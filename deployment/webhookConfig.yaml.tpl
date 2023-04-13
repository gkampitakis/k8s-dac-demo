apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: webhook-demo
webhooks:
  - name: webhook-server.webhook-demo.svc
    failurePolicy: Ignore # or Fail more information at https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#failure-policy
    sideEffects: None
    admissionReviewVersions: ["v1", "v1beta1"]
    clientConfig:
      service:
        name: webhook-server
        namespace: webhook-demo
        path: "/validate"
      caBundle: ${CA_PEM_B64}
    rules:
      - operations: ["CREATE"]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["pods"]
        scope: "Namespaced"
