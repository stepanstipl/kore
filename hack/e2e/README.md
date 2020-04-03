## **Running E2E**

### **Prerequisites**

- bats (bash tests)
- jq
- korectl
- kubectl

You must have some valid credentials in the e2e folder

- e2eci/gke-credentials.yml (containing a gkecredentials{} and allocation called gke for all teams).
  A example can be found in the examples/gcp-credentials.yml

##  **Running the suite locally**

- Bring up the dependencies via `make compose` or `make demo`
- If your not using the `make demo` bring up the kore-apiserver locally via `bin/kore-apiserver --verbose`; sourcing
  in any environment variables you usually do.
- Login via the `korectl login` command which will provision your user locally
- Ensure if you using multiple profiles your pointing to the local instance `korectl profiles ls`

You can then run the checks via:

```shell
$ hack/e2e/check-suite.sh
```

Note: any clusters it will create are being the environment definition; `CLUSTER="ci-${CIRCLE_BUILD_NUM:-$USER}"`.
