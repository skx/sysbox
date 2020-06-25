[![Go Report Card](https://goreportcard.com/badge/github.com/skx/sysbox)](https://goreportcard.com/report/github.com/skx/sysbox)
[![license](https://img.shields.io/github/license/skx/sysbox.svg)](https://github.com/skx/sysbox/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/skx/sysbox.svg)](https://github.com/skx/sysbox/releases/latest)


* [SysBox](#sysbox)
  * [Installation](#installation)
  * [Bash Completion](#bash-completion)
* [Overview](#overview)
* [Tools](#tools)
* [Future Additions?](#future-additions)
* [Github Setup](#github-setup)


# SysBox

This repository is the spiritual successor to my previous [sysadmin-utils repository](https://github.com/skx/sysadmin-util)

The idea here is to collect simple utilities and package them as a single binary, written in go, in a similar fashion to the `busybox` utility.

## Installation

Installation upon a system which already contains a go-compiler should be as simple as:

```
$ go get github.com/skx/sysbox
```

If you prefer you can find binary releases upon our [download page](https://github.com/skx/sysbox/releases).


## Bash Completion

The [subcommand library](https://github.com/skx/subcommands) this application uses has integrated support for the generation of bash completion scripts.

To enable this add the following to your bash configuration-file:

```
source <(sysbox bash-completion)

```


# Overview

This application is built, and distributed, as a single-binary named`sysbox`, which implements a number of sub-commands.

You can either run the tools individually, taking advantage of the [bash completion](#bash-completion) support to complete the subcommands and their arguments:

    $ sysbox foo ..
    $ sysbox bar ..

Or you can create symlinks to allow specific tool to be executed without the need to specify a subcommand:

    $ ln -s sysbox foo
    $ foo ..

This process of creating symlinks can be automated via the use of the `sysbox install` sub-command, which would allow you to install the tools globally like so:

    $ go get github.com/skx/sysbox
    $ $GOPATH/bin/sysbox install -binary $GOPATH/bin/sysbox -directory /usr/local/bin | sudo sh



# Tools

The tools in this repository started out as being simple ports of the tools in my [previous repository](https://github.com/skx/sysadmin-util), however I've now started to expand them and fold in things I've used/created in the past.

You can view a summary of the available subcommands via:

    $ sysbox help

More complete help for each command should be available like so:

    $ sysbox help sub-command

Examples are included where useful.


## calc

A simple calculator, which understands floating point-operations, unlinke `expr`.

The calculator supports either execution of sums via via the command-line, or as an interactive REPL environment:

```
$ sysbox calc 3.1 + 2.7
5.800000

$ ./sysbox calc
calc> let a = 1/3
0.333333
calc> a * 9
3
calc> exit
$
```

## choose-file

This subcommand presents a console-based UI to select a file.  The file selected will be displayed upon STDOUT.  The list may be filtered via an input-field.

Useful for launching videos, emulators, etc:

* `sysbox choose-file -execute="xine -g --no-logo --no-splash -V=40 {}" ~/Videos`
  * Choose a file, and execute `xine` with that filename as one of the arguments.
* `xine $(sysbox choose-file ~/Videos)`
  * Use the STDOUT result to launch instead.

The first form is preferred, because if the selection is canceled nothing happens.  In the second-case `xine` would be launched with no argument.


## choose-stdin

Almost identical to `choose-file`, but instead of allowing the user to choose from a filename it allows choosing from the contents read on STDIN.  For example you might allow choosing a directory:

```
$ find ~/Repos -type d | sysbox choose-stdin -execute="firefox {}"
```


## chronic

The chronic command is ideally suited to wrap cronjobs, it runs the command you specify as a child process and hides the output produced __unless__ that process exits with a non-zero exit-code.


## collapse

This is a simple tool which will read STDIN, and output the content without any extra white-space:

* Leading/Trailing white-space will be removed from every line.
* Empty lines will be skipped entirely.


## cpp

Something _like_ the C preprocessor, but supporting only the ability to include files, and run commands via `#include` and `#execute` respectively:

    #include "file/goes/here"
    #execute ls -l | wc -l


## env-template

Perform expansion, via environmental variables, on simple golang templates.

You can freely use the built-in golang template facilities, for example please see the sample template here [cmd_env_template.tmpl](cmd_env_template.tmpl), and the the examples included in the [text/template documentation](https://golang.org/pkg/text/template/).

> As an alternative you can consider the `envsubst` binary contained in your system's `gettext{-base}` package.


## exec-stdin

Read STDIN, and allow running a command for each line.  You can refer to
the line read either completely, or by fields.

For example:

```
$ ps -ef | sysbox exec-stdin echo field1:{1} field2:{2} line:{}
```

See the usage-information for more details (`sysbox help exec-stdin`), but consider this a simple union of `awk`, `xargs`, and GNU parallel (since we can run multiple commands in parallel).


## fingerd

A trivial finger-server.


## httpd

A simple HTTP-server.  Allows serving to localhost, or to the local LAN.


## http-get

Very much "curl-lite", allows you to fetch the contents of a remote URL.  SSL errors, etc, are handled, but only minimal options are supported.


## install

This command allows you to install symlinks to the sysbox-binary, for ease of use:

```
$ sysbox install -binary=$(pwd)/sysbox -directory=~/bin | sh
```


## ips

This tool lets you easily retrieve a list of local, or global, IPv4 and
IPv6 addresses present upon your local host.  This is a little simpler
than trying to parse `ip -4 addr list`, although that is also the
common approach.


## make-password

This tool generates a single random password each time it is executed, it is designed to be quick and simple to use, rather than endlessly configurable.


## peerd

This deamon provides the ability to maintain a local list of available cluster-members, via the JSON file located at `/var/tmp/peerd.json`.

See the usage-information for more (`sysbox help peerd`).


## run-directory

Run every executable in the given directory, optionally terminate if any command returns a non-zero exit-code.

> The exit-code handling is what inspired this addition; the Debian version of `run-parts` supports this, but the CentOS version does not.


## splay

This tool allows sleeping for a random amount of time.  This solves the problem when you have a hundred servers all running a task at the same time, triggered by `cron`, and you don't want to overwhelm a central resource that they each consume.


## ssl-expiry

A simple utility to report upon the number of hours, and days, until a given TLS certificate (or any intermediary in the chain) expires.

Ideal for https-servers, but also TLS-protected SMTP hosts, etc.


## timeout

Run a command, but kill it after the given number of seconds.  The command is executed with a PTY so you can run interactive things such as `top`, `mutt`, etc.


## torrent

Simple bittorrent client, which allows downloading a magnet-based torrent.  For example to download an Ubuntu ISO:

```
$ sysbox torrent magnet:?xt=urn:btih:ZOCMZQIPFFW7OLLMIC5HUB6BPCSDEOQU
```


## tree

Trivial command to display the contents of a filesystem, as a nested tree.  This is similar to the standard `tree` command, without the nesting and ASCII graphics.


## urls

Extract URLs from the named files, or STDIN.  URLs are parsed naively with a simple regular expression and only `http` and `https` schemes are recognized.


## validate-json

Validate `*.json` files from the current working-directory, or the named directory, recursively.


## validate-yaml

Validate `*.yaml`/`*.yml` files from the current working-directory, or the named directory, recursively.


## with-lock

Allow running a command with a lock-file to prevent parallel executions.

This is perfect if you fear your cron-jobs will start slowing down and overlapping executions will cause problems.



# Future Additions?

Unlike the previous repository I'm much happier to allow submissions of new utilities, or sub-commands, in this repository.


# Github Setup

This repository is configured to run tests upon every commit, and when
pull-requests are created/updated.  The testing is carried out via
[.github/run-tests.sh](.github/run-tests.sh) which is used by the
[github-action-tester](https://github.com/skx/github-action-tester) action.

Releases are automated in a similar fashion via [.github/build](.github/build),
and the [github-action-publish-binaries](https://github.com/skx/github-action-publish-binaries) action.

Steve
--
