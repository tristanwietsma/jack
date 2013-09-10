JackDB v0.1
===========

A Concurrent Key-Value Store

Requires
--------

* Go 1.1
* Python 2.7.x (for tests)

Install
-------

Assuming you have set your GOPATH, use the go tool:

    go get github.com/tristanwietsma/jackdb

Run
---

    cd $GOPATH/bin
    ./jackdb [--port <int>]

Run Tests
---------

To evaluate get and set (with server running):

    cd $GOPATH/src/github.com/tristanwietsma/jackdb/tests
    sh run-tests.sh

This will run 50 clients for 200 commands. The time required will be displayed for each action. On my i7, I get the following:

    $ sh run-tests.sh 
    .038331616
    .047540910

By comparison, the equivalent Redis benchmarks are around 0.12 for both get and set.

Roadmap
-------

* Testing

* More testing

* Might run some tests

* API
