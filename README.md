JackDB
======

A Concurrent Key-Value Store

Requires
--------

* Go 1.1
* Python 2.7.x (for tests)

Install
-------

    git clone git@github.com:tristanwietsma/jackdb.git
    cd jackdb
    go build jackdb.go

Run
---

    ./jackdb [--port <int>]

Run Tests
---------

To evaluate get and set,

    cd jackdb/tests
    sh run-tests.sh

This will run 50 clients for 200 commands. The time required will be displayed for each action.
