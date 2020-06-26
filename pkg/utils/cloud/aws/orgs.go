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

package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"
)

// CheckOUExist checks if the OU specified for an OU exists
func CheckOUExist(s *session.Session, ouName string) (bool, error) {
	ouID, err := GetOUID(s, ouName)
	if err != nil {

		return false, err
	}
	if ouID != nil {

		return true, nil
	}

	return false, nil
}

// GetOUID obtains the OU ID
func GetOUID(s *session.Session, ouName string) (*string, error) {
	// First we need the organisation root ID:
	roots, err := GetRoots(s)
	if err != nil {

		return nil, err
	}
	l := len(roots)
	if l != 1 {

		return nil, fmt.Errorf("can only support one organisational root but found %d", l)
	}
	rootID := aws.StringValue(roots[0].Id)
	ous, err := ListOUsFromParent(s, rootID)
	if err != nil {

		return nil, err
	}
	var ouID *string
	for _, ou := range ous {
		if *ou.Name == ouName {
			if ouID != nil {

				return nil, fmt.Errorf("more than one ou matching %s from root with id %s", ouName, rootID)
			}

			ouID = ou.Id
		}
	}
	if ouID == nil {

		return nil, fmt.Errorf("invalid OU Specified - %s", ouName)
	}

	return ouID, nil
}

// ListOUsFromParent will obtain the OUs from a given parent ID (id of root or OR)
func ListOUsFromParent(s *session.Session, parentID string) ([]*organizations.OrganizationalUnit, error) {
	orgSvc := organizations.New(s)
	ouo, err := orgSvc.ListOrganizationalUnitsForParent(&organizations.ListOrganizationalUnitsForParentInput{
		ParentId: aws.String(parentID),
	})
	if err != nil {

		return nil, fmt.Errorf("unbale to list organisational units %w", err)
	}
	return ouo.OrganizationalUnits, nil
}

// GetRoots gets all the Root IDs in an organisation
func GetRoots(s *session.Session) ([]*organizations.Root, error) {
	// First we need the organisation root ID:
	orgSvc := organizations.New(s)
	lro, err := orgSvc.ListRoots(&organizations.ListRootsInput{})
	if err != nil {

		return nil, fmt.Errorf("unbale to list organisation roots - %w", err)
	}
	return lro.Roots, nil
}
