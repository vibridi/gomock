
v3.7.0 / 2026-04-17
==================

  * Add pkgs option

v3.6.1 / 2026-04-17
==================

  * Fix: readd notice when printing output to file

v3.6.0 / 2026-04-17
==================

  * Add writer/file unit test
  * Write to output file without truncating
  * Refactor writer package to better encapsulate types
  * Rename ToFuncDef to AppendFuncDef directly modifying receiver
  * Add basic method and type docs

v3.5.0 / 2026-04-16
==================

  * Increase unit test coverage to 80%
  * Merge pull request #12 from vibridi/generics
  * Update README file
  * Report unit test coverage
  * Modernize github workflow
  * Support generics in struct mode
  * Remove helper package
  * Support mocking generic interfaces
  * Add TypeParamFields to MockData
  * Add FuncDef documentation
  * Update Go version to 1.26.2

v3.4.0 / 2026-04-15
===================

  * Add -p option to README
  * Add -p option to merge package name and interface name

v3.3.0 / 2024-05-30
===================

  * Bump version to v3.3.0
  * Add -d option to disambiguate withFunc identifiers
  * Bump version to v3.2.1
  * Fix module path when setting VERSION and GOVERSION via ldflags in Makefile
  * Add goreport.com badge

v3.3.0 / 2024-05-30
==================

  * Add -d option to disambiguate withFunc identifiers

v3.2.1 / 2024-05-03
==================

  * Fix module path when setting VERSION and GOVERSION via ldflags in Makefile
  * Add goreport.com badge

v3.2.0 / 2024-02-23
==================

  * Update README
  * Remove install target, locally it's possible to install just with 'go install'
  * Add --utype flag to map named types to underlying types via CLI
  * Upgrade required Go version to 1.22.0

v3.1.2 / 2023-09-21
==================

  * Upgrade Go version to 1.21.1 (fixes #10)

v3.1.1 / 2023-07-24
==================

  * Add status notice to README
  * Upgrade libraries

v3.0.0 / 2020-06-07
==================

  * Create go.yml
  * Update Readme with breaking changes from v2
  * Add 'name' option to override the name of the interface
  * Invert '-q' flag. Qualify with package name by default, and require specific flag to opt out instead. Assumption is that you don't need mocks within their own package.
  * Improve usage tip, title and flag help messages
  * Allow short options (one dash followed by N options)
  * Remove implemented unnamed params TODO from Readme

v2.0.0 / 2020-04-25
==================

  * Change writer return type to []byte and autoformat
  * Several changes: - refactor TemplateData into templates package - add Go template for struct style output - instantiate writer with New to avoid passing around the options object
  * Support --struct-style option
  * Upgrade urfave/cli to v2

v1.1.0 / 2020-03-16
==================

  * Outputting unnamed arguments is opt-in
  * Pass flags to Write with an options struct

v1.0.0 / 2020-02-20
==================

  * Support composed interfaces
  * Properly print messages to stderr

v0.2.3 / 2020-02-14
==================

  * Write default and with* functions with unnamed arguments

v0.2.2 / 2019-11-22
==================

  * Fix message of -i option and update readme

v0.2.1 / 2019-11-22
==================

  * Fix help message of '-x' flag

v0.2.0 / 2019-10-08
==================

  * Introduce changelog
  * Qualify exported identifiers from same package
  * Migrate to gomodules
