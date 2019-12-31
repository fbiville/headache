/*
 * Copyright 2018-2019 Florent Biville (@fbiville)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
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
	"strconv"
	. "strings"
	"time"
)

type VersioningClient interface {
	GetChanges(revision string) ([]FileChange, error)
	AddMetadata(changes []FileChange, clock Clock) ([]FileChange, error)
	GetClient() Vcs
}

type Client struct {
	Vcs Vcs
}

type FileChange struct {
	Path            string
	CreationYear    int
	LastEditionYear int
}

type FileHistory struct {
	CreationYear    int
	LastEditionYear int
}

const (
	duplicatedRenamedContents = "R100"
	duplicatedCopiedContents  = "C100"
)

func (client *Client) GetChanges(revision string) ([]FileChange, error) {
	vcs := client.Vcs
	committedChanges, err := GetCommittedChanges(vcs, revision)
	if err != nil {
		return nil, err
	}
	uncommittedChanges, err := GetUncommittedChanges(vcs)
	if err != nil {
		return nil, err
	}
	return merge(committedChanges, uncommittedChanges), nil
}

func (client *Client) AddMetadata(changes []FileChange, clock Clock) ([]FileChange, error) {
	for i, change := range changes {
		history, err := GetFileHistory(client.Vcs, change.Path, clock)
		if err != nil {
			return nil, err
		}
		change.CreationYear = history.CreationYear
		change.LastEditionYear = history.LastEditionYear
		changes[i] = change
	}
	return changes, nil
}

func (client *Client) GetClient() Vcs {
	return client.Vcs
}

func GetCommittedChanges(vcs Vcs, revision string) ([]FileChange, error) {
	revisions := fmt.Sprintf("%s..HEAD", revision)
	output, err := vcs.Diff("--name-status", revisions)
	if err != nil {
		return nil, err
	}
	result := make([]FileChange, 0)
	for _, line := range Split(output, "\n") {
		if line == "" {
			continue
		}
		statusName := SplitN(line, "\t", 2)
		status := Trim(statusName[0], " ")
		switch {
		case status == "D":
			// ignore
		case HasPrefix(status, "R"):
			statusName := SplitN(line, "\t", 3)
			result = append(result, FileChange{
				Path: Trim(statusName[2], " "),
			})
		default:
			result = append(result, FileChange{
				Path: Trim(statusName[1], " "),
			})
		}
	}
	return result, nil
}

func GetUncommittedChanges(vcs Vcs) ([]FileChange, error) {
	output, err := vcs.Status("--porcelain")
	if err != nil {
		return nil, err
	}
	result := make([]FileChange, 0)
	if output == "" {
		return result, nil
	}
	for _, line := range Split(output, "\n") {
		if line == "" {
			continue
		}
		statusName := SplitN(Trim(line, " "), " ", 2)
		statuses := Trim(statusName[0], " ")
		if Index(statuses, "D") != -1 {
			continue
		}
		result = append(result, FileChange{
			Path: Trim(statusName[1], " "),
		})
	}
	return result, nil
}

func GetFileHistory(vcs Vcs, file string, clock Clock) (*FileHistory, error) {
	output, err := vcs.Log("--follow", "--name-status", "--format=%at", "--", file)
	if err != nil {
		return nil, err
	}
	timestamps, err := getCommitTimestamps(file, output)
	if err != nil {
		return nil, err
	}
	defaultYear := clock.Now().Year()
	history := FileHistory{
		CreationYear:    defaultYear,
		LastEditionYear: defaultYear,
	}

	if len(timestamps) > 0 {
		minTimestamp := timestamps[len(timestamps)-1]
		maxTimestamp := timestamps[0]
		history.CreationYear = time.Unix(minTimestamp, 0).Year()
		history.LastEditionYear = time.Unix(maxTimestamp, 0).Year()
	}

	return &history, nil
}

func getCommitTimestamps(file string, log string) ([]int64, error) {
	var result []int64
	lines := Split(Replace(log, "\n\n", "\n", -1), "\n")
	lines = lines[0 : len(lines)-1]
	for i := 1; i < len(lines); i += 2 {
		line := lines[i]
		nameStatus := Split(line, "\t")[0]
		if nameStatus == duplicatedRenamedContents || nameStatus == duplicatedCopiedContents {
			continue
		}
		timestamp, err := strconv.ParseInt(lines[i-1], 10, 64)
		if err != nil {
			errorMsg := "could not parse timestamp (line %d) of file %q history. Full commit log below\n%s"
			return nil, fmt.Errorf(errorMsg, i, file, log)
		}
		result = append(result, timestamp)
	}
	return result, nil
}

func merge(changes []FileChange, changes2 []FileChange) []FileChange {
	set := make(map[FileChange]struct{}, len(changes))
	for _, change := range changes {
		set[change] = struct{}{}
	}

	for _, change := range changes2 {
		if _, ok := set[change]; !ok {
			set[change] = struct{}{}
		}
	}
	return keys(set)
}

func keys(set map[FileChange]struct{}) []FileChange {
	i := 0
	result := make([]FileChange, len(set))
	for key := range set {
		result[i] = key
		i++
	}
	return result
}
