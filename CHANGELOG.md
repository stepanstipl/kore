## Release v0.1.0 - 2020-04-24
## Added
add check for release notes (and generator) ([e8cc014](https://github.com/appvia/kore/commit/e8cc0140a6e662bac6679f795a68d495e44ad19b))
- fixing the documentation ([49771ad](https://github.com/appvia/kore/commit/49771ad8924db686727a0586ad01548de1cd79ad))
- fixing the linting issue ([10ddc90](https://github.com/appvia/kore/commit/10ddc9022875960c1ba37744d7e576d29147641f))
- adding the ability to delete as well ([458f2eb](https://github.com/appvia/kore/commit/458f2eb5f9b5ff0526343c75699293fa5d7c322d))
- fixing up the commend on the method ([743bd7d](https://github.com/appvia/kore/commit/743bd7de15d582224ac4e7a2b22a91460e6a7153))
- fixing up what was highlighted in the review ([4e6958b](https://github.com/appvia/kore/commit/4e6958b75981a207a6998495504a159def0069c0))
- fixing the linting issue ([0ac2dae](https://github.com/appvia/kore/commit/0ac2daea5163007d53b3283cb978fdee0a1e4eb7))
- fixing the wording ([b0805f2](https://github.com/appvia/kore/commit/b0805f2d4b0da96a0992a823b79ce5a6a613fb3f))
- taking the opportunity to use sjson instead and adding some examples ([24ca17a](https://github.com/appvia/kore/commit/24ca17a67ab974340b1731a5f28a9ecd0142e5e9))
- fixing the parsing of the parameter as json was failing ([c722a20](https://github.com/appvia/kore/commit/c722a20f1c2979f468f8634cfec22643c892ac35))
- updating the path to the kore configuration file ([c0addb9](https://github.com/appvia/kore/commit/c0addb9b879cea3ee6bef5ffab49f0e9865823d1))
Handle validation errors in kore ([81a23d2](https://github.com/appvia/kore/commit/81a23d253920786fe02e99912cfe2d7fa03a97d1))
When logging in with Kore, start to callback http server on localhost, to avoid firewall warnings ([773b9c3](https://github.com/appvia/kore/commit/773b9c3693de0080c967b869df9f20a2bae9f4e6))
Refactor Kore CLI to create an API client from a resource ([2e255e9](https://github.com/appvia/kore/commit/2e255e9f9f1e161785eef95e64ea15b3101ef17e))
- testing out on the branch ([2e5b2ac](https://github.com/appvia/kore/commit/2e5b2acadea733c9fd51ea6f036458d6414262de))
The deployment wasn't changing, fixing up so we always use the sha1 ([38837c5](https://github.com/appvia/kore/commit/38837c57676df17f22989e2f998d96bbb98fe843))
- fixing up the review comments ([4a4c0ec](https://github.com/appvia/kore/commit/4a4c0ec129ce2a14e938028950ced8edbf855c4d))
UI end-to-end testing ([1810498](https://github.com/appvia/kore/commit/1810498aaaae24f5f49323010f4802b67485d4de))
- fixing up the references in the UI ([38b1105](https://github.com/appvia/kore/commit/38b11055122ea87a723b971f1709ec5427011de7))
- fixing up the auth-proxy check ([0eb9a78](https://github.com/appvia/kore/commit/0eb9a78b3e511959f876676cf1d1ae2a0a84d26d))
Fixing incorrect plan description on AWS cloud config ([2c463e0](https://github.com/appvia/kore/commit/2c463e07eb5424a622813fe261403e897321b9d8))
- fixing the check as it was being skipped due to the order ([5e3c546](https://github.com/appvia/kore/commit/5e3c54607af8911a80841788238fdc82df561aea))
- changing all the references in kore command ([d7ad029](https://github.com/appvia/kore/commit/d7ad029e5fb2bda926534b69bb2c68930196eee6))
- removing all the references to viewer ([ed78716](https://github.com/appvia/kore/commit/ed7871694dbcd09b34ceadf645123cb20b72eea0))
- dropping the team role paramater and defaulting to use of the plan params ([3adab96](https://github.com/appvia/kore/commit/3adab96cba1958f6660f9ca8f320b38ca4bb7ad7))
Add and edit plans ([cc84693](https://github.com/appvia/kore/commit/cc846932812b81b0d5762f37703086cfce6421b4))
- fixing up the routes in the UI ([7fc5fa8](https://github.com/appvia/kore/commit/7fc5fa81fe62bf7056fdeb75c7f082d7d8cdae88))
- fixing up the plans to the view role ([140567b](https://github.com/appvia/kore/commit/140567b113f2cda3b44c673a5239c721c01fc691))
- using a single printer ([bd0c339](https://github.com/appvia/kore/commit/bd0c3390fcedd2ba81b79b704d751a46ebe72d48))
- fixing up the error being thrown ([a35f8a4](https://github.com/appvia/kore/commit/a35f8a4487f73f2bdb5754b3beadb83bccdfb3f9))
- fixing the wording ([6796c60](https://github.com/appvia/kore/commit/6796c608b4b326b3eb5fc5f7a52de1de2215b60a))
- fixing the check on the all flag ([d13b9db](https://github.com/appvia/kore/commit/d13b9db7ed5163a08dcb4750ebc9d294a4f6bc3e))
- updaing the swagger apiclient ([9984af8](https://github.com/appvia/kore/commit/9984af8d8776810b7fb5d5042fe241ba5b3ee1f9))
- fixing up the audit to be a plural ([0be1c92](https://github.com/appvia/kore/commit/0be1c92c1af79d5d2c55f280ddcd3d1c52ed2a29))
- adding the get audit command ([c60cd6c](https://github.com/appvia/kore/commit/c60cd6c26ff9b9e34e6722d2796ae6690d0a23bd))
- fixing a bug in the api client ([acda2bb](https://github.com/appvia/kore/commit/acda2bb20d2b24ca865d23e9e55b7cea8adcc692))



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
