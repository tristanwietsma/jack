JackDB (beta)
=============

[![Build Status](https://travis-ci.org/tristanwietsma/jackdb.png?branch=master)](https://travis-ci.org/tristanwietsma/jackdb)

What is JackDB?
---------------

JackDB is a concurrent key-value server. It supports get, set, publish, subscribe, and delete. **That is it.** Five commands.

The underlying data structure is [MetaStore](https://github.com/tristanwietsma/metastore), which is an abstraction over a string map that divides the key-space into buckets for finer lock resolution.

What isn't JackDB?
------------------

Jack is not Redis. Jack isn't a fan of large APIs.

The project is currently in beta and should be considered volatile. Contributions are welcome.

Installation
------------

You can use the Go tool to install the libary and dependencies:

    export GOPATH=<where you store your Go code>
    go get -u github.com/tristanwietsma/jackdb

The project currently ships with the server, **jackd**, and a command line interface, **jack-cli**. Since the Go tool doesn't like multiple build targets in the same project, you need to build them separately:

    cd $GOPATH/src/github.com/tristanwietsma/jackdb
    make

The 'make' will build both  **jackd** and **jack-cli**, as well as move them to *$GOCODE/bin*.
