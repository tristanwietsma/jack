JackDB (beta)
=============

[![Build Status](https://travis-ci.org/tristanwietsma/jackdb.png?branch=master)](https://travis-ci.org/tristanwietsma/jackdb)

JackDB is a concurrent key-value server. It supports get, set, publish, subscribe, and delete. The project is currently in beta.

Installation
------------

JackDB is built on top of [MetaStore](https://github.com/tristanwietsma/metastore). You can use the Go tool to install the libary and dependencies:

    export GOPATH=<where you store your Go code>
    go get -u github.com/tristanwietsma/jackdb

The project currently ships with the server, *jackd*, and a command line interface, *jack-cli*. Since the Go tool doesn't like multiple build targets in the same project, you need to build them separately:

    cd $GOPATH/src/github.com/tristanwietsma/jackdb
    make

The 'make' will build both  *jackd* and *jack-cli*, as well as move them to **$GOCODE/bin**.
