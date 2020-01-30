/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package bootstrap

const (
	// BootstrapDeploymentTemplate is the template for deploying Kore cluster management
	BootstrapDeploymentTemplate = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kore-clusterman
  namespace: kore-system
spec:
  replicas: 1
  selector:
    matchLabels:
      name: kore-clusterman
  strategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        name: kore-clusterman
    spec:
	  containers:
      - name: kore
        image: {{ .KoreImage }}
        command:
        - /kore-clusterman
        env:
        - name: IN_CLUSTER
          value: "true"
`
)
