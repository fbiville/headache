# Header management

`headache` manages license headers.
It is biased towards Golang but should work on any other language (provided a compatible code style is implemented).

## Build status

[![Build Status](https://travis-ci.org/fbiville/headache.svg?branch=master)](https://travis-ci.org/fbiville/headache)

## Example

By default, `headache` looks for a file named `headache.json` in the current directory:

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

`license-header.txt`:
```
Copyright {{.Year}} {{.Owner}}

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```


## Settings

Setting            | Type                    | Definition                                             |
| ---------------- |:----------------------: | -----------------------------------------------------: |
| `headerFile`     | string                  | Path to the parameterized license header. Parameters are referenced with the following syntax: {{.PARAMETER-NAME}}               |
| `style`          | string                  | One of: `SlashStar` (`/* ... */`), `SlashSlash` (`// ...`) |
| `includes`       | array of strings        | File globs to include (`*` and `**` are supported)     |
| `excludes`       | array of strings        | File globs to exclude (`*` and `**` are supported)     |
| `data`           | map of string to string | Key-value pairs, matching the parameters used in `headerFile`.<br>Please note that `{{.Year}}` is a reserved parameter and will automatically be computed based on the files versioning information.  |


## Run

By default, a file named `headache.json` must be present in the current directory.

All you have to do then is:
```shell
 $ go get -u github.com/fbiville/headache
 $ $(GOBIN)/headache
```

As a result, source files will be changed and `.headache-run` will be generated to keep track of `headache` last execution.
This file must be versioned as well.

## Run with custom configuration

Alternatively, the configuration file can be explicitly provided:
```shell
 $ go get -u github.com/fbiville/headache
 $ $(GOBIN)/headache --configuration /path/to/configuration.json
```

