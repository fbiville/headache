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
	"encoding/json"
	"flag"
	"fmt"
	. "github.com/fbiville/header/core"
	"io/ioutil"
	"os"
)

func main() {
	configFile := flag.String("configuration", "license.json", "Path to configuration file")
	dryRun := flag.Bool("dry-run", false, "Dumps the execution to a file instead of altering the sources")
	dumpFile := flag.String("dump-file", "", "Path to the dry-run dump")

	flag.Parse()

	if *dumpFile != "" && *dryRun {
		panic("cannot simultaneously use --dump-file and --dry-run")
	}
	executionMode := RegularRunMode
	if *dryRun {
		executionMode = DryRunMode
	} else if *dumpFile != "" {
		executionMode = RunFromFilesMode
	}

	rawConfiguration := readConfiguration(configFile)
	configuration, err := ParseConfiguration(rawConfiguration, executionMode, dumpFile)
	if err != nil {
		panic(err)
	}
	if executionMode == DryRunMode {
		file, err := DryRun(configuration)
		displayDryRunResult(file, err)
	} else {
		Run(configuration)
	}
}

func readConfiguration(configFile *string) Configuration {
	flag.Parse()
	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}
	result := Configuration{
		VcsImplementation: "git",
		VcsRemote:         "origin",
		VcsBranch:         "master",
	}
	json.Unmarshal(file, &result)
	return result
}

func displayDryRunResult(file string, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occurred during the execution, see below:")
		panic(err)
	} else {
		fmt.Printf("See dry-run result in file printed below:\n%s\n", file)
	}
}
