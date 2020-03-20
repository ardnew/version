# version
### Go module to easily embed semantic versioning compliance with change history

Note that this module is _not designed to parse and compare semantic versions_. There is [much better software](https://github.com/Masterminds/semver) designed for that task.

The goal is to provide a consistent interface for defining and reporting the version and changes of a Go package.

## Usage
Simply defining the global `ChangeLog` variable will conveniently record work history and doubles as package version definition. No function calls required. Similar to `godoc`, it's as simple to use as good comments.

Alternatively, you can just call `Set()` to set package version.

## Features
- [x] Compliant with [Semantic Versioning](https://semver.org/) (2.0.0)
- [x] Can parse and generate changelog for release notes
- [ ] Can automatically integrate with `flag` package (e.g., `-version`, `-changes`, and other command-line flags)

