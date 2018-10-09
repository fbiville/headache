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

package versioning

import (
	"fmt"
	. "github.com/fbiville/headache/helper"
	"os/exec"
	"strings"
)

type Git struct{}

func (Git) Status(args []string) (string, error) {
	return runGit(PrependString("status", args))
}

func (Git) Diff(args []string) (string, error) {
	return runGit(PrependString("diff", args))
}

func (Git) Log(args []string) (string, error) {
	return runGit(PrependString("log", args))
}

func (Git) ShowContentAtRevision(path string, revision string) (string, error) {
	fullRevision, err := revParse(revision)
	if err != nil {
		return "", err
	}
	fullRevision = strings.Trim(fullRevision, "\n")
	return runGit([]string{"cat-file", "-p", fmt.Sprintf("%s:%s", fullRevision, path)})
}

func revParse(revision string) (string, error) {
	return runGit([]string{"rev-parse", revision})
}

func runGit(args []string) (string, error) {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
