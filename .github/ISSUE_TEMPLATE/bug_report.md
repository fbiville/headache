---
name: Bug report
about: Create a report to help us improve headache
title: ''
labels: bug
assignees: fbiville

---

**Environment**

 - headache version: 
 - Git repository URL (if possible): 
 - contents of `.headache-run`: 
```
.headache-run contents to insert here
```
 - latest revision of `.headache-run` (`git --no-pager log --format='%H' -1 -- .headache-run`): `edit this`
 - contents of JSON configuration (default name: `headache.json`):
```
JSON configuration contents to insert here
```
- contents of the license header template:
```
header template contents to insert here
```

**Protips**

_the following only applies to duplicated or malformed header issues_

`headache` can be quite sensitive to small differences between the header template and the actual headers in source files, especially at the first run!

Here a few things to check before reporting an issue:

 - [ ] the configuration JSON file and header template have not been changed since `.headache-run` last revision
 - [ ] there are no differences between the affected source file's header and the header in configured template file

**Describe the bug**
A clear and concise description of what the bug is.
If it concerns duplicated/malformed headers, please include the full contents of at least one affected source file:
```
source file contents to insert here
```

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'

**Expected behavior**
A clear and concise description of what you expected to happen.

**Additional context**
Add any other context about the problem here.
