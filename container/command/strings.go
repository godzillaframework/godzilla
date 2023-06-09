/**
@author: Krisna Pranav, GodzillaFrameworkDevelopers
@filename: container/command/strings.go

Copyright [2021 - 2023] [Krisna Pranav, GodzillaFrameworkDeveloeprs]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package command

import (
	"fmt"
	"strings"
)

type Strings struct {
	target   *[]string
	defval   *[]string // default value
	optional bool
}

type stringsValue struct {
	inner *Strings
	set   bool
}

func (v *stringsValue) verify(displayName string) error {
	if v.set {
		return nil
	} else if v.inner.defval != nil {
		*v.inner.target = *v.inner.defval
		return nil
	} else if v.inner.optional {
		return nil
	}
	return fmt.Errorf("missing %s", displayName)
}

func (v *Strings) Default(values ...string) {
	v.defval = &values
}

func (v *stringsValue) Set(val string) error {
	*v.inner.target = append(*v.inner.target, val)
	v.set = true
	return nil
}

func (v *stringsValue) String() string {
	if v.inner == nil {
		return ""
	} else if v.set {
		return strings.Join(*v.inner.target, ", ")
	} else if v.inner.defval != nil {
		return strings.Join(*v.inner.defval, ", ")
	}
	return ""
}
