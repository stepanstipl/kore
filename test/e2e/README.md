## **Running E2E**

### **Prerequisites**

- bats (bash tests)
- jq
- kore
- kubectl

You must have some valid credentials in the e2e folder

 * e2eci/gke-credentials.yml (containing a GKECredentials object and Allocation called gke for all teams).
   A example can be found in the examples/gcp-credentials.yml
 * e2eci/aks-credentials.yml (containing an AKSCredentials object and Allocation called aks for all teams).
   A example can be found in the examples/aks-credentials.yml

##  **Running the suite locally**

- Bring up the dependencies via `make compose` or `make demo`
- If your not using the `make demo` bring up the kore-apiserver locally via `bin/kore-apiserver --verbose`; sourcing
  in any environment variables you usually do.
- Login via the `kore login` command which will provision your user locally
- Ensure if you are using multiple profiles your pointing to the local instance `kore profiles ls`

You can then run the checks via:

```shell
$ test/e2e/check-suite.sh
```

Note: any clusters created will use the environment definition; `export CLUSTER="ci-${CIRCLE_BUILD_NUM:-$USER}"`.
