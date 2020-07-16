/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package register

import (
	// importing the aws controllers
	_ "github.com/appvia/kore/pkg/controllers/cloud/aws/awsaccount"
	_ "github.com/appvia/kore/pkg/controllers/cloud/aws/awsaccountclaims"
	_ "github.com/appvia/kore/pkg/controllers/cloud/aws/awsorganization"
	_ "github.com/appvia/kore/pkg/controllers/cloud/aws/credentials"
	_ "github.com/appvia/kore/pkg/controllers/cloud/aws/eks"
	_ "github.com/appvia/kore/pkg/controllers/cloud/aws/eksnodegroup"
	_ "github.com/appvia/kore/pkg/controllers/cloud/aws/eksvpc"

	// import the gcp controllers
	_ "github.com/appvia/kore/pkg/controllers/cloud/gcp/gke"
	_ "github.com/appvia/kore/pkg/controllers/cloud/gcp/gkecredentials"
	_ "github.com/appvia/kore/pkg/controllers/cloud/gcp/organization"
	_ "github.com/appvia/kore/pkg/controllers/cloud/gcp/projectclaims"
	_ "github.com/appvia/kore/pkg/controllers/cloud/gcp/projects"
	_ "github.com/appvia/kore/pkg/controllers/cloud/gcp/projectsweeper"

	// import the azure controllers
	_ "github.com/appvia/kore/pkg/controllers/cloud/azure/aks"
	_ "github.com/appvia/kore/pkg/controllers/cloud/azure/akscredentials"

	// importing the management controller
	_ "github.com/appvia/kore/pkg/controllers/management/cluster"
	_ "github.com/appvia/kore/pkg/controllers/management/clusterbindings"
	_ "github.com/appvia/kore/pkg/controllers/management/clusterconfig"
	_ "github.com/appvia/kore/pkg/controllers/management/clusterroles"
	_ "github.com/appvia/kore/pkg/controllers/management/kubernetes"
	_ "github.com/appvia/kore/pkg/controllers/management/namespaceclaims"
	_ "github.com/appvia/kore/pkg/controllers/management/podpolicy"

	// import secret controllers
	_ "github.com/appvia/kore/pkg/controllers/secrets/gcp"
	_ "github.com/appvia/kore/pkg/controllers/secrets/generic"

	// importing the user controllers
	_ "github.com/appvia/kore/pkg/controllers/user/allocations"
	_ "github.com/appvia/kore/pkg/controllers/user/teams"

	// importing the service controllers
	_ "github.com/appvia/kore/pkg/controllers/servicecredentials"
	_ "github.com/appvia/kore/pkg/controllers/serviceproviders"
	_ "github.com/appvia/kore/pkg/controllers/services"

	// importing the security controllers
	_ "github.com/appvia/kore/pkg/controllers/security"

	// import the feature controllers
	_ "github.com/appvia/kore/pkg/controllers/korefeatures"
)
