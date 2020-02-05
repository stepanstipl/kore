/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package register

import (
	// importing the cloud providers
	_ "github.com/appvia/kore/pkg/controllers/cloud/aws/eks"
	_ "github.com/appvia/kore/pkg/controllers/cloud/gcp/gke"
	_ "github.com/appvia/kore/pkg/controllers/cloud/gcp/gkecredentials"

	// importing the management controller
	_ "github.com/appvia/kore/pkg/controllers/management/bootstrap"
	_ "github.com/appvia/kore/pkg/controllers/management/clusterbindings"
	_ "github.com/appvia/kore/pkg/controllers/management/clusterconfig"
	_ "github.com/appvia/kore/pkg/controllers/management/clusterroles"
	_ "github.com/appvia/kore/pkg/controllers/management/kubernetes"
	_ "github.com/appvia/kore/pkg/controllers/management/podpolicy"
	_ "github.com/appvia/kore/pkg/controllers/management/namespaceclaims"

	// importing the user controllers
	_ "github.com/appvia/kore/pkg/controllers/user/allocations"
	_ "github.com/appvia/kore/pkg/controllers/user/teams"
)
