## Release v0.2.3 - 2020-06-29

Bugfix release to enable [Namespaces not appearing on cluster page](https://github.com/appvia/kore/issues/1016)

## Release v0.2.2 - 2020-06-26

Bugfix release to enable [GSuite token refresh to continue to work for Kore CLI after first login](https://github.com/appvia/kore/issues/1006)

## Release v0.2.1 - 2020-06-25

Bugfix release to enable [GSuite token refresh to work for Kore CLI](https://github.com/appvia/kore/issues/999)

## Release v0.2.0 - 2020-05-30

This release of Kore adds the following key features:

- Security Overview of plans and clusters.
- Cloud Services (including S3 and managed binding of services into cluster namespaces via Kore-provisioned secrets).
- GCP Account Management - auto-creation of projects in a configurable way as teams request clusters.
- GCP GKE Improvements - node pools, auto-scale/sizing re-configuration, release channels and versions are now fully supported.
- Enhanced management of cluster plans/parameters to allow cluster users and node groups/pools to be managed easily through the UI.

### Added

- **Security Overview of plans and clusters**
    - [Introduce Security Controller to scan plans and clusters](https://github.com/appvia/kore/issues/437)
    - [UI visibility and overviews](https://github.com/appvia/kore/issues/436)
- **Cloud Services**
    - [AWS Broker](https://github.com/appvia/kore/issues/707)
    - [S3](https://github.com/appvia/kore/issues/628)
    - [Open Service Broker Client](https://github.com/appvia/kore/issues/710)
    - [Service Credentials](https://github.com/appvia/kore/issues/644)
    - [Plans & Services](https://github.com/appvia/kore/issues/642)
- **GCP Account Management**
    - [Update Kore GCP Account Management](https://github.com/appvia/kore/issues/648)
    - [Add CLI support](https://github.com/appvia/kore/issues/649)
    - [Add UI support](https://github.com/appvia/kore/issues/650)
    - [Support Multiple Clusters in Accounts](https://github.com/appvia/kore/issues/828)
- **Cluster / Plan Management UI Improvements**
    - [Add cluster users custom control for GKE and EKS clusters](https://github.com/appvia/kore/issues/539) ([PR](https://github.com/appvia/kore/pull/823))
    - [EKS Node Group management and validation](https://github.com/appvia/kore/pull/835)
- **GCP GKE Improvements**
    - [GKE node pool support](https://github.com/appvia/kore/issues/539) - node pools can now be added, edited and removed on 
      both new and existing clusters
    - [GKE auto-scale, sizing, release channel and version management support](https://github.com/appvia/kore/pull/876) - these 
      properties can now be set for new clusters and edited on existing ones.
- **Other Additions**
    - [Added verbosity flags to helm chart](https://github.com/appvia/kore/issues/815)
    - [Added feature gate options to the helm chart](https://github.com/appvia/kore/issues/807)
    - [Added the ability to lookup allocation on CLI](https://github.com/appvia/kore/issues/766)

### Fixes
- [Cleaning up the EKS IAM Roles](https://github.com/appvia/kore/issues/820)
- [Dryrun flag cleanup on Kore CLI](https://github.com/appvia/kore/issues/801)
- [Fixed an issue in the Accounting update for GCP](https://github.com/appvia/kore/pull/805)
- [Fixed an issue with deletion block on namespaces](https://github.com/appvia/kore/issues/800)
- [Fixed logging format on security controller](https://github.com/appvia/kore/issues/796)
- [Fixed up the age format on kore cli to human readable](https://github.com/appvia/kore/issues/614)
- [Fixed an issue with streaming on the auth-proxy](https://github.com/appvia/kore/issues/704)
- [Fixed the overwriting of inbuilt plans](https://github.com/appvia/kore/issues/755)
- [Added a team check on the CLI](https://github.com/appvia/kore/pull/733)
- [Fixed bug in the create secret CLI](https://github.com/appvia/kore/issues/743)
- [Fixed up the EKS reconciliation controller to be state driven](https://github.com/appvia/kore/pull/576)
- [Renamed packages for store](https://github.com/appvia/kore/pull/719)
- [Fixed verbosity in the CLI config](https://github.com/appvia/kore/issues/697)
- [Fixed issue of local overwriting config](https://github.com/appvia/kore/issues/698)
- [Added kind to allocation names to prevent allocation naming clashes](https://github.com/appvia/kore/issues/679)
- [Fixed Kore UI setup wizard can't be completed for AWS](https://github.com/appvia/kore/issues/600)
- [Fixed AWS EKS Nodegroup desired size should be validated](https://github.com/appvia/kore/issues/611)
- [Fixed Minor: Panic getting cluster when not logged in](https://github.com/appvia/kore/issues/542)

## Release v0.1.0 - 2020-04-27

The first beta release of Kore delivers the following key themes:
- Added GCP Projects for teams (CLI & UI).
- Single account EKS support.
- Clusters by Plans only. Ability for operators to control the shape of the plans which teams consume.
- A formalized release process for Kore.
- The addition of a team audit trail of operations, view-able from CLI and UI.
- A new CLI - `kore`.

### Added

- **GCP Projects for teams (CLI & UI)**
    - GCP Projects [PR #327](https://github.com/appvia/kore/pull/327)
    - GCP Organization Credentials Verification [PR #457](https://github.com/appvia/kore/pull/457)
    - GCP IAM Permissions Check [PR #368](https://github.com/appvia/kore/pull/368)
- **EKS Support**
    - Basic EKS Cluster Support [#348](https://github.com/appvia/kore/issues/348)
    - Add EKS Cluster Build support in the UI [#402](https://github.com/appvia/kore/issues/402)
- **Plans and Policies**
    - limit cluster creation to be only available via a plan [#343](https://github.com/appvia/kore/issues/343)
    - As a kore admin be able to get, create, update plans in the UI [#442](https://github.com/appvia/kore/issues/442)
    - Be able to edit an existing cluster configuration once it's built [#567](https://github.com/appvia/kore/issues/567)
    - Create plan policies associated with teams and plans [#536](https://github.com/appvia/kore/issues/536)
    - View cluster plans [#420](https://github.com/appvia/kore/issues/420)
    - Cluster Deployment from Plans [#346](https://github.com/appvia/kore/issues/346)
- **Audit**
    - Improved Auditing [PR #331](https://github.com/appvia/kore/pull/331)
- **Stable Kore Releases**
    - A release process for kore [#345](https://github.com/appvia/kore/issues/345)
    - Testing: Create basic infrastructure in which we can produce API-level tests [#334](https://github.com/appvia/kore/issues/334)
- **Kore CLI**
    - Add 'local' command to kore cli [#590](https://github.com/appvia/kore/issues/590)
    - Tab Completion of Resource [#640](https://github.com/appvia/kore/issues/640)
    - Stop building korectl [#636](https://github.com/appvia/kore/issues/636)
    - Korectl Bash Autocompletion [#424](https://github.com/appvia/kore/issues/424)

- **Other Additions**
    - UI: See detailed status of a cluster and its underlying objects [#494](https://github.com/appvia/kore/issues/494)
    - UI: Allow deletion of a cluster in failed state [#493](https://github.com/appvia/kore/issues/493)
    - have a way of knowing how and where to authenticate to when i have a cluster [#341](https://github.com/appvia/kore/issues/341)
    - Know the API endpoint for `kore` via the UI [#443](https://github.com/appvia/kore/issues/443)
    - remove alpha from kore local documentation [#687](https://github.com/appvia/kore/issues/687)
    - make the version displayed at the bottom of every page discretely [#686](https://github.com/appvia/kore/issues/686)
    - Create & Delete Admin [PR #483](https://github.com/appvia/kore/pull/483)
    - Korectl Cluster Kubeconfig [#444](https://github.com/appvia/kore/issues/444)
    - GKE Kubernetes Version [PR #462](https://github.com/appvia/kore/pull/462)
    - as an owner of a team i want to be able to delete my team [#430](https://github.com/appvia/kore/issues/430)
    - Apply From Stdin [#408](https://github.com/appvia/kore/issues/408)
    - invite an admin user to kore  [#360](https://github.com/appvia/kore/issues/360)
    - promote a user to be an admin of kore [#359](https://github.com/appvia/kore/issues/359)

### Changes
- Rename API endpoint ekss/ to eks/ to match other resources [#675](https://github.com/appvia/kore/issues/675)
- As a kore administrator I expect to be able to deploy kore using helm with a URL [#624](https://github.com/appvia/kore/issues/624)

### Fixes
- Store Client hiding errors [#671](https://github.com/appvia/kore/issues/671)
- Kubernetes Deletion Bug [#655](https://github.com/appvia/kore/issues/655)
- Team Deletion Bug [#654](https://github.com/appvia/kore/issues/654)
- Always check on the team API endpoints whether a team exists [PR #651](https://github.com/appvia/kore/pull/651)
- Handle validation errors in kore [PR #637](https://github.com/appvia/kore/pull/637)
- EKS resource type doesn't work in the `kore` CLI [#607](https://github.com/appvia/kore/issues/607)
- Create/delete admin is broken in the `kore` CLI [#606](https://github.com/appvia/kore/issues/606)
- Handle trailing slashes for kore profile configure command [#638](https://github.com/appvia/kore/issues/638)
- Oauth Redirect URL Bug Fix [PR #490](https://github.com/appvia/kore/pull/490)
- Helm Chart HMAC  [#476](https://github.com/appvia/kore/issues/476)
- Fix Max Attempts on Resource Waits [PR #452](https://github.com/appvia/kore/pull/452)
- Deleting GKE before Cluster Fix [PR #435](https://github.com/appvia/kore/pull/435)
- Controller Runtime Client GVK [PR #434](https://github.com/appvia/kore/pull/434)
- Fix namespace claim wait [PR #431](https://github.com/appvia/kore/pull/431)
- GKE Bootstrap Fix [PR #454](https://github.com/appvia/kore/pull/454)
- Team Removal with Clusters Bug [PR #426](https://github.com/appvia/kore/pull/426)
- Namespace required when setting Allocation via API [#417](https://github.com/appvia/kore/issues/417)
- Fix Cloud Provider failing state [PR #407](https://github.com/appvia/kore/pull/407)
- fix invalid gke status when no credentials [PR #405](https://github.com/appvia/kore/pull/405)
- Fix invalid status for GKE nodes [#404](https://github.com/appvia/kore/issues/404)
- Make relevant audit visible to the correct personas [#299](https://github.com/appvia/kore/issues/299)
- korectl cls auth doesn't give any context help [#296](https://github.com/appvia/kore/issues/296)
- korectl cls auth on it's own gives 404 [#295](https://github.com/appvia/kore/issues/295)
- Use the refresh token to refresh expired credentials [PR #325](https://github.com/appvia/kore/pull/325)
- When creating a team with an invalid name, it returns 500 [#367](https://github.com/appvia/kore/issues/367)

### Other
- Examples use different teams and invalid plans which make them hard to use [#633](https://github.com/appvia/kore/issues/633)
- Kubernetes Images [PR #463](https://github.com/appvia/kore/pull/463)
- Default Option Value [PR #428](https://github.com/appvia/kore/pull/428)
- Add kubernetes resource config, clean up custom path override logic [PR #414](https://github.com/appvia/kore/pull/414)
- Makefile Images Stage Fix [PR #409](https://github.com/appvia/kore/pull/409)
- Kore Admin Delete Check [PR #397](https://github.com/appvia/kore/pull/397)
- Enable simple kore UI releases with kore [#389](https://github.com/appvia/kore/issues/389)
- Auto Generated Resources  [PR #388](https://github.com/appvia/kore/pull/388)
- CLI Cleanup & Resource Waits [PR #386](https://github.com/appvia/kore/pull/386)
- Fix deepcopy-gen [PR #383](https://github.com/appvia/kore/pull/383)
- JSONPath Dependency [PR #382](https://github.com/appvia/kore/pull/382)
- WhoAmI Command Rendering Fix [PR #381](https://github.com/appvia/kore/pull/381)
- Team Name Regex [PR #375](https://github.com/appvia/kore/pull/375)
- Long Description Format [PR #374](https://github.com/appvia/kore/pull/374)
- Profile Delete Command [PR #373](https://github.com/appvia/kore/pull/373)
- Potential Segmentation Fault [PR #370](https://github.com/appvia/kore/pull/370)
- UI Endpoint [PR #353](https://github.com/appvia/kore/pull/353)
- Enable APIs [PR #349](https://github.com/appvia/kore/pull/349)
- Logrus Caller Report  [#330](https://github.com/appvia/kore/issues/330)

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
