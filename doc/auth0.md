# Auth0

Auth0, found [here](https://auth0.com/), provides an enterprise SAAS identity provider.

- Sign up for an account from the [home page](https://auth0.com)
- From the dashboard side menu choose 'Applications' and then 'Create Application'
- Given the application a name and choose 'Regular Web Applications'
- Once provisioned click on the 'Settings' tab and scroll down to 'Allowed Callback URLs'. These are the permitted redirects for the applications. If we are running the application locally these will be `http://localhost:3000/auth/callback` and `http://localhost:10080/oauth/callback` (Note the comma separation in the Auth0 UI.
- Scroll to the bottom of the settings and click the 'Show Advanced Settings'
- Choose the 'OAuth' tab from the advanced settings and ensure that the 'JsonWebToken Signature Algorithm' is set to RS256 and 'OIDC Conformant' is toggled on.
- Select the 'Endpoints' tab and note down the 'OpenID Configuration'.
- You can then scroll back to the top and note down the 'ClientID' and 'Client Secret'
