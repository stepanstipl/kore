- [**Getting Started**](#getting-started)

## Contents
- [Contents](#contents)
  - [Supported Cloud Providers](#supported-cloud-providers)
  - [What is required?](#what-is-required)
  - [Identity Broker](#identity-broker)
  - [Configuring Auth0](#configuring-auth0)
  - [Running the Demo](#running-the-demo)
  - [Provisioning Credentials](#provisioning-credentials)

![Demo Video](doc/images/demo.gif)

The following provides a quick start guide for rolling out and playing with the product locally; please ensure you have the following installed on your machine

- Docker: install instructions can be found [here]([https://docs.docker.com/install/](https://docs.docker.com/install/))</em>
- Docker Compose: installation instructions can found [here](https://docs.docker.com/compose/install/)

While Kore are be run locally off a laptop for testing there are components which need to externally accessible; namely the identity provider due to the requirement for the clusters to access the endpoints.

### Supported Cloud Providers

The aim of Kore is to enable teams to provision clusters. The supported cloud providers are:

+ Google Cloud Provider (GCP)
+ Azure - `Coming Soon`
+ AWS - `Coming Soon`

There is automated account provisioning for AWS and Google, where an isolated user account can be created that maps to a specific team. The Account or Project account provisioning uses least-privilege and will create a project or AWS Account service account, that gives it enough permissions to create other accounts or projects. From that point on, it will create another service account inside the child account or project, for just managing Kubernetes and the related resources, (GKE or EKS). It is this account that is then used by Kore to provision the Kubernetes services, of which, options are controlled by the plans defined by the administrators.

As of alpha release the only cloud provider support is Google Kubernete Engine (GKE), though we will roll out support for Amazon’s EKS, Azure AKS and Cluster API in the near future.

### What is required?

Naturally given you are building in cloud you require the permissions to create a service account in the GCP or IAM permissions in the AWS.

- An external facing identity provider; anything that supports OpenID (Keycloak, Auth0, ForgeRock etc)
- Permissions to setup or acquire cloud credentials neccessary to provision infrastructure in the GCP and or AWS.

### Identity Broker

Appvia Kore is designed to use an external identity provider for user management. The product is feature flagged and bundled with Dex IDP but given the requirement for an IDP to be externally facing, unless your deploying the product in full you will need to use your current IDP or use a SAAS product for testing purposes. Below gives a quick how to setup an provision Auth0 for testing.

### Configuring Auth0

Auth0, found [here](https://auth0.com/), provides an enterprise SAAS identity provider

- Sign up for an account from the [home page](https://auth0.com)
- From the dashboard side menu choose 'Applications' and then 'Create Application'
- Given the application a name and choose 'Regular Web Applications'
- Once provisioned click on the 'Settings' tab and scroll down to 'Allowed Callback URLs'. These are the permitted redirects for the applications. Since we are running the application locally off the laptop are and add `http://localhost:3000/callback` and `http://127.0.0.1:10080/oauth/callback` (Note the comma separation in the Auth0 UI.
- Scroll to the bottom of the settings and click the 'Show Advanced Settings'
- Choose the 'OAuth' tab from the advanced settings and ensure that the 'JsonWebToken Signature Algorithm' is set to RS256 and 'OIDC Conformant' is toggled on.
- Select the 'Endpoints' tab and note down the 'OpenID Configuration'.
- You can then scroll back to the top and note down the 'ClientID' and 'Client Secret'

Once you have the three pieces of the information *(ClientID, Client Secret and the OpenID endpoint)* you can substitute these settings on the [demo.yml](https://github.com/appvia/kore/blob/master/hack/compose/demo.yml); mapping to to ClientID, Client Secret and Discovery URL.

The next logical step would be to return to the dashboard of Auth0 and create one or more test users under the 'Users & Roles' settig

### Running the Demo

Once you have the above configured you

```shell

# change the following the hack/compose/demo.yml
KORE_CLIENT_ID: <YOUR_CLIENT_ID>
KORE_CLIENT_SECRET: <YOUR_CLIENT_SECRET>
DISCOVERY_URL: <OPENID_ENDPOINT>
```

You can then run the `make demo` command from the root directory; which will bring up the dependencies within docker-compose. From here you can open up the browser and point it at http://localhost:3000; the password defaults for 'password' but can be found on the KORE_ADMIN_PASS environment variable. You can then configure specific cloud providers, for GCP, this will be a project credential for GKE that should be a service account credential with privileges for GKE.

Once this is configured, all teams can use these downstream to provision clusters using the defined plans.

### Provisioning Credentials

The self-service nature in Kore is provided via the use of allocations (`kubectl get crd allocations.config.kore.appvia.io` api group). Adminstrator create shared credentials which are then allocated out to one or more teams to self-serve; at present the alpha release reuses the credentials across the teams, though we are currently integrating the provisioning of cloud account management into the product.

