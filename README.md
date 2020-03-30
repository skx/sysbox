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

* Document tools here.

## Additions?

Unlike the previous repository I'm much happier to allow submissions of new utilities, or sub-commands, in this repository.


Steve
