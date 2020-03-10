# Local Quick Start Guide (Alpha )

In this guide, we'll walk you through how to use the Appvia Kore CLI to set up a sandbox team environment locally.

We'll showcase how Appvia Kore can give you a head start with setting up clusters, team members and environments.

## Kubernetes

You'll need a Kubernetes instance to work through this guide. We simplify this by helping you set up a project on GKE.

**Please Note**: Created GKE clusters are for demo purposes only. They're tied to a local environment and will be orphaned once the local Kore instance is stopped.

## Team Access

Appvia Kore uses an external identity provider, like Auth0, to manage team member identity and authenticate members.

To keep things simple, we'll help you get set up on Auth0 to configure team access.

## Getting Started

- [Docker](#docker)
- [Google Cloud account](#google-cloud-account)
- [Configure Team Access](#configure-team-access)
- [Start Kore Locally with CLI](#use-cli-to-start-kore-locally)

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
- `http://localhost:3000/auth/callback` and 
- `http://localhost:10080/oauth/callback` (Note the comma separation in the Auth0 UI).

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
