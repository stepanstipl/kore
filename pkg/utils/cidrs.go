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

package utils

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
)

// GetSubnet will get the first network from a new mask
func GetSubnet(ip net.IP, bitMask int) (*net.IPNet, error) {
	sn, err := GetSubnets(ip, bitMask, 1)
	if err != nil {
		return nil, err
	}
	return sn[0], nil
}

// GetSubnetFromLast will get the next network given a new mask
func GetSubnetFromLast(baseNet *net.IPNet, newBitMask int) (*net.IPNet, error) {
	sn, no := cidr.NextSubnet(baseNet, newBitMask)
	if no {
		return nil, fmt.Errorf("cannot create network of x.x.x.x/%d from %s", newBitMask, baseNet.String())
	}
	return sn, nil
}

// GetSubnetsFromCidr will get a sequence of subnets from within a larger network
func GetSubnetsFromCidr(baseCidr string, newBitMask, count int) ([]*net.IPNet, error) {
	_, baseNet, err := net.ParseCIDR(baseCidr)
	if err != nil {
		return nil, err
	}
	return GetSubnets(baseNet.IP, newBitMask, count)
}

// GetSubnets will retrieve a sequence from a startin IP
func GetSubnets(ip net.IP, newBitMask, count int) ([]*net.IPNet, error) {
	// Work out base of each subnet
	nets := make([]*net.IPNet, count)
	_, newNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", ip, newBitMask))
	if err != nil {
		return nil, fmt.Errorf("error working out first network from base %s and new bitmask of %d - %s", ip.String(), newBitMask, err)
	}
	for i := 0; i < count; i++ {
		nets[i] = newNet
		if i < count {
			var notposs bool
			newNet, notposs = cidr.NextSubnet(newNet, newBitMask)
			if notposs {
				return nil, fmt.Errorf("cannot create network of x.x.x.x/%d from %s", newBitMask, newNet.IP)
			}
		}
	}
	return nets, nil
}
