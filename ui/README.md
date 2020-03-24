# Kore UI

The UI is build using next.js with an ExpressJS server.

### Install and run

**Pre-reqs**

If you have not set up your local kore environment, run the following minimal steps to set it up:

```bash
# Move to root of the kore git clone (i.e. cd ../ from this directory)
cd ../
make korectl
bin/korectl local configure
```

For a full guide to configuring Kore, see [the kore quickstart guide](https://github.com/appvia/kore/blob/master/doc/alpha-local-quick-start.md).

To run the kore API server and its dependencies, in a separate terminal, do the following:

```bash
# Move to root of the kore git clone (i.e. cd ../ from this directory)
cd ../
# build and run dependencies
make run
# check it's running
curl http://localhost:10080/api/v1alpha1/whoami -H 'Authorization: Bearer password'
```

To run the UI, return to the ui directory and run:

```bash
npm install
make run
```

Visit http://localhost:3000 in the browser.
Login using the default admin user credentials: admin / password or your configured auth provider.

**Production**

To run in production mode, do the following

```bash
npm install
npm run build
npm start
```
