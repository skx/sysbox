[![Go Report Card](https://goreportcard.com/badge/github.com/skx/sysbox)](https://goreportcard.com/report/github.com/skx/sysbox)
[![license](https://img.shields.io/github/license/skx/sysbox.svg)](https://github.com/skx/sysbox/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/skx/sysbox.svg)](https://github.com/skx/sysbox/releases/latest)


* [SysBox](#sysbox)
  * [Installation](#installation)
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



# Overview

This application is built, and distributed, as a single-binary named`sysbox`, which implements a number of sub-commands.

You can either run the tools individually:

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

You can view a brief list of the commands via:

    $ sysbox help

More complete help for each command should be available like so:

    $ sysbox help sub-command

Examples are included where useful.


## calc

A simple calculator, which understands floating points, which `expr` never does.

Portions of this code from [Porting Eval to Go](https://thorstenball.com/blog/2016/11/16/putting-eval-in-go/), by Thorsten Ball.  (I expanded it to support parenthesis, for precedence, and the use of floating-point numbers rather than integers.)


## chronic

The chronic command is ideally suited to wrap cronjobs, it runs the command you specify as a child process and hides the output produced __unless__ that process exits with a non-zero exit-code.


## collapse

This is a simple tool which will read STDIN, and output the content without any extra whitespace:

* Leading/Trailing whitespace will be removed from every line.
* Empty lines will be skipped entirely.


## env-template

Perform expansion, via environmental variables, on simple golang templates.

You can freely use the built-in golang template facilities, for example please see the sample template here [cmd_env_template.tmpl](cmd_env_template.tmpl), and the the examples included in the [text/template documentation](https://golang.org/pkg/text/template/).

> As an alternative you can consider the `envsubst` binary contained in your system's `gettext{-base}` package.


## exec-stdin

Read STDIN, and allow running a command for each line.  You can refer to
the line read either completely, or by fields.

For example:

   $ ps -ef | sysbox exec-stdin echo field1:{1} field2:{2} line:{}


## fingerd

A trivial finger-server.


## httpd

A simple HTTP-server.  Allows serving to localhost, or to the local LAN.



## http-get

Very much "curl-lite", allows you to fetch the contents of a remote URL.  SSL errors, etc, are handled, but there are zero configuration options.


## install

This command allows you to install symlinks to the binary, for ease of use:

    $ sysbox install -binary=$(pwd)/sysbox -directory=~/bin | sh


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
