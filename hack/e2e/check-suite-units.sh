#!/bin/bash
#
# Copyright 2020 Appvia Ltd <info@appvia.io>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
source hack/e2e/environment.sh || exit 1

announce "running the integretion suite"

RET=0
BATS_OPTIONS=""

cd ${PLATFORM_DIR}/integration
announce "performing the setup for unit test "
bats ${BATS_OPTIONS} setup.bats || exit 1
announce "performing the checks for proiles"
bats ${BATS_OPTIONS} profiles.bats || exit 1
announce "performing the checks for users"
bats ${BATS_OPTIONS} users.bats || exit 1
announce "performing the checks for whoami"
bats ${BATS_OPTIONS} whoami.bats || exit 1
announce "performing the checks for teams"
bats ${BATS_OPTIONS} teams.bats || exit 1
announce "performing the tests for plans"
bats ${BATS_OPTIONS} plans.bats || exit 1
announce "performing the gke credentials"
bats ${BATS_OPTIONS} gke-credentials.bats || exit 1
announce "performing checks on gke clusters"
bats ${BATS_OPTIONS} gke.bats || exit 1
announce "performing checks on clusterappman"
bats ${BATS_OPTIONS} clusterappman.bats || exit 1
announce "performing checks on clusterroles"
bats ${BATS_OPTIONS} clusterroles.bats || exit 1
announce "performing checks on clusterconfig"
bats ${BATS_OPTIONS} clusterconfig.bats || exit 1
announce "performing checks on namespaces"
bats ${BATS_OPTIONS} namespaces.bats || exit 1
announce "performing checks on cluster delete"
bats ${BATS_OPTIONS} gke-deletion.bats || exit 1

exit $RET
