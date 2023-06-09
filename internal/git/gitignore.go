/**
@author: Krisna Pranav, GodzillaFrameworkDevelopers
@filename: internal/ignore/gitignore.go

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

package gitignore

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

// ignoreable values.
var alwaysIgnore = []string{
	"node_modules",
	".git",
	".DS_Store",
	".vscode/",
	".idea/",
	"bud",
}

// default ignores values
var defaultIgnores = append([]string{"/bud"}, alwaysIgnore...)
var defaultIgnore = gitignore.CompileIgnoreLines(defaultIgnores...).MatchesPath

// from filesystem
func fromFS(fsys fs.FS) (ignore func(path string) bool) {
	code, err := fs.ReadFile(fsys, ".gitignore")
	if err != nil {
		return defaultIgnore
	}
	lines := strings.Split(string(code), "\n")
	lines = append(lines, alwaysIgnore...)
	ignorer := gitignore.CompileIgnoreLines(lines...)
	return ignorer.MatchesPath
}

// from
func from(dir string) (ignore func(path string) bool) {
	code, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		return defaultIgnore
	}
	lines := strings.Split(string(code), "\n")
	lines = append(lines, alwaysIgnore...)
	ignorer := gitignore.CompileIgnoreLines(lines...)
	return ignorer.MatchesPath
}
