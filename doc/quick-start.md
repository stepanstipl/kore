## **Getting Started**

### Contents
- [Supported Cloud Providers](#supported-cloud-providers)
- [What is required?](#what-is-required)
- [Configuring your cloud provider](#configuring-your-cloud-provider)
- [Identity Broker](#identity-broker)
- [Configuring Auth0](#configuring-auth0)
- [Configuring test users](#configuring-test-users)
- [Running Kore](#running-kore)
- [Provisioning Credentials](#provisioning-credentials)
- [Provisioning a team cluster](#provisioning-a-team-cluster)

The following is a quick start guide for running Kore locally to provision clusters on cloud platforms.

However, you'll still need access to an online identity provider to manage cluster access endpoints. See [Identity Broker](#identity-broker).

Please ensure you have the following installed on your machine,

- Docker: install instructions can be found [here]([https://docs.docker.com/install/](https://docs.docker.com/install/))
- Docker Compose: installation instructions can found [here](https://docs.docker.com/compose/install/)

### Supported Cloud Providers

Kore enables teams to provision clusters. Supported cloud providers include:

+ Google Cloud Provider (GCP)
+ Azure - `Coming Soon`
+ AWS - `Coming Soon`

### What is required?

First, you require permissions to create a service account [on GCP](https://cloud.google.com/iam/docs/service-accounts).

- Permissions to setup or acquire cloud credentials neccessary to provision infrastructure in the GCP and or AWS.
- An external facing identity provider; anything that supports OpenID (Keycloak, Auth0, ForgeRock etc)

### Configuring your cloud provider

There is automated account provisioning for AWS and Google, where an isolated user account can be created that maps to a specific team. The Account or Project account provisioning uses least-privilege and will create a project or AWS Account service account, that gives it enough permissions to create other accounts or projects.

From that point on, it will create another service account inside the child account or project, for just managing Kubernetes and the related resources, (GKE or EKS). It is this account that is then used by Kore to provision the Kubernetes services, of which, options are controlled by the plans defined by the administrators.

### Identity Broker

Kore is designed to use an external identity provider for user management. You can bring your own IDP, but for this quick start guide we're using Auth0.

See below for how to setup and provision Auth0 for testing.

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

Once you have the above configured you

```shell

# change the following the hack/compose/demo.yml
KORE_CLIENT_ID: <YOUR_CLIENT_ID>
KORE_CLIENT_SECRET: <YOUR_CLIENT_SECRET>
DISCOVERY_URL: <OPENID_ENDPOINT>
```

To launch the Kore server, from the root directory, run

```shell
make demo
```

Open up a web browser and point it at http://localhost:3000; the password defaults for 'password' but can be changed using the KORE_ADMIN_PASS environment variable.

You can now configure specific cloud providers. For GCP, this will require a Project Id, that includes created a service account that has GKE privilages.

Once this is configured, all teams can use these downstream to provision clusters using the defined plans.

### Provisioning Credentials

The self-service nature in Kore is provided via the use of allocations (`kubectl get crd allocations.config.kore.appvia.io` api group). An adminstrator creates shared credentials which are then allocated out to one or more teams to self-serve; at present the alpha release reuses the credentials across the teams, though we are currently integrating the provisioning of cloud account management into the product.

### Provisioning a team cluster

It's time to setup a team and the provision a dedicated cluster on GCP.

This video illustrates how to do using the Kore's CLI (korectl).

![Demo Video](images/demo.gif)