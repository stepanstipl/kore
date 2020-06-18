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

import "sort"

// Contains checks a list has a value in it
func Contains(v string, l []string) bool {
	for _, x := range l {
		if v == x {
			return true
		}
	}

	return false
}

// ChunkBy breaks a slice into chunks
func ChunkBy(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}

// Unique removes any duplicates from a slice
func Unique(items []string) []string {
	var list []string

	found := make(map[string]bool)

	for _, x := range items {
		if ok := found[x]; !ok {
			list = append(list, x)
		}
	}

	return list
}

// StringsSorted returns a 'copy' of a sorted list of strings
func StringsSorted(list []string) []string {
	v := make([]string, len(list))
	copy(v, list)

	sort.Strings(v)

	return v
}

func StringSliceFrom(v interface{}) ([]string, bool) {
	if v == nil {
		return nil, true
	}

	switch vt := v.(type) {
	case []string:
		return vt, true
	case []interface{}:
		var res []string
		for _, e := range vt {
			s, ok := e.(string)
			if !ok {
				return nil, false
			}
			res = append(res, s)
		}
		return res, true
	default:
		return nil, false
	}
}

type StringSet []string

func (s *StringSet) Add(v string) {
	for _, e := range *s {
		if e == v {
			return
		}
	}
	*s = append(*s, v)
}

func (s *StringSet) Remove(v string) {
	for i, e := range *s {
		if e == v {
			*s = append((*s)[0:i], (*s)[i+1:]...)
		}
	}
}

func (s *StringSet) MemberIf(v string, cond bool) {
	if cond {
		s.Add(v)
	} else {
		s.Remove(v)
	}
}

func (s *StringSet) Contains(v string) bool {
	for _, e := range *s {
		if e == v {
			return true
		}
	}
	return false
}
