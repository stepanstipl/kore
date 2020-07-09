// +build !ignore_autogenerated

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1beta1

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Continent) DeepCopyInto(out *Continent) {
	*out = *in
	if in.Regions != nil {
		in, out := &in.Regions, &out.Regions
		*out = make([]Region, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Continent.
func (in *Continent) DeepCopy() *Continent {
	if in == nil {
		return nil
	}
	out := new(Continent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContinentList) DeepCopyInto(out *ContinentList) {
	*out = *in
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Continent, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContinentList.
func (in *ContinentList) DeepCopy() *ContinentList {
	if in == nil {
		return nil
	}
	out := new(ContinentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Cost) DeepCopyInto(out *Cost) {
	*out = *in
	out.Resource = in.Resource
	in.From.DeepCopyInto(&out.From)
	in.To.DeepCopyInto(&out.To)
	in.RetrievedAt.DeepCopyInto(&out.RetrievedAt)
	if in.CostElements != nil {
		in, out := &in.CostElements, &out.CostElements
		*out = make([]CostElement, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Cost.
func (in *Cost) DeepCopy() *Cost {
	if in == nil {
		return nil
	}
	out := new(Cost)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CostElement) DeepCopyInto(out *CostElement) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CostElement.
func (in *CostElement) DeepCopy() *CostElement {
	if in == nil {
		return nil
	}
	out := new(CostElement)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CostEstimate) DeepCopyInto(out *CostEstimate) {
	*out = *in
	if in.CostElements != nil {
		in, out := &in.CostElements, &out.CostElements
		*out = make([]CostEstimateElement, len(*in))
		copy(*out, *in)
	}
	in.PreparedAt.DeepCopyInto(&out.PreparedAt)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CostEstimate.
func (in *CostEstimate) DeepCopy() *CostEstimate {
	if in == nil {
		return nil
	}
	out := new(CostEstimate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CostEstimateElement) DeepCopyInto(out *CostEstimateElement) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CostEstimateElement.
func (in *CostEstimateElement) DeepCopy() *CostEstimateElement {
	if in == nil {
		return nil
	}
	out := new(CostEstimateElement)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CostList) DeepCopyInto(out *CostList) {
	*out = *in
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Cost, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CostList.
func (in *CostList) DeepCopy() *CostList {
	if in == nil {
		return nil
	}
	out := new(CostList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceType) DeepCopyInto(out *InstanceType) {
	*out = *in
	if in.Prices != nil {
		in, out := &in.Prices, &out.Prices
		*out = make(map[PriceType]int64, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceType.
func (in *InstanceType) DeepCopy() *InstanceType {
	if in == nil {
		return nil
	}
	out := new(InstanceType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceTypeList) DeepCopyInto(out *InstanceTypeList) {
	*out = *in
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]InstanceType, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceTypeList.
func (in *InstanceTypeList) DeepCopy() *InstanceTypeList {
	if in == nil {
		return nil
	}
	out := new(InstanceTypeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Region) DeepCopyInto(out *Region) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Region.
func (in *Region) DeepCopy() *Region {
	if in == nil {
		return nil
	}
	out := new(Region)
	in.DeepCopyInto(out)
	return out
}