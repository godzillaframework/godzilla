/**
@author: Krisna Pranav, GodzillaFrameworkDevelopers
@filename: internal/color/color.go

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

package color

import (
	"os"

	"github.com/aybabtme/rgbterm"
)

var noColor = os.Getenv("NO_COLOR") != ""

func paint(msg string, r, b, g uint8) string {
	if noColor {
		return msg
	}

	return rgbterm.FgString(msg, r, g, b)
}

func dim(msg string) string {
	if noColor {
		return msg
	}

	return "\033[37m" + msg + "\033[0m"
}

func bold(msg string) string {
	if noColor {
		return msg
	}

	return "\033[1m" + msg + "\033[0m"
}

// colors
func white(msg string) string {
	return paint(msg, 226, 232, 240)
}
