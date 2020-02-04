# User Management

While the identity provider provides authentication, the hub removes the need for complication roles permission being maintained in the user claims. Instead we use the identity provider for authentication and maintain our own internal control over the permissions via roles the user can permit. Note at present the scope on early release is to border permissions between teams and treat all members in the team admin inside that team as admin. In the near will will then break down responisibilies via roles within the teams themselves. 

Pretty much all the administrative features can be handled by the CLI. Once you have authenticated against the API ($ korectl auth) and global admin can then manipulate the teams and user management via the $ korectl teams <subcommands>

## Creating Teams

Teams are created by any member of the admin group; though this will change in the near future to be less restrictive. $ korectl teams apply -f <path>

## Adding Users
Adding users is again carried out via the CLI: 

`korectl teams mb add --team <name> --user <name>` 

**Note as of now the user must exist in the hub.**

## Cluster and User Roles

Kore provides the ability to manage an assortment of kubernetes clusters roles and bindings across the estate. These much like other policies are applied in the order or Global, Team and cluster. An example of a managed cluster role can be found in the examples/policies folder.

