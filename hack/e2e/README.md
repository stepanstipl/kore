## **Running E2E**

### **Prerequisites**

- kubectl
- korectl
- bats (bash tests)

You must have some validate credentials in the e2e folder

- e2e/gke-credentials (containing a gkecredentials{} and allocation called gke for all teams)

##  **Running the suite locally**

- Bring up the dependencies via `make compose` or `make demo`
- If your not using the `make demo` bring up the kore-apiserver locally via `bin/kore-apiserver --verbose`; sourcing
  in any environment variables you usually do.
- Login via the `korectl login` command which will provision your user locally

```shell
korectl 

```

- Ensure if you using multiple profiles your pointing to the local instance `korectl profiles ls`

You can then run the checks via:

```shell
# The prefix
$ hack/e2e/check-suite.sh
```

Note: any clusters it will create are being the environment definition; `CLUSTER="ci-${CIRCLE_BUILD_NUM:-$USER}"`.
