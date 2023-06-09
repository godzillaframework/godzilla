/**
@author: Krisna Pranav, GodzillaFrameworkDevelopers
@filename: internal/errors/errors.go

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

package errs

import (
	"errors"
	"fmt"
	"strings"
)

type Errors interface {
	Errors() []error
}

func joi(errs ...error) error {
	var agg error
	for _, err := range errs {
		if err == nil {
			continue
		} else if agg == nil {
			agg = err
			continue
		} else if errors.Is(err, agg) {
			agg = fmt.Errorf("%w. %s", agg, err)
		} else {
			agg = fmt.Errorf("%s. %s", agg, err)
		}
	}
	return agg
}

func format(err error) string {
	lines := strings.Split(err.Error(), ". ")
	lineLen := len(lines)
	stack := make([]string, lineLen)
	j := lineLen - 1
	for i := 0; i < lineLen; i++ {
		line := lines[j]
		if i > 0 {
			// line = " " + interncolor.dim(line)
			line = " " + line
		}
		stack[i] = line
		j--
	}
	return strings.Join(stack, "\n")
}
