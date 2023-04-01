apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: webhook-demo
webhooks:
  - name: webhook-server.webhook-demo.svc
    sideEffects: None
    admissionReviewVersions: ["v1", "v1beta1"]
    clientConfig:
      service:
        name: webhook-server
        namespace: webhook-demo
        path: "/validate"
      caBundle: ${CA_PEM_B64}
    rules:
      # TODO: failure policy
      # TODO: namespaces to skip
      - operations: ["CREATE"]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["pods"]
