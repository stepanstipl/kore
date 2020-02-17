# Clusterman

Clusterman is a binary built as part of Kore that provides:
1. Simple single deployment for bootstrapping Kore capabilities into kubernetes clusters
1. A capability for reporting the status of all Kore provided components / **cluster apps** in a Kubernetes cluster (**coming soon** - see [appvia/kore/#89](https://github.com/appvia/kore/issues/89)

`kore-clusterman` is deployed automatically by the Kore API server in a new Kubernetes cluster after it is created.

It is delivered as part of the Kore container.

## Kore Cluster Components

Each Kore component that is deployed to a Kore managed Kubernetes cluster falls into the following categories:

1. [An Embedded Kore Component](#embeded-kore-components)
1. [Upstream Kore Components](#upstream-kore-components)

### Embeded Kore Components

These are qualified as:
- Built as part of the Kore source repo (makes sense rather than publishing each seperatly)
- Provide base capability for deploying and monitoring other Kore "**cluster apps**"

The components below are embedded Kore components (manifests):
1. [Kubernetes Application Controller](https://github.com/kubernetes-sigs/application#kubernetes-applications) -  monitoring componant Readyness and Health
1. [Captain](https://github.com/alauda/captain#captain) (a Helm 3 operator) - provides Helm artefact deployments using [CRD's](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
1. [Appvia Helm Repo](https://github.com/appvia/kore-helm-repo) provides a source of Chart artefacts tested with Kore clusters.
1. Kore auth proxy (**coming soon** - see [update kore auth proxy deployment](https://github.com/appvia/kore/issues/92))

Together these components will provide:
- The ability to deploy further Upstream Kore Componets.
- SSO using the Kore Auth Proxy

#### Mainatining Embedded Kore Components

To add or update an embedded Kore component:
1. Update the files here: `pkg/clusterman/manifests/...`. To ease of maintenance, we generate static code from these manifest files.
1. Ensure a Kubernetes Application resource is created, see [Kubernetes Application](https://github.com/kubernetes-sigs/application#kubernetes-applications)
1. Update the var (and array of manifest structures) called mm in `pkg/clusterman/clusterman.go`

These will be deployed automatically after any Kore kubernetes cluster is built.

#### Monitoring

To make sure what we deploy is reporting at least `Ready`, we deploy the [Kubernetes Application Controller](https://github.com/kubernetes-sigs/application#kubernetes-applications).

The Kubernetes Appication Controller also provides:
> - The ability to describe an applications metadata (e.g., that an application like WordPress is running)
> - A point to connect the infrastructure, such as Deployments, to as a root object. This is useful for tying things together and even cleanup (i.e., garbage collection)
> - Information for supporting applications to help them query and understand the objects supporting an application
Application level health checks.

At this time we build and deploy it to all Kore clusters ahead of this becoming part of Kubernetes.

We use this capability for simple `Readiness` monitoring of a collection of underlying resources. We aim to also integrate further health checks either to `Kubernetes Applications` or in another way.

### Upstream Kore Components

These are provided as a seperate deployment Artefacts (not compiled as part of Kore directly)

As the defacto artefact format for Kubernetes manifests, we have decided to use [Helm Charts](https://helm.sh/)

We plan to automatically reconcile these componets when issue [Reconcile Kubernetes Deployment Artefacts](https://github.com/appvia/kore/issues/87) is addressed.

#### Adding an Upstream Component

A chart compiled into the [kore-helm-repo](https://github.com/appvia/kore-helm-repo) with:
- An upstream published chart at specified version
- Any Kore customisations like static parameters or manifest updates
- A mandated Kubernetes Application resource (to ensure monitoring)

For any Clusterapp, the following simple steps add the capability:

A chart for any required operator is added to the Helm Chart Repo
A cluster app is defined in the clusterman with:
- A CRD is added to the embedded manifests, templated with the correct parameters from the params configmap
- An application resource is included with manifests when in cluster resources are created by the relevant CRD operator
- A jsonpath? or another way to monitor arbitrary CRD status when required
