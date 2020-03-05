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

package main

import (
	"flag"
	"fmt"
	. "github.com/fbiville/headache/internal/pkg/core"
	"github.com/fbiville/headache/internal/pkg/fs"
	"log"
	"os"
)

func main() {
	log.Print("Starting...")

	// dependency graph - begin
	environment := DefaultEnvironment()
	fileSystem := environment.FileSystem
	configLoader := &ConfigurationFileLoader{
		Reader:         fileSystem.FileReader,
		SchemaLocation: environment.SchemaLocation,
		SchemaLoader:   &JsonSchemaFileLoader{},
	}
	executionTracker := &ExecutionVcsTracker{
		Versioning:   environment.VersioningClient.GetClient(),
		FileSystem:   fileSystem,
		Clock:        environment.Clock,
		ConfigLoader: configLoader,
	}
	configurationResolver := &ConfigurationResolver{
		Environment:      environment,
		ExecutionTracker: executionTracker,
		PathMatcher:      &fs.ZglobPathMatcher{},
	}
	headache := &Headache{Fs: fileSystem}
	// dependency graph - end

	execute(headache, configLoader, configurationResolver, executionTracker)
}

func execute(headache *Headache, configLoader *ConfigurationFileLoader, configurationResolver *ConfigurationResolver, executionTracker *ExecutionVcsTracker) {
	configFile := flag.String("configuration", "headache.json", "Path to configuration file")
	checkMode := flag.Bool("check", false, "Checks if headers are up-to-date")
	flag.Parse()

	exitCode := 0
	configuration := loadConfiguration(configFile, configLoader, configurationResolver)
	if len(configuration.Files) > 0 {
		if *checkMode {
			if diff := headache.DryRun(configuration); diff != "" {
				_, _ = fmt.Fprintf(os.Stderr, "Headers are not up-to-date! See details in file: TODO")
				// TODO: write in tmp file
				_, _ = fmt.Fprintf(os.Stderr, diff)
				exitCode = 1
			} else {
				fmt.Print("Check successful!")
			}
		} else {
			headache.Run(configuration)
			if err := executionTracker.TrackExecution(configFile); err != nil {
				log.Printf("headache warning, could not save current execution, see below for details\n\t%v\n", err)
			}
		}
	} else {
		log.Print("No files to process")
	}

	log.Print("Done!")
	os.Exit(exitCode)
}

func loadConfiguration(configFile *string, configLoader *ConfigurationFileLoader, configResolver *ConfigurationResolver) *ChangeSet {
	userConfiguration, err := configLoader.ValidateAndLoad(*configFile)
	if err != nil {
		log.Fatalf("headache configuration error, cannot load\n\t%v\n", err)
	}
	configuration, err := configResolver.ResolveEagerly(userConfiguration)
	if err != nil {
		log.Fatalf("headache configuration error, cannot parse\n\t%v\n", err)
	}
	return configuration
}
