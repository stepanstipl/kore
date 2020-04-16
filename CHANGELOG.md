## Release v0.0.23 - 2020-04-16

This release targets being able to use Kore to create, manage and destroy EKS Kubernetes clusters in AWS, with
all dependent resources managed (VPC, Subnets, node groups, etc). 

As this is the first formal release of Kore, it only shows the major features added since the previous pre-release alpha.

### Added

- [Kore cluster automatically creates and uses dependent EKS objects](https://github.com/appvia/kore/issues/450)
  Creating an EKS cluster now creates the required VPC infrastructure.
- [Kore deletes Kore-managed VPC infrastructure when deleting EKS cluster](https://github.com/appvia/kore/issues/492)
  Where Kore has created a VPC and subnets for an EKS cluster, it will now remove them when the cluster is deleted.
- [Secured EKS Endpoints](https://github.com/appvia/kore/issues/514) 
- [Proxy protocol annotation added to AWS Auth Proxy](https://github.com/appvia/kore/issues/505)
- [Configure AWS EKS in the UI](https://github.com/appvia/kore/issues/488)
- [Modify plan parameters when creating a cluster through UI](https://github.com/appvia/kore/issues/489)
- [Packaged release for korectl](https://github.com/appvia/kore/issues/324)

### Improved & Changed

- EKS improvements:
  - https://github.com/appvia/kore/pull/544
  - https://github.com/appvia/kore/pull/527
- [Expose egress IP addresses from EKS VPC object status](https://github.com/appvia/kore/issues/517)

### Fixed

- [Server endpoint incorrect for korectl login from UI](https://github.com/appvia/kore/issues/571)
- [Non-recoverable transient error when trying to create EKS VPC](https://github.com/appvia/kore/issues/548)
