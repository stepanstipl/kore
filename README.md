# Kore
Kubernetes for Teams

[![Build_Status](https://circleci.com/gh/appvia/kore.svg?style=svg)](https://circleci.com/gh/appvia/kore)

---
We make it easy for teams to get Kubernetes across multiple clouds. By using Kore an administrator can setup Kubernetes cluster best practices and allow teams to self-provision clusters, environmnets, login to their organisations SSO and deploy their applications.
---

## Overview

Kore is designed to make Kubernetes a commodity for your organisation, allowing any team to be able to get Kubernetes simply and easily, without relying on specialist resources to provision it for you, set it up and give you access credentials.

Features:
- Support for Github, Openid and SAML authentication
- Create Kubernetes cluster plans with best practice that define teams cluster makeup, (Production plans, Developer plans or Machine learning plans etc.)
- Let Kore provision a google project per team for better cost visibliity and security
- Let Kore provision an AWS Account per team, for better cost visiblity and security
- Support for AWS, Google with support for Azure and VMware coming soon
- Create your own team
- Invite users your team
- Allow anyone in your team to create Kubernetes clusters
- Provision namespaces, (environments) and sign in using single sign on and start deploying applications


## Getting Started

### Administrator Setup

*Note* This is just to get up and running locally, not suitable to run as production like this!

First checkout the repository:

`git checkout git@github.com/appvia/kore`

And run make demo

```
cd kore
make demo
```

*Note*: You will not be able to login to the clusters, as it will require an external IDP, (for that you will need to host kore with a consumable endpoint for authentication)/


Steps:
1. Setup cloud
2. Setup SSO
3. Setup plans

Once complete developer teams can login and create teams and clusters in your cloud provider.


There are two ways to setup Kore:

1. [With the CLI]()
2. [Using the web frontend]()


#### With the CLI


`git checkout git@github.com/appvia/kore` 

#### Using the web frontend


### Developer User


[Who is it For]

[Architecture]

[Cloud Support]

[Contributing]

[User Management]

[Policies]

[Global Configuration]

[License]


[Overview]: #Overview
[Get Started]: doc/getting-started.md
[Who is it For]: doc/users.md
[Architecture]: doc/architecture.md
[Cloud Support]: doc/cloud.md
[Contributing]: contributing.md
[User Management]: doc/user-management.md
[Policies]: doc/policies.md
[Global Configuration]: doc/global-configuration.md
[License]: LICENSE
