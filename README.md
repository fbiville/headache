# Golang header management

`header` manages license headers for you.

## Settings

Setting            | Type                    | Definition                                             |
| ---------------- |:----------------------: | -----------------------------------------------------: |
| `headerFile`     | string                  | Path to the parameterized license header               |
| `style`          | string                  | One of: SlashStar (`/* ... */`), SlashSlash (`// ...`) |
| `includes`       | array of strings        | File globs to include (`*` and `**` are supported)     |
| `excludes`       | array of strings        | File globs to exclude (`*` and `**` are supported)     |
| `data`           | map of string to string | Key-value pair, matching the parameters in headerFile  |


## Example

By default, `header` looks for a file named `license.json` in the current directory:

```json
{
  "headerFile": "./license-header.txt",
  "style": "SlashStar",
  "includes": ["**/*.go"],
  "excludes": ["vendor/**/*"],
  "data": {
    "Year": "2018",
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

All you have to do then is to run:
```shell
 $ go get -u github.com/fbiville/header
 $ $(GOBIN)/header
```

Alternatively, the configuration file can be explicitly provided:
```shell
 $ go get -u github.com/fbiville/header
 $ $(GOBIN)/header /path/to/configuration.json
```