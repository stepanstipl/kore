# Cloud Providers

Kore is designed to work with multiple cloud providers, with a view of supporting on-premise cloud solutions in the future. 

By leveraging cloud provided Kubernetes services, it helps offload some of the operational responsibility back to the cloud operators and help reduce some of the technical responsibility from the organisation. 

Kore provides a consitent way to work with cloud providers as well as their security features.

## Cloud Providers

Kore currently supports
- GCP
- AWS
- Azure (Comming soon)

Within the cloud providers, we will support their relative Kubernetes service offerings: 
- GKE
- EKS
- AKS

## Cloud Security Feature Comparison

Below is a list of known security features we want to enable as part of Kore across different cloud providers. This will be a growing list over time, but hopefully will help people understand the feature comparison between clouds with Kore.

| Security Features | Description | GKE | EKS | AKS |
|:-----------------:|:-----------:|:---:|:---:|:---:|
| Role Based Access Control | Enable RBAC inside the Clusters | x | TBD | TBD |  
| Audit | Enable audit logging of OS and Kubernetes | x | TBD | TBD |
| PSP | Enable pod security policies | x | TBD | TBD |
| Disabled legacy Authentication | Disable basic authentication and Client certificates | x | TBD | TBD |
| OS Image Authenticity | Enable OS image and Kernel module authenticity through signing via a vTPM | x | TBD | TBD |
| Minimal privilege | Make service accounts in GKE minimal privilege for both the nodes and the container workloads | x | TBD | TBD |
| Kubernetes Network Policies | Restrict traffic between pods | x | TBD | TBD |
| Kubernetes Secret encryption | Encrypt secrets inside of Kubernetes with KMS keys | TBD | TBD | TBD |
| Private nodes and Private / Restricted API Endpoint | Make the Kubernetes API and worker nodes private and restricted to known customer networks | TBD | TBD | TBD |
| Node network traffic logs | Traffic logging between nodes and pods | TBD | TBD | TBD |
| Auto node upgrade | Upgrade Kubernetes versions automatically | TBD | TBD | TBD |

## Additional Cloud Security Detail

For a more detailed view of the implementation of the cloud providers security controls, please click on the below
- [GKE Security Information](security-gke.md)
