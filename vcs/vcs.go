/*
 * Copyright 2018 Florent Biville (@fbiville)
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

package vcs

import (
	"fmt"
	. "github.com/fbiville/headache/helper"
	"os/exec"
	"strings"
)



type Vcs interface {
	Status(args ...string) (string, error)
	Diff(args ...string) (string, error)
	LatestRevision(file string) (string, error)
	Log(args ...string) (string, error)
	ShowContentAtRevision(path string, revision string) (string, error)
	Root() (string, error)
}

type Git struct{}
func (*Git) Status(args ...string) (string, error) {
	return git(PrependString("status", args)...)
}
func (*Git) Diff(args ...string) (string, error) {
	return git(PrependString("diff", args)...)
}
func (g *Git) LatestRevision(file string) (string, error) {
	result, err := g.Log("-1", `--format=%H`, "--", file)
	if err != nil {
		return "", err
	}
	return strings.Trim(result, "\n"), nil
}
func (*Git) Log(args ...string) (string, error) {
	return git(PrependString("log", args)...)
}
func (*Git) ShowContentAtRevision(path string, revision string) (string, error) {
	if revision == "" {
		return "", nil
	}
	fullRevision, err := revParse(revision)
	if err != nil {
		return "", err
	}
	fullRevision = strings.Trim(fullRevision, "\n")
	return git("cat-file", "-p", fmt.Sprintf("%s:%s", fullRevision, path))
}
func (*Git) Root() (string, error) {
	result, err := git("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.Trim(result, "\n"), nil
}

func revParse(revision string) (string, error) {
	return git("rev-parse", revision)
}

func git(args ...string) (string, error) {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
