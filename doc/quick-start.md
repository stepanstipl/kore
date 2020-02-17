## **Quick Start**

### Contents
- [Supported Cloud Providers](#supported-cloud-providers)
- [What is required?](#what-is-required)
- [Google Cloud account](#google-cloud-account)
- [Identity Broker](#identity-broker)
- [Configuring Auth0](#configuring-auth0)
- [Configuring test users](#configuring-test-users)
- [Running Kore](#running-kore)
- [Provisioning a team cluster](#provisioning-a-team-cluster)

The following is a quick start guide for running Kore locally to provision clusters on cloud platforms.

However, you'll still need access to an online identity provider to manage user authentication. See [Identity Broker](#identity-broker).

Please ensure you have the following installed on your machine,

- Docker: install instructions can be found [here]([https://docs.docker.com/install/](https://docs.docker.com/install/))
- Docker Compose: installation instructions can found [here](https://docs.docker.com/compose/install/)

### Supported Cloud Providers

Kore enables teams to provision clusters. Supported cloud providers include:

+ Google Cloud Provider (GCP)
+ Azure - `Coming Soon`
+ AWS - `Coming Soon`

### What is required?

- A GCP account with at least one Project and Service Account.
- An external facing identity provider that supports OpenID (Keycloak, Auth0, ForgeRock etc).

### Google Cloud account

We assume you're already setup as a Google Cloud user.

If not, grab a credit card and go to https://cloud.google.com/. Then, click the “Get started for free” button. Finally, choose whether you want a business account or an individual one.

#### Single cluster, multiple environments

For the purpose of this quick start, we're going to create a single cluster.

This cluster will use [Kubernetes Namespaces](https://kubernetes.io/docs/tasks/administer-cluster/namespaces/) to enable different environements for development, testing and production.

Next step: On GCP, select or create your target project.

#### Enabling the GKE API

(You can skip this step if GKE API is already enabled for this project)

With the a GCP Project selected or created,

- Head to the [Google Developer Console](https://console.developers.google.com/apis/api/container.googleapis.com/overview).
- Enable the GKE API.

#### Create a Service Account

(You can skip this step if you already have a Service Account setup)

With the a GCP Project selected or created,

- Head to the [IAM Console](https://console.cloud.google.com/iam-admin/serviceaccounts).
- Click `Create service account`.
- Fill in the form with details with your team's service account.

#### Configure your Service Account permissions

(You can skip this step if you're Service Account has the `Kubernetes Engine Admin` role)

- Assign the `Kubernetes Engine Admin` role to your Service account.

#### Create a key and download it (as JSON)

(You can skip this step if you already have your Service Account key downloaded in JSON format)

Kore will use this key to access the Service Account.

This is the last step, create a key and download it in JSON format.

### Identity Broker

Kore ships with the [`Dex` identity provider](https://github.com/dexidp/dex) or it can use an external identity provider for user management.

For this quick start guide we're using Auth0. See below for how to set it up and provision it.

#### Configuring Auth0

Auth0, found [here](https://auth0.com/), provides an enterprise SAAS identity provider.

- Sign up for an account from the [home page](https://auth0.com)
- From the dashboard side menu choose 'Applications' and then 'Create Application'
- Given the application a name and choose 'Regular Web Applications'
- Once provisioned click on the 'Settings' tab and scroll down to 'Allowed Callback URLs'. These are the permitted redirects for the applications. Since we are running the application locally off the laptop add `http://localhost:3000/callback` and `http://127.0.0.1:10080/oauth/callback` (Note the comma separation in the Auth0 UI.
- Scroll to the bottom of the settings and click the 'Show Advanced Settings'
- Choose the 'OAuth' tab from the advanced settings and ensure that the 'JsonWebToken Signature Algorithm' is set to RS256 and 'OIDC Conformant' is toggled on.
- Select the 'Endpoints' tab and note down the 'OpenID Configuration'.
- You can then scroll back to the top and note down the 'ClientID' and 'Client Secret'

Once you have the three pieces of the information *(ClientID, Client Secret and the OpenID endpoint)* you can substitute these settings on the [demo.yml](https://github.com/appvia/kore/blob/master/hack/compose/demo.yml); mapping to to ClientID, Client Secret and Discovery URL.

#### Configuring test users

Return to the Auth0 dashboard. From the side menu select 'Users & Roles' setting.

- Create a user by selecting 'Users'.
- Create a role by selecting 'Roles'.
- Add the role to the user.

### Running Kore

Once you have the above configured update the `demo.yml`:

```shell
KORE_CLIENT_ID: <YOUR_CLIENT_ID>
KORE_CLIENT_SECRET: <YOUR_CLIENT_SECRET>
KORE_DISCOVERY_URL: <OPENID_ENDPOINT>
```

To launch the Kore server, from the root directory, run

```shell
make demo
```

### Provisioning a team cluster

It's time to setup a team and the provision a dedicated cluster on GCP.

This video illustrates how to do using the Kore's CLI (korectl).

![Demo Video](images/demo.gif)