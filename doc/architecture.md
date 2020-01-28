## The Architecture

### Overview

Appvia Kore uses the Kubernetes framework and enhances it to provide a more enriched set of features as well as an improved and simplified developer and operations experience. These enhancements that we have created are:
+ Kubernetes Cluster Plans
+ Team management and creation
+ SSO and authentication with your organisation IDP
+ Auditability on user actions, cluster creation and access management

Each enhancement works under the operator framework. The operators are domain specific features, such as team management or SSO configuration. To bring each domain specific operator together, we have the Kore API, which bridges each service and manages the coordination of data into each operator on your behalf.

All of the components run as a set of containers, so as long as there is Docker, you can run this either locally or in cloud. As it is using the Kubernetes framework, Kubernetes is a prerequisite to host Appvia Kore.

Note Appvia Kore is deemed an early release, the project is not regarded as production ready and is under rapid development; thus expect new features to rollout.
