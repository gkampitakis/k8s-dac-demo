# Minikube Setup

This is a guide for running and testing the DAC demo locally.

## Prerequisites

You will need some tools installed.

- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
- [docker](https://docs.docker.com/engine/)
- [minikube](https://minikube.sigs.k8s.io/docs/start/)
- [go](https://go.dev/doc/install) (_optional if you want to experiment with webhook server_ )

## Get started

```sh
# start kubernetes cluster locally
minikube start

# enable needed addons
minikube addons enable metrics-server
minikube addons enable registry
```

We are setting up a local image registry so we don't rely on pushing/pulling the webhook server container image on a remote registry e.g. [docker hub](https://hub.docker.com/).

```sh
# run local registry
docker run --rm -it --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:$(minikube ip):5000"
```

All the components for running and testing the webhook now should be ready.

## Building and Deploying

The next steps are for building and deploying the webhook server alongside the webhook configuration.

1. You will need to build and push the container to your local registry. You can do that by running

```sh
make docker-image-local
```

2. You are ready to deploy the validation webhook in your local kubernetes cluster

```sh
make deploy
```

## Verifying everything works

You should see two healthy running pods from the webhook-demo deployment with

```sh
kubectl get pods -n webhook-demo
```

example output

```
NAME                              READY   STATUS    RESTARTS   AGE
webhook-server-79c79b5877-d624n   1/1     Running   0          11m
webhook-server-79c79b5877-dpf4b   1/1     Running   0          11m
```

## Experimenting and applying changes

For experimenting/changing the [webhook-server](./webhook-server/main.go), you can
build a new docker image and push it on your local registry with

```sh
make docker-image-local
```

Then you need to restart the k8s deployment for pulling the new changes with

```sh
kubectl rollout restart deploy/webhook-server
```

On the other hand for changing the [ValidatingWebhookConfiguration](./deployment/webhookConfig.yaml.tpl) you need to re-deploy.

You can do that with

```sh
make clean

make deploy
```
