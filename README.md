# Golang header management

`headache` manages license headers of Go files.

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
| `style`          | string                  | One of: SlashStar (`/* ... */`), SlashSlash (`// ...`) |
| `includes`       | array of strings        | File globs to include (`*` and `**` are supported)     |
| `excludes`       | array of strings        | File globs to exclude (`*` and `**` are supported)     |
| `data`           | map of string to string | Key-value pairs, matching the parameters used in `headerFile`.<br>Please note that `{{.Year}}` is a reserved parameter and will automatically be computed based on the files versioning information.  |



## Custom configuration

By default, a file named `headache.json` must be present in the current directory.

Alternatively, the configuration file can be explicitly provided:
```shell
 $ go get -u github.com/fbiville/headache
 $ $(GOBIN)/headache --configuration /path/to/configuration.json
```

All the examples below support that option.

## Dry run

All you have to do then is to simulate the run:
```shell
 $ go get -u github.com/fbiville/headache
 $ $(GOBIN)/headache --dry-run
```

The command will output the file in which the actual diff summary is appended to.

For instance:
```
See dry-run result in file printed below:
/path/to/headache-dry-runXXX
```

The dump file aggregates the diff for each file that would be modified by the execution.

### List the files

If you want to get a list of the Go files possibly changed the future
execution, you can run something like:

```
 $ ./headache --dry-run | tail -n 1 | xargs cat | grep '^file:.*\.go' | sed s/file:// | sort
```


### Exclude files

Copyright years should only be updated after a
significant change is made (read this
[Stack Overflow post](https://stackoverflow.com/questions/2390230/do-copyright-dates-need-to-be-updated)
for more information).

To exclude files from being unnecessarily updated, locate the corresponding line, prefixed by `file:`,
followed by the file name and replace `file:` by `xfile:`.

Then, the modified dump file can be fed back to `headache`, as described just below.

## Run

### From dry-run dump
Once you have successfully run `headache --dry-run` and
possibly edited the dump file (see above to see how), all you have to do then is to run:

```shell
 $ go get -u github.com/fbiville/headache --dump-file /path/to/headache-dry-runXXX
 $ $(GOBIN)/headache
```

This will update only the files for which names are prefixed by `file:`.

### Direct run

All you have to do then is to run:
```shell
 $ go get -u github.com/fbiville/headache
 $ $(GOBIN)/headache
```


## Unsupported

`headache` currently does **not** support text changes **other than**:

 * parameter value updates
 * comment style changes
