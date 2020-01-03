# Header management

`headache` is an opinionated license header updater (see the "Approach" section below).
It allows to consistently insert and update license headers in source files or even change licenses completely!

## Build status

[![Build Status](https://github.com/fbiville/headache/workflows/CI/badge.svg)](https://github.com/fbiville/headache/actions)

## Quick start

By default, `headache` looks for a configuration file named `headache.json` in the directory in which it is invoked:

```json
{
  "headerFile": "./license-header.txt",
  "style": "SlashStar",
  "includes": ["**/*.go"],
  "excludes": ["vendor/**/*"],
  "data": {
    "Owner": "The original author or authors"
  }
}
```

`license-header.txt` (note the absence of `YearRange` parameter in the configuration file):
```
Copyright {{.YearRange}} {{.Owner}}

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

### Run

All you have to do then is:
```shell
 $ cd $(mktemp -d) && go get -u github.com/fbiville/headache/cmd/headache && cd -
 $ ${GOBIN:-`go env GOPATH`/bin}/headache
```

As a result, source files will be changed and `.headache-run` will be generated to keep track of `headache` last execution.
This file must be versioned along with the source file changes.

### Run with custom configuration

Alternatively, the configuration file can be explicitly provided:
```shell
 $ go get -u github.com/fbiville/headache
 $ ${GOBIN:-`go env GOPATH`/bin}/headache --configuration /path/to/configuration.json
```

## Reference documentation

### Approach

`headache` approach to copyright is well explained in [this stackoverflow answer](https://stackoverflow.com/a/2391555/277128), 
read it first!

Now that you read this, here are two important points:

 - Copyright years have to be updated when a significant change occurs.
 
There is, to the author knowledge, no automatic solution to distinguish a trivial change from a significant one.

Based on this premise, `headache` will process all files matching the configuration and that have been changed since its last execution.
`headache` will then compute the copyright year, file by file, from their available versioning information (typically by retrieving the relevant dates from Git commits).

**It is up to the project maintainer to discard the generated changes if they are not relevant.**

 - > The [first] date on the [copyright] notice establishes how far back the claim is made.
 
This claim could predate any commit associated to the file (imagine a file copied from project to project for years).

`headache` will never overwrite the start date of the copyright year if it finds one, if and only if that date occurs earlier than the first commit date of the file.

### Configuration

`headache` relies on the emerging [JSON Schema standard](https://json-schema.org/) to validate its configuration.
`headache` schema is defined [here](https://fbiville.github.io/headache/schema.json).

In layman's terms, here are all the possible settings:

Setting            | Type                    | Definition                                             |
| ---------------- |:----------------------: | -----------------------------------------------------: |
| `headerFile`     | string                  | **[required]** Path to the parameterized license header. Parameters are referenced with the following syntax: {{.PARAMETER-NAME}}               |
| `style`          | string                  | **[required]** See all the possible names [here](https://fbiville.github.io/headache/schema.json). The lookup is case-insensitive. |
| `includes`       | array of strings        | **[required, min size=1]** File globs to include (`*` and `**` are supported)     |
| `excludes`       | array of strings        | File globs to exclude (`*` and `**` are supported)     |
| `data`           | map of string to string | Key-value pairs, matching the parameters used in `headerFile` except for the reserved parameters (see below section).


#### Reserved parameters

 - `{{.YearRange}}` (formerly `{{.Year}}`) is automatically substituted with either:
     - a single year if both years in the range are the same
     - a year range with the earliest commit's year (or earlier) and latest commit's year
 - `{{.StartYear}}` is substituted with the earliest commit's year (or earlier)
 - `{{.EndYear}}` is substituted with the latest commit's year
 
As explained earlier, if a file specifies a start date in its header that is earlier than any commit's year, then that
date is preserved.

If you want to avoid copyright dates like `2019-2019`, then rely on `{{.YearRange}}` instead of `{{.StartYear}}-{{.EndYear}}`.
If you need something like `2018-present`, then use `{{.StartYear}}-present` instead.
