const { LoginPage } = require('../page-objects/login')
const { IndexPage } = require('../page-objects/index')
const { NewTeamPage } = require('../page-objects/teams/new-team')
const { TeamPage } = require('../page-objects/teams')
const { NewClusterPage } = require('../page-objects/teams/clusters/new-cluster')
const { ClusterPage } = require('../page-objects/teams/clusters')
const canonical = require('../../../lib/utils/canonical')
const asyncForEach = require('../../../lib/utils/async-foreach')

const page = global.page

export class TeamCluster {
  constructor({ provider, plan, timeouts }) {
    this.provider = provider
    this.cloud = this.getCloud(provider)
    this.plan = plan
    this.clusterCreateTimeout = timeouts.create
    this.clusterDeleteTimeout = timeouts.delete
  }

  getCloud(provider) {
    return {
      'GKE': 'GCP',
      'EKS': 'AWS'
    }[provider.toUpperCase()]
  }

  run() {
    const testTeam = {
      // Randomise name of test team
      name: `${this.provider} test ${Math.random().toString(36).substr(2, 5)}`,
      description: 'Team created by automation testing'
    }
    const teamID = canonical(testTeam.name)
    const clusterName = teamID
    const namespaces = ['development', 'staging']

    const loginPage = new LoginPage(page)
    const indexPage = new IndexPage(page)
    const newTeamPage = new NewTeamPage(page)
    const teamPage = new TeamPage(page, teamID)
    const newClusterPage = new NewClusterPage(page, teamID)
    const clusterPage = new ClusterPage(page, teamID, clusterName)

    describe('Login', () => {
      it('logs in as an admin user', async () => {
        // admin user login
        await loginPage.visitPage()
        await loginPage.localUserLogin()

        // main dashboard
        indexPage.verifyPageURL()
      })
    })

    describe('Team with cluster and namespaces', () => {
      beforeAll(async () => {
        await indexPage.visitPage()
        indexPage.verifyPageURL()
      })

      describe('Creating team', () => {
        it('navigates to the new team page', async () => {
          await indexPage.clickMenuLink('New team')
        })

        it('creates a new team', async () => {
          await newTeamPage.populate(testTeam)
          await newTeamPage.save()
          await expect(page).toMatch('Team created')
          await newTeamPage.checkTeamID(teamID)
          await newTeamPage.skipToTeamDashboard()
          await indexPage.checkForMenuLink(testTeam.name)
          console.log('Team created', teamID)
        })

        it('show empty team dashboard', async () => {
          teamPage.verifyPageURL()
          await expect(page).toMatch('No clusters found for this team')
          console.log('No clusters found for team', teamID)
        })
      })

      describe('Creating cluster', () => {
        it('navigates to the new cluster page', async () => {
          await teamPage.newCluster()
          newClusterPage.verifyPageURL()
        })

        it('requests a cluster', async () => {
          await newClusterPage.populate({
            cloud: this.cloud,
            plan: this.plan,
            name: clusterName
          })
          await newClusterPage.save()
          console.log('Requested cluster with name', clusterName)
        })

        it('shows the cluster page with cluster status as pending', async () => {
          clusterPage.verifyPageURL()
          await clusterPage.checkClusterName()
          await clusterPage.checkForClusterStatus('Pending')
          await expect(page).toMatch('No namespaces found for this cluster')
          console.log('No namespaces found for cluster', clusterName)
        })

        // overall timeout of 15 minutes
        it('waits for cluster to be successful', async () => {
          // wait up to 12 minutes for the cluster to be created
          await clusterPage.waitForClusterStatus(['Success', 'Failure'], this.clusterCreateTimeout - 30)
          await clusterPage.checkForClusterStatus('Success')
          console.log('Cluster created successfully', clusterName)
        }, this.clusterCreateTimeout)
      })

      describe('Creating namespaces', () => {
        beforeEach(async () => {
          await clusterPage.closeAllNotifications()
        })

        it('creates namespaces', async () => {
          await asyncForEach(namespaces, async (namespace) => {
            await clusterPage.addNamespace()
            await clusterPage.populateNamespace(namespace)
            await clusterPage.saveNamespace()
            await expect(page).toMatch(`Namespace "${namespace}" requested`)
            console.log('Requested namespace with name', namespace)

            // wait up to 10 seconds for the namespace to be created
            await clusterPage.waitForNamespaceStatus(namespace, ['Success', 'Failure'], 10)
            await clusterPage.checkForNamespace(namespace, 'Success')

            await expect(page).toMatch(`Namespace "${namespace}" created`)
            console.log('Namespace created successfully', namespace)
          })
        })
      })

      describe('Post-creation checks', () => {
        it('navigates to the team page and shows the cluster and namespaces', async () => {
          await teamPage.visitPage()
          teamPage.verifyPageURL()
          await teamPage.checkForCluster({ name: clusterName, plan: this.plan, namespaces })
          console.log('Cluster found as expected on team page')
        })

        it('navigates to the cluster page and shows the namespaces', async () => {
          await clusterPage.visitPage()
          clusterPage.verifyPageURL()
          await asyncForEach(namespaces, async (namespace) => {
            await clusterPage.checkForNamespace(namespace)
          })
          console.log('Namespaces found as expected on cluster page')
        })
      })

      describe('Deleting namespaces', () => {
        beforeEach(async () => {
          await clusterPage.closeAllNotifications()
        })

        it('deletes namespaces', async () => {
          await asyncForEach(namespaces, async (namespace) => {
            await clusterPage.deleteNamespace(namespace)
            await clusterPage.deleteNamespaceConfirm()
            await expect(page).toMatch(`Namespace deletion requested: ${clusterName}-${namespace}`)
            console.log('Requested deletion of namespace with name', namespace)

            // wait up to 10 seconds for the namespace to be deleted
            await clusterPage.waitForNamespaceDeleted(namespace, 10)
            await expect(page).toMatch(`Namespace "${namespace}" deleted`)
            console.log('Namespace deleted successfully', namespace)
          })
        })

        it('navigates to the cluster page and shows no namespaces', async () => {
          await clusterPage.visitPage()
          clusterPage.verifyPageURL()
          await expect(page).toMatch('No namespaces found for this cluster')
          console.log('Cluster found as expected on team page')
        })
      })

      describe('Deleting cluster', () => {
        beforeAll(async () => {
          await teamPage.visitPage()
          teamPage.verifyPageURL()
        })

        // overall timeout of 10 minutes
        it('deletes the cluster', async () => {
          await teamPage.deleteCluster(clusterName)
          await teamPage.deleteClusterConfirm()
          await expect(page).toMatch(`Cluster deletion requested: ${clusterName}`)
          await teamPage.checkForClusterStatus(clusterName, 'Deleting')
          console.log('Requested deletion of cluster', clusterName)

          // wait up to 9 minutes for the cluster to be deleted
          await teamPage.waitForClusterDeleted(clusterName, this.clusterDeleteTimeout - 30)
          await expect(page).toMatch(`Cluster successfully deleted: ${clusterName}`)
          console.log('Cluster deleted successfully', clusterName)
        }, this.clusterDeleteTimeout)

        it('navigates to the team page and shows no cluster', async () => {
          await teamPage.visitPage()
          teamPage.verifyPageURL()
          await expect(page).toMatch('No clusters found for this team')
          console.log('No clusters found for team', teamID)
        })
      })

      describe('Deleting team', () => {
        it('deletes the team', async () => {
          await teamPage.deleteTeam()
          await teamPage.deleteTeamConfirm()
          await expect(page).toMatch(`Team "${teamID}" deleted`)
          console.log('Team deleted', teamID)
        })

        it('shows the index page', async () => {
          indexPage.verifyPageURL()
        })
      })

    })
  }
}
