# Local Quick Start Guide

In this guide, we'll walk you through how to use the Appvia Kore CLI to set up a sandbox team environment locally and deploy a sample application.

We'll showcase how Appvia Kore can give you a head start with setting up [clusters](https://www.redhat.com/en/topics/containers/what-is-a-kubernetes-cluster), team members and environments.

**NOTE** The installation of Kore created by the `kore local` command in this quick start is suitable for testing and proof-of-concept work only. Bootstrapping a production installation of Kore is coming soon with [issue/340](https://github.com/appvia/kore/issues/340).

## Kubernetes

You'll need a Kubernetes provider to work through this guide. We simplify this by helping you set up a project on [GKE](https://cloud.google.com/kubernetes-engine).

**Please Note**: Created GKE clusters are for demo purposes only. They're tied to a local environment and will be orphaned once the local Kore instance is stopped.

## Team Access

Appvia Kore uses an external identity provider to manage team member identity and authenticate members.

For this guide, we'll help you to get set up on Auth0 to configure team access.

## Getting Started

- [Docker](#docker)
- [Google Cloud account](#google-cloud-account)
- [Configure Team Access](#configure-team-access)
- [Start Kore Locally with CLI](#start-kore-locally-with-cli)
- [Login as Admin with CLI](#login-as-admin-with-cli)
- [Create a Team with CLI](#create-a-team-with-cli)
- [Enable Kore to Set up Team Environments on GKE](enable-kore-to-set-up-team-environments-on-gke)
- [Provision a Sandbox Env with CLI](#provision-a-sandbox-env-with-cli)
- [Deploy An App to the Sandbox](#deploy-an-app-to-the-sandbox)
- [Cleaning Up](#cleaning-up)

### Docker

Please ensure you have the following installed on your machine,

- Docker: installation instructions can be found [here]([https://docs.docker.com/install/](https://docs.docker.com/install/))
- Docker Compose: installation instructions can found [here](https://docs.docker.com/compose/install/)

### Google Cloud account

If you don't have a Google Cloud account, grab a credit card and go to https://cloud.google.com/. Then, click the “Get started for free” button. Finally, choose whether you want a business account or an individual one.

Next step: On GCP, select an existing project or create a new one.

#### Enable the GKE API

(You can skip this step if GKE API is already enabled for this project)

With a GCP Project selected or created,

- Head to the [Google Developer Console](https://console.developers.google.com/apis/api/container.googleapis.com/overview).
- Enable the 'Kubernetes Engine API'.
- Enable the 'Cloud Resource Manager API'
- Enable the 'Compute Engine API'
- Enable the 'IAM Service Account Credentials API'

Alternatively you can enable these from the [gcloud](https://cloud.google.com/sdk/gcloud) command line;

```shell
# Setup if required
gcloud auth login (assuming you've not authenticated)
gcloud config set project <project_id>

# Enable the APIs
gcloud services enable cloudresourcemanager.googleapis.com
gcloud services enable iam.googleapis.com
gcloud services enable compute.googleapis.com
gcloud services enable container.googleapis.com
```

#### Create a Service Account

(You can skip this step if you already have a Service Account setup)

With the a GCP Project selected or created,

- Head to the [IAM Console](https://console.cloud.google.com/iam-admin/serviceaccounts).
- Click `Create service account`.
- Fill in the form with details with your team's service account.

#### Configure your Service Account permissions

(You can skip this step if you're Service Account has the `Owner` role)

- Assign the `Owner` role to your Service account.

#### Create a key and download it (as JSON)

(You can skip this step if you already have your Service Account key downloaded in JSON format)

Appvia Kore will use this key to access the Service Account.

This is the last step, create a key and download it in JSON format.

### Configure Team Access

Using Appvia Kore, team IAM (Identity and Access management) [is greatly simplified](security-gke.md#rbac).

Kore uses an external identity provider, like Auth0 or an enterprise's existing SSO system, to directly manage team member access to the team's provisioned environment.

For this guide, we'll be using Auth0 to configure team access.

#### Configure Auth0

[Auth0](https://auth0.com/), provides an enterprise SAAS identity provider.

Sign up for an account from the [home page](https://auth0.com).

From the dashboard side menu choose `Applications` and then `Create Application`

Give the application a name and choose `Regular Web Applications`

Once provisioned click on the `Settings` tab and scroll down to `Allowed Callback URLs`.
These are the permitted redirects for the applications. Since we are running the application locally off the laptop set
```
http://localhost:10080/oauth/callback,http://localhost:3000/auth/callback
```

Please make a note of the [__*Domain, Client ID, and Client Secret*__].

Scroll to the bottom of the settings and click the `Show Advanced Settings`

Choose the `OAuth` tab from the advanced settings and ensure that the `JsonWebToken Signature Algorithm` is set to RS256 and `OIDC Conformant` is toggled on.

#### Configuring test users

Return to the Auth0 dashboard. From the side menu select 'Users & Roles' setting.

- Create a user by selecting 'Users'.
- Create a role by selecting 'Roles'.
- Add the role to the user.

### Start Kore Locally with CLI

We'll be using our CLI, `kore`, to help us set up Kore locally.

#### Install the kore CLI

Find the latest kore release from https://github.com/appvia/kore/releases for your machine architecture and download it to a suitable location.

For example:

```shell script
KORE_VERSION=v0.1.0
curl -L https://github.com/appvia/kore/releases/download/${KORE_VERSION}/kore-cli-darwin-amd64 --output kore
chmod +x kore

# Confirm you have a working CLI:
./kore version
# kore version v0.1.0 (git+sha: aaaaaaa, built: 01-01-2020)
```

#### Configure Appvia Kore

You'll need access to the following details created earlier:

- Auth0 ClientID.
- Auth0 Client Secret.
- Auth0 domain.

Make sure you fill in the OpenID endpoint as `https://[Auth0 domain]/`, including the trailing `/`.

Once you have everything, run,

```shell script
./kore local configure
# What are your Identity Broker details?
# ✗ Client ID :
# ...
```

When configured correctly, you should see

```shell script
# ✅ A 'local' profile has been configured in ~/.kore/config
```

#### Start locally

```shell script
./kore local start
# ...Starting Kore.
# ...Kore is now started locally and is ready on http://127.0.0.1:10080
```

- Stop: To stop, run `./kore local stop`

- Logs: To view local logs, run `./kore local logs`

At this point, you can use the CLI as detailed below, and also browse to the Kore UI at http://127.0.0.1:3000/

### Login as Admin with CLI

You now have to login to be able to create teams and provision environments.

This will use our Auth0 set up for IDP. As you're the only user, you'll be assigned Admin privileges.

```shell script
./kore login
# Attempting to authenticate to Appvia Kore: http://127.0.0.1:10080 [local]
# Successfully authenticated
```

### Create a Team with CLI

Let's create a team with the CLI. In local mode, you'll be assigned as team member to this team.

As a team member, you'll be able to provision environments on behalf of team.

```shell script
./kore create team --description 'The Appvia product team, working on project Q.' team-appvia
# "team-appvia" team was successfully created
```

To ensure the team was created,

```shell script
./kore get teams team-appvia
# Name            Description
# team-appvia     The Appvia product team, working on project Q.
```

### Enable Kore to Set up Team Environments on GKE

We now need to give Kore the credentials it needs to build a cluster on our behalf. This command imports a set of credentials into kore 
and allows your new team to use them to make clusters. 

We'll then use these to create a cluster to host our sandbox environment. You'll need the following details which you set up earlier:

- GKE Project ID.
- Path to the service account key JSON file.

```shell script
./kore create gkecredentials gke --description "GKE Credentials" -p <gcp-project-id> --cred-file <path-to-json-service-account> --allocate team-appvia
# Storing credentials in Kore
# Waiting for resource "gke" to provision (you can background with ctrl-c)
# Successfully provisioned the resource: "gke"
# Storing credential allocation in Kore
# Waiting for resource "gke" to provision (you can background with ctrl-c)
# Successfully provisioned the resource: "gke"
```

### Provision a Sandbox Env with CLI

Its time to use the Kore CLI To provision our Sandbox environment,

```shell script
./kore create cluster appvia-trial -t team-appvia --plan gke-development -a gke --namespaces sandbox
# Attempting to create cluster: "appvia-trial", plan: gke-development
# Waiting for "appvia-trial" to provision (usually takes around 5 minutes, ctrl-c to background)
# Cluster appvia-sdbox has been successfully provisioned
# --> Attempting to create namespace: sandbox

# You can update your kubeconfig via: $ kore kubeconfig -t team-appvia
# Then use 'kubectl' to interact with your team's cluster
```

There's a lot to unpack here. So, lets walk through it,

- `create cluster`, we create a [cluster](https://www.redhat.com/en/topics/containers/what-is-a-kubernetes-cluster) to host our sandbox environment.

- `appvia-trial`, the name of the cluster.

- `-t team-appvia`, the team for which we are creating the sandbox environment.

- `--plan gke-development`, a Kore predefined plan called `gke-development`. This creates a cluster ideal for non-prod use.

- `-a gke`, the `gke` allocated credential to use for creating this cluster.

- `--namespace sandbox`, creates an environment called `sandbox` in the `appvia-trial` where we can deploy our apps, servers, etc..

You now have a sandbox environment locally provisioned for your team. 🎉

### Deploy An App to the Sandbox

We'll be using `kubectl`, the Kubernetes CLI, to make the deployment. If you don't have it already, [please install and setup kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl-on-macos).

Now we have to configure our `kubectl` kubeconfig in ~/.kube/config with our new GKE cluster.

```shell script
./kore kubeconfig -t team-appvia
# Successfully added team [team-appvia] provisioned clusters to your kubeconfig
# Context        Cluster
# appvia-trial   appvia-trial
```

Switch the current `kubectl` context to `appvia-trial`,

```shell script
kubectl config use-context appvia-trial --namespace=sandbox
# + kubectl config use-context appvia-trial --namespace=sandbox
# Switched to context "appvia-trial".
```

Deploy the GKE example web application container available from the Google Cloud Repository

```shell script
kubectl create deployment hello-server --image=gcr.io/google-samples/hello-app:1.0
# + kubectl create deployment hello-server --image=gcr.io/google-samples/hello-app:1.0
# deployment.apps/hello-server created

kubectl expose deployment hello-server --type LoadBalancer --port 80 --target-port 8080
# + kubectl expose deployment hello-server --type LoadBalancer --port 80 --target-port 8080
# service/hello-server exposed
```

Get the `EXTERNAL-IP` for `hello-server` service

```shell script
kubectl get service hello-server
# + kubectl get services
# NAME           TYPE           CLUSTER-IP     EXTERNAL-IP          PORT(S)        AGE
# hello-server   LoadBalancer   10.70.10.119   <35.242.154.199>     80:31319/TCP   23s
```

Now navigate to the `EXTERNAL-IP` as a url

```shell script
open http://35.242.154.199
```

You should see this on the webpage

```text
Hello, world!
Version: 1.0.0
Hostname: hello-server-7f8fd4d44b-hpxls
```

### Cleaning Up

To avoid incurring charges to your Google Cloud account for the resources used in this quickstart, follow these steps.

#### Delete the app from the sandbox environment

```shell script
kubectl delete service hello-server
```

#### Delete the sandbox environment

You can now use kore to destroy the cluster:

```shell script
./kore delete --team team-appvia cluster appvia-trial
# "appvia-trial" was successfully deleted
```

You can check for the cluster deletion completing by retrieving the cluster:

```shell script
./kore get cluster appvia-trial --team team-appvia
# Name            Kind    API Endpoint           Auth Proxy Endpoint    Status
# appvia-trial    GKE     https://1.2.3.4        5.6.7.8                Deleting
```

Once the deletion is complete, the cluster will disappear from Kore:

```shell script
./kore get cluster appvia-trial --team team-appvia
# Error: "appvia-trial" does not exist
```

Finally, after waiting for your cluster to delete, you may stop your local kore environment:

```shell script
./kore local stop
# ...Stopping Kore.
# ...Kore is now stopped.
```
