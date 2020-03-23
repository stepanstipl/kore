Running tests
=============

To run the tests:

1. Ensure the API is running with the localjwt and admintoken authentication plugins 
   enabled (i.e. run with `--kore-authentication-plugin localjwt` and 
   `--kore-authentication-plugin admintoken`)
2. Ensure the API is using `--local-jwt-public-key` with a public key that 
   matches the private key in api_test.go. See below for how to generate a new
   key if needed.
3. If the API is not running on localhost, export KORE_API_HOST with the hostname.
3. Run make api-test.

Adding new tests
================

Each API endpoint should be represented by a file in this directory. The file should include tests for:

* Authentication
* Authorisation
* Validation
* Success (if possible)

There are a number of test users and teams configured in setupdata_test.go. Any other data required
by your test should be configured in your test or added to setupdata_test.go.

Generating new keys for localjwt
================================

To generate new keys, uncomment generateJWTKeys in `api_integration_test.go` and run go test, this 
will output a new keypair to the console. Place both keys in api_integration_test.go and provide the 
public key to the API server on start-up via the --local-jwt-public-key parameter.