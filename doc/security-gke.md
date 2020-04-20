# GKE Security With Kore

At the end of the cluster create, we should be able to say that we:

+ [RBAC](#rbac) - Enable RBAC inside the cluster
+ [Audit](#audit) - Enable audit logging of the OS: SSH, kubectl exec's etc.
+ [PSP](#psp-pod-security-policy) - Enable pod security policies
+ [Disabled Legacy Authentication](#disabled-legacy-authentication) - Disable legacy authentication methods
+ [Shielded Node Image](#shielded-nodes) for authenticity validation and rootkit avoidance (beta)
+ [Hardened Node Image](#hardened-image-node) - this is default
+ [Minimal Privilege - Node Accounts / Workload Identities](#minimal-privilege-node-accounts) - (this should be default by GCP)
+ [Network Policies Between Pods](#network-policies-between-pods) these should be enabled
+ Secret encryption using cloud KMS
+ Private nodes and control plane access (with IP address / firewalled)
+ Intranode visibility (for flow logs)
+ Auto node upgrade (to keep up to date)

## RBAC

GKE enables Role Bases Access Control (RBAC) as default in GKE, see Google documentation; [role-based-access-control](https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control):

> Kubernetes RBAC is enabled by default.

### Without Kore

1. &#x2611; RBAC is integrated by Google with gcloud IAM users.
1. &#x2612; Every user of GKE would need a Goolge project access.
1. &#x2612; Google IAM would need to be managed - see [Predefined GKE roles](https://cloud.google.com/kubernetes-engine/docs/how-to/iam#predefined) to manage authorization.

### With Kore

1. &#x2611; RBAC authorization is integrated with Kore Users and Teams and the configured [IDP](./idp.md).
1. &#x2611; No Google project access required with Kore.
1. &#x2611; Access controls are provided natively by Kubernetes subjects we manage through Teams.

## Audit

With GKE clusters audit Logs are visible through [stackdriver](https://cloud.google.com/stackdriver).

Different levels of audit are avialble with GKE clusters and Google Cloud Projects GCP access.

Googles own actions on your GKE clusters and projects are configured separately - [Access Transparency](https://cloud.google.com/logging/docs/audit/access-transparency-overview).

### Without Kore

1. &#x2611; [Admin Activity Logging](https://cloud.google.com/kubernetes-engine/docs/how-to/audit-logging#audit_logs_in_your_project).
1. &#x2612; [Data Access Logs](https://cloud.google.com/logging/docs/audit/configure-data-access) are not automatically enabled.
1. &#x2612; Node Access metrics - [Manually by deploying a logging daemonset](https://cloud.google.com/kubernetes-engine/docs/how-to/linux-auditd-logging#deploying_the_logging_daemonset).


### With Kore

1. &#x2611; Admin Activity Logging](https://cloud.google.com/kubernetes-engine/docs/how-to/audit-logging#audit_logs_in_your_project).
1. &#x2611; [Data Access Logs](https://cloud.google.com/logging/docs/audit/configure-data-access) are enabled in Kore's [default plans](https://github.com/appvia/kore/blob/a6a1c3e38bdd2b83150cc8eb5434d7303b6498b3/pkg/kore/assets/plans.go#L67-L68). [*](#note-kore-examples)
1. &#x2612; Node Access metrics - [Manually by deploying a logging daemonset](https://cloud.google.com/kubernetes-engine/docs/how-to/linux-auditd-logging#deploying_the_logging_daemonset). See issue [kore/issues/150](https://github.com/appvia/kore/issues/150).

#### Note Kore Examples

Our [Examples](./examples/) do not enable StackDriver audit or logging (as they are designed for local test only and cost effectiveness).

## PSP - Pod Security Policies

### Without Kore

- &#x2612; PSP Controller - [Manually enable PSP controller](https://cloud.google.com/kubernetes-engine/docs/how-to/pod-security-policies#enabling_podsecuritypolicy_controller) - This is a manual step.
- &#x2612; PSP Policies - [Manually define PSP policies](https://cloud.google.com/kubernetes-engine/docs/how-to/pod-security-policies#define_policies)

### With Kore

- &#x2611; PSP Controller - [Manually enable PSP controller](https://cloud.google.com/kubernetes-engine/docs/how-to/pod-security-policies#enabling_podsecuritypolicy_controller) - This is automatically configured.
- &#x2611; PSP Policies - [Manually define PSP policies](https://cloud.google.com/kubernetes-engine/docs/how-to/pod-security-policies#define_policies) - This is automatically configured.

## Disabled Legacy Authentication

The following legacy authentication methods should be disabled as they represent a security risk (long lived credentials).
- Client Certificate Issuer
    These should be disabled as they represent long lived credentials that are hard to rotate and as such they are being discontinued by [GKE from 1.12+](https://cloud.google.com/kubernetes-engine/docs/how-to/hardening-your-cluster#restrict_authn_methods)
- Basic Authentication

### Without Kore

- &#x2612; Client Certificate Issuer is disabled from 1.12+
- &#x2612; Basic Authentication is disabled from 1.12+

### With Kore

- &#x2611; Client Certificate Issuer is disabled by default
- &#x2611; Basic Authentication is disabled by default

#### Kore Admin User

The Kore API (not the kubernetes API) has a Basic Authentication capability for the admin user - see [kore/issues/151](https://github.com/appvia/kore/issues/151).

## Shielded nodes

This provides GKE node authenticity validation and rootkit avoidance - see [Shielded GKE nodes](https://cloud.google.com/kubernetes-engine/docs/how-to/shielded-gke-nodes).

### Without Kore

- &#x2611; Not enabled by default

### With Kore

- &#x2611; Not enabled - see [kore/issues/152](https://github.com/appvia/kore/issues/152)

## Hardened Node Image

Google recommend a [Hardened Node Image](https://cloud.google.com/kubernetes-engine/docs/how-to/hardening-your-cluster#containerd) which is enabled by default.

### Without Kore

- &#x2611; Enabled by default and optional

### With Kore

- &#x2611; Enabled by default and currently optional - see [kore/issues/153](https://github.com/appvia/kore/issues/153)

## Minimal Privilege Node Accounts

The "default" service account used by nodes should now have least privelage as documented here:
- [GKE - Changes in access scopes](https://cloud.google.com/kubernetes-engine/docs/how-to/access-scopes#changes)
- [GKE - 245.0.0 (2019-05-07) - Breaking Changes](https://cloud.google.com/sdk/docs/release-notes#breaking_changes_23)

This is further documented in [appvia/kore/issues/161](https://github.com/appvia/kore/issues/161).

Google now recommend using [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity) but this has some limitations most notably:

> Workload Identity can't be used with Pods running in the host network.

### Without Kore

- &#x2611; **Least Privilege Node Accounts** yes, see above
- &#x2612; **Workload Identities**phe§ are not enabled by default

### With Kore

- &#x2611; **Least Privilege Node Accounts** yes, see above
- &#x2612; **Workload Identities** not enabled by default - see [kore/issues/165](https://github.com/appvia/kore/issues/165).

## Network Policies Between Pods

### By default

You need to apply Network Policies in order to enforce them.

### Without Kore

- &#x2611; Network Policy is Enabled

### With Kore

- &#x2611; Network Policy is Enabled

There is a single issue related to Calico configuration with Kore clusters when pods are created directly: [kore/issues/179](https://github.com/appvia/kore/issues/179).
