# Golang header management

`headache` manages license headers of Go files.

## Example

By default, `headache` looks for a file named `license.json` in the current directory:

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
| `vcs`            | string                  | Versioning system, only `"git"` is supported for now and is the default value.  |
| `vcsRemote`      | string                  | Remote of the current branch, defaults to `"origin"`.  |
| `vcsBranch`      | string                  | Current branch, defaults to `"master"`.                |



## Custom configuration

By default, a file named `license.json` must be present in the current directory.

Alternatively, the configuration file can be explicitly provided:
```shell
 $ go get -u github.com/fbiville/headache
 $ $(GOBIN)/headache --configuration /path/to/configuration.json
```

All the examples below support that option.

## First run

Normal dry-run/run executions detect only recent changes, based on the versioning
configuration.

When you need to run the first time, you will have to run:
```shell
 $ go get -u github.com/fbiville/headache
 $ $(GOBIN)/headache --dry-run --init
```

And follow the steps just below to read, possibly edit and reinject
the generated dump file.

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

The dump file follows this structure:
```text
file:FILENAME
---
multi-line
colored diff 1
---
file:OTHERFILE
---
multi-line
colored diff 2
---
```

### View the colorized diff summary

If you want to see the dump contents, you can run:
```shell
 $ less -r /path/to/headache-dry-runXXX
```


### List the files

If you want to get a list of the Go files possibly* changed the future
execution, you can run something like:

```
 $ ./headache --dry-run | tail -n 1 | xargs cat | grep '^file:.*\.go' | sed s/file:// | sort
```

_\* the execution may result in no changes at all, or rather, the new written content
will be the same as the previous one_

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

`headache` currently does **not** support text changes other than:

 * parameter value updates
 * comment style changes
