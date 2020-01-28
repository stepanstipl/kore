# Supported Cloud Providers

The aim of Kore is to enable teams to provision clusters. The supported cloud providers are:
+ Google Cloud Provider (GCP)
+ Azure
+ AWS

There is automated account provisioning for AWS and Google, where an isolated user account can be created that maps to a specific team. The Account or Project account provisioning uses least-privilege and will create a project or AWS Account service account, that gives it enough permissions to create other accounts or projects. From that point on, it will create another service account inside the child account or project, for just managing Kubernetes and the related resources, (GKE or EKS). It is this account that is then used bar Kore to provision the Kubernetes services, of which, options are controlled by the plans defined by the administrators.

As of alpha alpha release the only cloud provider support is Google Kubernete Engine (GKE), though we will roll out support for Amazon’s EKS, Azure AKS and Cluster API in the near future.
