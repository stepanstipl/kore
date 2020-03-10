# Local Quick Start Guide (Alpha )

In this guide, we'll walk you through how to use the Appvia Kore CLI to set up a sandbox team environment locally and deploy a sample application.

We'll showcase how Appvia Kore can give you a head start with setting up [clusters](https://www.redhat.com/en/topics/containers/what-is-a-kubernetes-cluster), team members and environments.

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
- [Start Kore Locally with CLI](#use-cli-to-start-kore-locally)
- [Login as Admin with CLI](#login-as-admin-with-cli)
- [Create a Team with CLI](#create-a-team-with-cli)
- [Enable Kore to Set up Team Environments on GKE](enable-kore-to-set-up-team-environments-on-gke)
- [Provision a Sandbox Env with CLI](#provision-a-sandbox-env-with-cli)
- [Deploy An App to the Sandbox](#deploy-an-app-to-the-sandbox)

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
- Enable the GKE API.

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

Scroll to the bottom of the settings and click the `Show Advanced Settings`

Choose the `OAuth` tab from the advanced settings and ensure that the `JsonWebToken Signature Algorithm` is set to RS256 and `OIDC Conformant` is toggled on.

Select the `Endpoints` tab and note down the `OpenID Configuration`.

Please make a note of the [__*ClientID, Client Secret and the OpenID endpoint*__].

#### Configuring test users

Return to the Auth0 dashboard. From the side menu select 'Users & Roles' setting.

- Create a user by selecting 'Users'.
- Create a role by selecting 'Roles'.
- Add the role to the user.

### Start Kore Locally with CLI

We'll be using our CLI, `korectl`, to help us set up Kore locally. 

#### Install the korectl CLI

For the time being, we have to clone this repo to install the CLI.

```shell script
git clone git@github.com:appvia/kore.git kore-test

cd kore-test

make korectl

bin/korectl -v
# korectl version v0.0.12 (git+sha: cef3143, built: 06-03-2020)
```

#### Configure Appvia Kore

You'll need access to the following details created earlier:

- Auth0 ClientID.
- Auth0 Client Secret.
- Auth0 OpenID endpoint.
- GKE Project ID.
- GKE Region.
- Path to the service account key JSON file.

Once you have everything, run,

```shell script
bin/korectl local configure
# What are your Identity Broker details?
# ✗ Client ID :
# ...
```

When configured correctly, you should see

```shell script
# ✅ A 'local' profile has been configured in ~/.korectl/config
# ✅ Generated Kubernetes CRDs are now stored in <project root>/manifests/local directory.
```

#### Start locally

```shell script
bin/korectl local start
# ...Starting Kore.
# ...Kore is now started locally and is ready on http://127.0.0.1:10080
```

- Stop: To stop, run `bin/korectl local stop`

- Logs: To view local logs, run `bin/korectl local logs`

### Login as Admin with CLI

You now have to login to be able to create teams and provision environments.

This will use our Auth0 set up for IDP. As you're the only user, you'll be assigned Admin privileges.  

```shell script
bin/korectl login
# Attempting to authenticate to Appvia Kore: http://127.0.0.1:10080 [local]
# Successfully authenticated
```

### Create a Team with CLI

Let's create a team with the CLI. In local mode, you'll be assigned as team member to this team.

As a team member, you'll be able to provision environments on behalf of team.

```shell script
bin/korectl create team --description 'The Appvia product team, working on project Q.' team-appvia
# "team-appvia" team was successfully created
```

To ensure the team was created,

```shell script
bin/korectl get teams team-appvia
#Name            Description
#team-appvia     The Appvia product team, working on project Q.
```

### Enable Kore to Set up Team Environments on GKE

This command applies a set of manifests created when configuring Kore to run locally.

When applied, these manifests give Kore the credentials necessary to build a GKE cluster on behalf of our team. 

This cluster will in turn host our sandbox environment.

```shell script
bin/korectl apply -f manifests/local/gke-credentials.yml -f manifests/local/gke-allocation.yml
# gke.compute.kore.appvia.io/teams/kore-admin/gkecredentials/gke configured
# config.kore.appvia.io/teams/kore-admin/allocations/gke configured
```

### Provision a Sandbox Env with CLI

Its time to use the Kore CLI To provision our Sandbox environment,

```shell script
bin/korectl create cluster appvia-trial -t team-appvia --plan gke-development -a gke --namespace sandbox
# Attempting to create cluster: "appvia-trial", plan: gke-development
# Waiting for "appvia-trial" to provision (usually takes around 5 minutes, ctrl-c to background)
# Cluster appvia-sdbox has been successfully provisioned
# --> Attempting to create namespace: sandbox

# You can retrieve your kubeconfig via: $ korectl clusters auth -t team-appvia
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

We'll be using `kubectl`, the Kubernetes CLI, to make the deployment.

First we have to set up our kubeconfig in `~/.kube/config` with our new GKE cluster.

```shell script
bin/korectl clusters auth -t team-appvia
# Successfully updated your kubeconfig with credentials
```

Switch the current `kubectl` context to `appvia-trial`,

```shell script
kubectl config set-context appvia-trial
# + kubectl config set-context appvia-trial --namespace=sandbox
# Context "appvia-trial" modified.
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