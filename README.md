lbx
===

`lbx` is a set of tools to ease development in LHCb.

## Installation

```sh
$ go get github.com/lhcb-org/lbx
$ lbx version
0.1-20140528
```

## Usage

### Help

```sh
$ lbx help
lbx - tools for development.

Commands:

    init        initialize a local development project.
    pkg         add, remove or inspect sub-packages
    version     print out script version

Use "lbx help <command>" for more information about a command.
```

### init

```sh
$ lbx help init
Usage: lbx init [options] <project-name> [<project-version>]

init initialize a local development project.

ex:
 $ lbx init Gaudi
 $ lbx init Gaudi HEAD
 $ lbx init -name mydev Gaudi

Options:
  -c="x86_64-slc6-gcc48-opt": runtime platform
  -dev-dirs="": path-list to prepend to the projects-search path
  -lvl=0: message level to print
  -name="": name of the local project (default: <project>Dev_<version>)
  -nightly="": specify a nightly to use. e.g. slotname,Tue
  -user-area=".": use the specified path as User_release_area instead of ${User_release_area}
  -v=false: enable verbose mode
```

### pkg

```sh
$ lbx help pkg co
Usage: lbx pkg co [options] <pkg-uri> [<pkg-version>]

co adds a package to the current workarea.

ex:
 $ lbx pkg co MyPackage vXrY

Options:
  -go=false: use the go version
  -v=false: enable verbose output
```
