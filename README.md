
<h1><img src="https://www.appvia.io/hubfs/APPVIA_November2019/Images/appvia_logo.svg" alt="Appvia Kore" width="200"></h1><em>

### **[Getting Started](#getting-started)** • **[Architecture](#architecture)** • **[Contributing](#contributing)** • **[Roadmap](#roadmap)**

![Latest Release](https://img.shields.io/github/v/release/appvia/kore?style=for-the-badge&label=Latest%20Release&color=%23D1374D)
![License: Apache-2.0](https://img.shields.io/github/license/appvia/kore?style=for-the-badge&color=%23D1374D)

## Kubernetes for Teams via Appvia Kore
- **Cluster provisioning** provides secure and consistent provisioning of environments for teams.
- **Accounts & Account Users** provides a single source to access and control across the estate
- **Single Sign On**
- **Account Limits** to ensure quality of service and fairness when sharing a cluster
- **Namespace Templates** for secure tenant isolation and self-service namespace initialization
- **Multi-Cluster And Single Tenant Management** for sharing a pool of clusters ([coming soon](#roadmap))

<br>

![kore Demo Video](docs/website/images/demo.gif)

<br>

## Contents
- [Why Appvia Kore?](#why-appvia-kore)
  - [Developer](#the-developer)
  - [Devops](#the-devops)
- [Workflow & Interactions](#workflow--interactions)
- [Architecture](#architecture)
  - [Workflow & Interactions](#workflow--interactions)
  - [Custom Resources & Resource Groups](#custom-resources--resource-groups)
- [Getting Started](#getting-started)
  - [0. Requirements](#0-requirements)
  - [1. Install Kore](#1-install-kiosk)
  - [2. Configure Accounts](#2-configure-accounts)
  - [3. Working with Spaces](#3-working-with-spaces)
  - [4. Setting Account limits](#4-setting-account-limits)
  - [5. Working with Templates](#5-working-with-templates)
- [Uninstall Kore](#uninstall-kiosk)
- [Roadmap](#roadmap)
- [Contributing](#contributing)

## Why Appvia Kore?

Appvia Kore is designed to make Kubernetes a commodity for your organisation and teams. Allowing any team to be able to get Kubernetes simply and easily, without relying on specialist resources to provision it for you, set it up and give you access credentials.

## The Developer

>-   **Self serve Kubernetes clusters with best practice**
>-   **On board and off board your teammates**
>-   **Manage your Role Based Access Controls**

You know how to work with containers or you’re starting that journey, but learning Kubernetes and Cloud is complicated and not necessarily a good use of your time. You just want to start deploying applications and using in-cluster or cloud services to start showing off your great work!

With Kore an administrator can configure a set of default known cluster best practices so you and your teammates can provision clusters safely and securely. You can on-board and off-board your own users, with them automatically getting access so you can iterate and deploy applications easily.

You can see what versions are deployed where, generate robot tokens for your Continuous Integration system and start iterating quickly and securely! With out-of-the-box developer roles and policies, you know that your applications are meeting the security requirements they need to!

## The DevOps

>-   **Define best practice clusters as plans for developer teams to self-serve**
>-   **Push user and application policies to clusters to enforce organisation wide controls**
>-   **Iterate cluster validation as Code in your CI pipelines**
>-   **Improve Kubernetes cluster security with security policies**

You know what your developers need, they are also not shy of telling you! But you want to make sure that they get reliable, scalable and secure services before they consume them. Configuring all of this takes time and some of the components can be very complicated and time consuming to do correctly.

Not only that, but there is the cloud provisioning and architecture, that adds another hurdle to get over before teams can even consume Kubernetes and other cloud services! With Kore, we decided to enable cloud setup, as cloud accounts or projects are free, you can isolate teams to accounts, making cloud cost management, security and access controls simpler!

You can also setup default Kubernetes cluster plans that developer teams can consume without you needing to run CI, code or scripts. You can validate the plans work as they should once and then make them available. This might be anything from a Developer plan that is a single availability zone with cheap instance types, through to production plans with enhanced security settings across multiple availability zones for resilience and on-demand instances.

You can also define sensible policies globally, enhancing the security footprint once across all of your teams, so you know that everyone is working consistently and safely!

## Architecture

Appvia Kore reuses the Kubernetes framework and enhances it to provide a more enriched set of features as well as an improved and simplified developer and operations experience. These enhancements that we have created are:

> -  **Kubernetes Cluster Plans**
>-   **Team management and creation**
>-   **SSO and authentication with your organisation IDP**
> -  **Auditability on user actions, cluster creation and access management**

Each enhancement works under the operator framework. The operators are domain specific features, such as team management or SSO configuration. To bring each domain specific operator together, we have the Kore API, which bridges each service and manages the coordination of data into each operator on your behalf.

All of the components run as a set of containers, so as long as there is Docker, you can run this either locally or in cloud. As it is using the Kubernetes framework, Kubernetes is a prerequisite to host Appvia Kore.

Note Appvia Kore is deemed an early release, the project is not regarded as production ready and is under rapid development; thus expect new features to rollout.

## Getting Started

The following provides a quick start guide for rolling out and playing with the product locally; please ensure you have the following installed on your machine

- Docker: install instructions can be found [here]([https://docs.docker.com/install/](https://docs.docker.com/install/))</em>
- Docker Compose: installation instructions can found [here](https://docs.docker.com/compose/install/)

