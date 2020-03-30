

* [SysBox](#sysbox)
  * [Overview](#overview)
  * [Tools](#tools)
     * [collapse](#collapse)
     * [make-password](#make-password)
     * [splay](#splay)
     * [ssl-expiry](#ssl-expiry)
     * [with-lock](#with-lock)
  * [Future Additions?](#future-additions)


# SysBox

This repository is the spiritual successor to my previous [sysadmin-utils repository](https://github.com/skx/sysadmin-util)

The idea here is to collect simple utilities and package them as a single binary, written in go, in a similar fashion to the `busybox` utility.


## Overview

Much like busybox there will be a single binary, `sysbox`, which implements a number of sub-commands.

You can either run the tools individually:

    $ sysbox foo ..
    $ sysbox bar ..

Or you can create symlinks to allow specific tool to be executed without the need to specify a subcommand:

    $ ln -s sysbox foo
    $ foo ..


## Tools

The tools in this repository are currently those which are being ported from the [previous repository](https://github.com/skx/sysadmin-util).   Full help is provided via:

    $ sysbox help

And:

    $ sysbox help sub-command

For each command you may create a symlink to allow it to be invoked without the prefix:

    $ ln -s sysbox with-lock
    $ ./with-lock ls


### collapse

This is a simple tool which will read STDIN, and output the content without any extra whitespace:

* Leading/Trailing whitespace will be removed from every line.
* Empty lines will be skipped entirely.


### make-password

This tool generates a single random password each time it is executed, it is designed to be quick and simple to use, rather than endlessly configurable.

### splay

This tool allows sleeping for a random amount of time.  This solves the problem when you have a hundred servers all running a task at the same time, triggered by cron, and you don't want to overwhelm a central host that they each talk to.

### ssl-expiry

A simple utility to report upon the number of hours, and days, until a given TLS certificate (or any intermediary in the chain) expires.

Ideal for https-servers, but also TLS-protected SMTP hosts, etc.

### with-lock

Allow running a command with a lock-file to prevent parallel executions.

This is perfect if you fear your cron-jobs will start slowing down and overlapping executions will cause problems.


## Future Additions?

Unlike the previous repository I'm much happier to allow submissions of new utilities, or sub-commands, in this repository.


Steve
