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

	// dependency graph - begin
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
	configurationResolver := &ConfigurationResolver{
		SystemConfiguration: systemConfig,
		ExecutionTracker:    executionTracker,
		PathMatcher:         &fs.ZglobPathMatcher{},
	}
	headache := &Headache{Fs: fileSystem}
	// dependency graph - end

	configFile, configuration := loadConfiguration(configLoader, configurationResolver)
	if len(configuration.Files) > 0 {
		headache.Run(configuration)
		if err := executionTracker.TrackExecution(configFile); err != nil {
			log.Printf("headache warning, could not save current execution, see below for details\n\t%v\n", err)
		}
	} else {
		log.Print("No files to process")
	}

	log.Print("Done!")
}

func loadConfiguration(configLoader *ConfigurationLoader, configResolver *ConfigurationResolver) (*string, *ChangeSet) {
	configFile := flag.String("configuration", "headache.json", "Path to configuration file")
	flag.Parse()

	userConfiguration, err := configLoader.ReadConfiguration(configFile)
	if err != nil {
		log.Fatalf("headache configuration error, cannot load\n\t%v\n", err)
	}
	configuration, err := configResolver.ResolveEagerly(userConfiguration)
	if err != nil {
		log.Fatalf("headache configuration error, cannot parse\n\t%v\n", err)
	}
	return configFile, configuration
}
