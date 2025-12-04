# dev-template

The natron Stack Template for bootstrapping new applications. It is based on the basic implementation of [koda](https://github.com/natrontech/koda).

## Setup

First you need to create `.env` files according to the `.env.example` files in the `services/orchestrator` and `services/portal` directories.

Then you can start the services with `make run-orchestrator` and `make run-portal`.

You need to have a Kubernetes cluster with a working kubeconfig file. For development purposes you can use `minikube` or `kind`.
