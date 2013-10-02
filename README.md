JackDB
======

![screenshot](https://raw.github.com/tristanwietsma/jackdb/master/docs/screenshot.png)

What is JackDB?
---------------

JackDB is a concurrent key-value server. It supports get, set, publish, subscribe, and delete

 **That is it.** Five commands.

The underlying data structure is [MetaStore](https://github.com/tristanwietsma/metastore), which is an abstraction over a string map that divides the key-space into buckets for finer lock resolution.

What isn't JackDB?
------------------

Jack is not Redis. Jack isn't a fan of large APIs and does not support persistence.

The project is currently in beta and should be considered volatile. Contributions are welcome.

Installation
------------

You can use the Go tool to install the library and dependencies:

    export GOPATH=<where you store your Go code>
    go get -u github.com/tristanwietsma/jackdb

The project currently ships with the server, **jackd**, and a command line interface, **jack-cli**. Since the Go tool doesn't like multiple build targets in the same project, you need to build them separately:

    cd $GOPATH/src/github.com/tristanwietsma/jackdb
    make

The 'make' will build both  **jackd** and **jack-cli**, as well as move them to *$GOCODE/bin*.

Usage
-----

To start the server, run jackd:

    $ jackd
    2013/10/02 15:26:25 created storage with 1000 buckets
    2013/10/02 15:26:25 server started on port 2000
    ...

To start the commandline tool (which is still *very* young):

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

The Future
----------

* Currently, the cli splits arguments on space (so you can't have spaces in your value). This will get better; the cli was hacked together pretty fast for testing purposes.

* I want to test it and add some benchmarks in the near future. I'm sure there is room for optimization. I hear a good argument in favor of a particular bell or whistle, features are capped.
