## **Running E2E**

### **Prerequisites**

- kubectl
- korectl
- bats (bash tests)

You must have some validate credentials in the e2e folder

- e2e/gke-credentials (containing a gkecredentials{} and allocation called gke for all teams)

You can then run the checks via:

```shell
$ hack/e2e/check-suite.sh
```
