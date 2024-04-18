## [0.2.0] - 2023-12-04

### Added

#### 🚀 Features

- Extend the Python interpreter. See [How to extend the Python interpreter](docs/HOWTO-EXTEND.md).
- Concurrency primitives:
  - `Py.GoEnterPython`, `Py.GoLeavePython`: see [Accessing the interpreter](docs/HOWTO-EXTEND.md#accessing-the-interpreter).
  - `Py.GoNewFlow`: see [Concurrency](docs/HOWTO-EXTEND.md#concurrency).
  - `Py.GoClosureBuddy`: see [GoClosure and its buddy](docs/HOWTO-EXTEND.md#goclosure-and-its-buddy).
- Execute CI tests also with Python debug versions so to catch more bugs.
- Augment tests with property-based testing.
- Support for Python 3.12.

#### 🐛 Bug fixes

- Goroutines are now pinned to their OS threads, Python uses thread-local
  storage and does not tolerate goroutines migration to different threads.
- Do not reuse threads that accessed the Python interpreter, Go runtime
  needs absolute control of its threads and Python could modify them in
  unexpected ways.

### Modified

- All the Go APIs now start with the `Go` prefix, ex. `py.Args` and `py.KwArgs`
  are now named `py.GoArgs` and `py.GoKwArgs`.
- `Py.Object_Length` now returns error by a proper `error` value
  instead of just `-1`.
- `*py.GoError` (not `py.GoError`) now implements the `error` interface,
  incorrect type casting is now detected during compilation.
- Python context for the embedded interpreter is now created by
  `py.GoEmbed`, `py.Py{}` is an invalid and unusable context. See
  [Initialization and finalization](docs/HOWTO-EMBED.md#initialization-and-finalization)
  for the updates.

### Removed

- `Py.Tuple_New` and `Py.Tuple_SetItem` ([#16](https://gitlab.com/pygolo/py/-/issues/16)),
  use `Py.Tuple_Pack` instead.
- Support for Go 1.9, it does not prevent thread reuse when a
  pinned goroutine returns. Scientific Linux 7 is dropped from the test matrix.

## [0.1.1] - 2023-08-28

#### 🐛 Bug fixes

- Fix build issue with Go 1.21

## [0.1.0] - 2023-07-04

💥 First release ever! 💥

This release focused on readying basic life support and solid ground for growth.

### Added

#### 🚀 Features

- Embed the Python interpreter
- Convert basic types to/from Python
- Handle Python exceptions as regular Go errors
- Import modules
- Call functions and objects
- Build with pyenv provided interpreters
- Run in venv environments

#### 📚 Documentation

- [Contributing](CONTRIBUTING.md)
- [How to embed](docs/HOWTO-EMBED.md)
- [Advanced topics](docs/ADVANCED-TOPICS.md)

### It Could Work!

<img src="http://www.frankensteinjunior.it/download/foto/1/big/FJ_015.jpg" alt="Gene Wilder exclaiming 'It Could Work!'" width=256 height=144>
