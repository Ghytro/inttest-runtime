In this page you should find all you need to start. If you think something
is missing, unclear or wrong please reach out.

[[_TOC_]]

### Issues

Ideally everything in the project should be tracked by a Gitlab issue,
it's a matter of transparency towards current and future contributors.

There is no special process in place, the issues are just tagged with any
of the following labels:

- **Bug**: issues reporting misbehaving code, incorrect documentation,
  or anything else that does not work as intended.

- **Good first issue**: these issues are a good starting point for
  newcomers, they do not require any (deep) knowledge of the project but
  allow gaining some along the way.

- **Help wanted**: used for issues that require expertise, tools, or
  something else that is not available within the project. It's likely
  that without that specific help the issue will remain open or unsolved.

- **More info**: marks issues that do not contain all the
  information needed to proceed, details are missing, further discussion
  is needed, etc.

### How to ask for help

At times it's not clear what to do, while it's useful to invest some time
on your own it's better to ask for help well before the frustration gets
unhealthy.

Join the [#pygolo/py](https://gitter.im/pygolo/py) chat, the
[r/pygolo](https://www.reddit.com/r/pygolo) Reddit community, observe our
[Code of conduct](CODE-OF-CONDUCT.md), feel free to ask. We'll do our best
to support you.

### Environment set up

Your OS is either Linux, macOS or FreeBSD. You need the following tools
installed:

- Go compiler (>= 1.10)
- Python (>= 3.6) and its development libraries
  (ex. `python3`, `python3-dev` or `python3-devel`)
- C compiler (`gcc` on Linux, pre-installed on FreeBSD)
- GNU Make (`make` on Linux, `gmake` on FreeBSD)
- `pkg-config` or its drop-in replacement `pkgconf`
- pytest framework (`python3-pytest` on Linux, `py39-pytest` on FreeBSD)
- pytest-benchmark (`python3-pytest-benchmark` on Linux, `py39-pytest-benchmark`
  on FreeBSD)
- Hypothesis library (`python3-hypothesis` on Linux, `py39-hypothesis` on FreeBSD)
- git

Notes for macOS:

- you need Command Line Tools (CLT) for Xcode; it provides the C compiler,
  `make`, `git`, and `python3`
- `go` and `pkgconf` are available on both [Homebrew](https://brew.sh)
  and [MacPorts](https://www.macports.org)
- if you want to build with the CLT provided Python interpreter you also
  need to install Xcode; versions of Xcode and CLT for Xcode need to match
- building with Python installations from Homebrew and MacPorts is also
  fine and supported
- `pytest`, `pytest-benchmark`, and `hypothesis` are available from MacPorts
  and other Python package managers such as [pip](https://pip.pypa.io/en/stable)
  and [poetry](https://python-poetry.org/docs)

#### Code checkout

Now clone the source tree and walk into it:

```shell
$ git clone https://gitlab.com/pygolo/py.git
$ cd py
```

Try to execute the tests:

```shell
$ make
go test -tags "py3.9" ./test
ok  	gitlab.com/pygolo/py/test	1.104s
go build -tags "py3.9 py_ext" -buildmode=c-shared -o test/ext.so ./test/ext
python3 -m pytest -q test
...
3 passed in 0.51s
```

If you get a similar output congratulations, your environment is up an
running! 🎉

Try also the examples (the output depends on your system):

```shell
$ make examples RUN=1
go build -tags "py3.9" -o examples/hello examples/hello.go
go build -tags "py3.9" -o examples/info examples/info.go
### hello ###
Hello, World!
Hi, cavokz!
### info ###
Python:
  version: "3.9"
  platform: "macosx-10.9-universal2"
  platlib: "/Library/Python/3.9/site-packages"
```

Check what `hello` is linked to:

```shell
# on Linux
$ ldd examples/hello
	linux-vdso.so.1 (0x00007ffe79690000)
	libpython3.11.so.1.0 => /lib/x86_64-linux-gnu/libpython3.11.so.1.0 (0x00007f2d36800000)
	libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f2d3661f000)
	libm.so.6 => /lib/x86_64-linux-gnu/libm.so.6 (0x00007f2d36fa7000)
	libz.so.1 => /lib/x86_64-linux-gnu/libz.so.1 (0x00007f2d36600000)
	libexpat.so.1 => /lib/x86_64-linux-gnu/libexpat.so.1 (0x00007f2d365d5000)
	/lib64/ld-linux-x86-64.so.2 (0x00007f2d37096000)
```

```shell
# on macOS
$ otool -L examples/hello
examples/hello:
	@rpath/Python3.framework/Versions/3.9/Python3 (compatibility version 3.9.0, current version 3.9.0)
	/usr/lib/libresolv.9.dylib (compatibility version 1.0.0, current version 1.0.0)
	/usr/lib/libSystem.B.dylib (compatibility version 1.0.0, current version 1319.100.3)
```

```shell
# on FreeBSD
$ ldd examples/hello
examples/hello:
	libpython3.9.so.1.0 => /usr/local/lib/libpython3.9.so.1.0 (0x8003bd000)
	libthr.so.3 => /lib/libthr.so.3 (0x800737000)
	libc.so.7 => /lib/libc.so.7 (0x800765000)
	libcrypt.so.5 => /lib/libcrypt.so.5 (0x800b6f000)
	libintl.so.8 => /usr/local/lib/libintl.so.8 (0x800b90000)
	libdl.so.1 => /usr/lib/libdl.so.1 (0x800b9e000)
	libutil.so.9 => /lib/libutil.so.9 (0x800ba2000)
	libm.so.5 => /lib/libm.so.5 (0x800bba000)
```

You are good to go. 🙌

#### Makefile targets

You already got it, `make` is part of the development flow. Let's see what
else you can do with it.

Most useful target is `test`, which is also the default; it runs all the
unit tests also executed by the CI/CD pipeline.

Then you have `lint` and `license-check`, they ensure all files have the
same formatting, code style, and license; `spell-check` searches typos in
the documentation. Best used in the pre-push git hook, they are executed
also in the CI/CD pipeline.

With `examples` you build the examples, if you add `RUN=1` they are also
executed.

`clean` removes the built executables, `mrproper` wipes also the Go caches.

`prereq` is maybe not too well defined, it's just used to install the
linting tools in the CI/CD pipeline. You can use it as well for the same
purpose but it's not really the one-stop for installing all the other
development dependencies as its name might suggest.

`pygolo-diags` gives you details about your environment and how `make`
adjusts itself. Try it with the `PYTHON` parameter (see below).

#### Makefile parameters

With `PYTHON=<some_python>` you execute the invoked target aiming at the
given Python. It can be an absolute path or just an executable's name.

For example `make examples RUN=1 PYTHON=python3.10` will use `python3.10`
instead of the default `python3` interpreter for building and running the
examples.  It's quite handy if you want to try alternative Python
installations (ex. pyenv) on your system.

Don't forget to clean the Go cache as explained [below](#go-build-cache).

With `V=1` you get increased output from the invoked target.

### Overview of the source tree

All the source code of this module is in the root directory except for the
tests, they are in the `test` directory.

Most of the Pygolo Go API has a 1:1 correspondence with the Python C API.
As example, the Go API function `Dict_New` matches C API function
`PyDict_New`; the Go API type `Object` matches the C API type `PyObject`.
Pygolo functions and types unrelated to the C API have prefix `Go`.

The Go function name prefix (ex. `Dict_`) determines the file name where
the function is implemented (ex. `dict.go`). Pygolo specific functions are
implemented in files named with prefix `go-` (ex. `go-util.go`).

In the `examples` directory you find small programs that use Pygolo,
`docs` and `images` directories are self-explanatory, `scripts` contains
auxiliary utilities used in the CI/CD pipeline.

### Code style

All the Go code is formatted with [gofmt](https://pkg.go.dev/cmd/gofmt)
(see [go fmt your code](https://go.dev/blog/gofmt)). We deliberately use
underscores (`_`) in the Go API names so to stay close to the naming in
the Python C API.

We have plenty of variables referring to Python objects, they are
generally named `o`. If more objects are referred in the same scope, the
main object (usually the first function parameter) is named `o` and all
the others have names starting with `o_`.

We are not scared by single-letter variable's names as long they help to
keep the code readable, otherwise more descriptive (but not too long!)
names are preferred. In case of doubts, your taste and `make lint` are
your best friends.

### Known issues

#### Go build cache

To save time the Go compiler caches the compilation of files and monitors
changes to their content. Because variations in the C environment are not
tracked, the Go compilation cache overlooks changes like in the version of
Python and its libraries.

If you are working with multiple versions of Python be sure to clean the
Go cache in between the versions switch with `go clean -cache -testcache`.

See [Build and test caching](https://pkg.go.dev/cmd/go#hdr-Build_and_test_caching)
and [go#24355](https://github.com/golang/go/issues/24355) for more details.

#### Go playgrounds

Regular Go playgrounds come with cgo disabled therefore it's not possible
to experiment with Pygolo without setting up a local environment.

#### Building extensions on Alpine Linux

It's not possible to build Python extensions on Alpine Linux or any other
musllibc based system. See https://gitlab.com/pygolo/py/-/issues/15.
