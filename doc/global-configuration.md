# Global Configuration
## Single Sign On

Currently Appvia Kore requires an external identity provider providing authenticate to both API, clusters built and applications provisioned. The options are easily configured on the command line of the api server via --client-id / --client-secret  and --discovery-url (OpenID’s discovery URL). User’s can authorize themselves via the CLI: korectl authorize which will create the user locally within the kore. Note the first user who authenticates is deemed an admin, while the rest will be placed into the kore-default team (a catch all for users).

The same authentication process is also used across the estate. Once in a team, administrators can control user access at a team, cluster or namespace level. Permissions are covered later but for a quick howto once a cluster has been built gaining access requires as korectl clusters auth. The command (assuming you’ve already authenticated) will retrieve a list of your clusters across all teams and provision you kubeconfig accordingly.


## Cloud Provider Configuration

At present the CLI is still under development and so we rely on standard kubectl for integration / administrative configuration; such as credentials, allocations and of sorts. The examples directory contains snippets of configuration for GKE and allocations to team. Note, since it’s been mentioned we should probably explain what an allocation is. In the platform we’ve siloed all the credentials all the main administrative team, thus cloud credentials, policy etc configured by the admins and ‘allocated’ to one or more or all teams. Essentially, it’s simply a means of providing a team with credentials without them having access to it. 

## Cluster Plans

Currently not in scope in the is release, but plans can be thought of as templates for desired state of a resource, be a cluster, a cloud provide or even a namespace. They simply provide the admin to craft a specific set of packages for the teams to consume.

