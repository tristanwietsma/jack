Jack
====

[![Build Status](https://travis-ci.org/tristanwietsma/jack.png?branch=master)](https://travis-ci.org/tristanwietsma/jack) [![GoDoc](https://godoc.org/github.com/tristanwietsma/jack?status.svg)](http://godoc.org/github.com/tristanwietsma/jack)

![screenshot](https://raw.github.com/tristanwietsma/jack/master/docs/screenshot.png)

What is Jack?
-------------

Jack is a proof-of-concept concurrent key-value server. It supports get, set, publish, subscribe, and delete

The underlying data structure is [MetaStore](https://github.com/tristanwietsma/metastore), which is an abstraction over a string map that divides the key-space into buckets for finer lock resolution.


Installation
------------

You can use the Go tool to install the library and dependencies:

    export GOPATH=<where you store your Go code>
    go get -u github.com/tristanwietsma/jack

The project currently ships with the server, **jackd**, and a command line interface, **jack-cli**. Since the Go tool doesn't like multiple build targets in the same project, you need to build them separately:

    cd $GOPATH/src/github.com/tristanwietsma/jack
    make

The 'make' will build both  **jackd** and **jack-cli**, as well as move them to *$GOCODE/bin*.

Usage
-----

To start the server, run jackd:

    $ jackd
    2013/10/02 15:26:25 created storage with 1000 buckets
    2013/10/02 15:26:25 server started on port 2000
    ...

To start the command-line tool:

    $ jack-cli
    jack> set key123 val567
    jack> 1
    jack> get key123
    jack> key123 := val567
    jack> pub key123 765lav
    jack> 1
    jack> sub key123
    ...

The commands are all three characters and not case sensitive. The following should be self-explanatory:

    jack> set <key> <value>

    jack> get <key> [<key> ...] // supports multiple keys

    jack> pub <key> <value>

    jack> del <key> [<key> ...] // supports multiple keys

    jack> sub <key> [<key> ...] // supports multiple keys
