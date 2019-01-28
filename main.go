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

package main

import (
	"flag"
	. "github.com/fbiville/headache/core"
	"github.com/fbiville/headache/fs"
	"log"
)

func main() {
	log.Print("Starting...")

	// poor man's dependency graph
	systemConfig := DefaultSystemConfiguration()
	fileSystem := systemConfig.FileSystem
	configLoader := &ConfigurationLoader{
		Reader: fileSystem.FileReader,
	}
	executionTracker := &ExecutionVcsTracker{
		Versioning:   systemConfig.VersioningClient.GetClient(),
		FileSystem:   fileSystem,
		Clock:        systemConfig.Clock,
		ConfigLoader: configLoader,
	}
	matcher := &fs.ZglobPathMatcher{}

	configFile := parseFlags()

	userConfiguration, err := configLoader.ReadConfiguration(configFile)
	if err != nil {
		log.Fatalf("headache configuration error, cannot load\n\t%v\n", err)
	}

	configuration, err := ParseConfiguration(userConfiguration, systemConfig, executionTracker, matcher)
	if err != nil {
		log.Fatalf("headache configuration error, cannot parse\n\t%v\n", err)
	}

	if len(configuration.Files) > 0 {
		Run(configuration, fileSystem)
		trackRun(configFile, executionTracker)
	}
	log.Print("Done!")
}

func parseFlags() *string {
	configFile := flag.String("configuration", "headache.json", "Path to configuration file")
	flag.Parse()
	return configFile
}

func trackRun(configFile *string, tracker ExecutionTracker) {
	err := tracker.TrackExecution(configFile)
	if err != nil {
		log.Printf("headache warning, could not save current execution, see below for details\n\t%v\n", err)
	}
}
