_Thanks for all the suggestions and early review! I've got a lot of work ahead of me._

_Update 9/25/2013: I've spun the backend storage structure into a separate project called [MetaStore](https://github.com/tristanwietsma/metastore). MetaStore splits the map into buckets for finer lock granularity; hashing with fnv. (Hat tip to bonekeeper)_

JackDB v0.01
============

Concurrent key-value server in Go

Supports GET, SET, DEL, PUB, and SUB

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
    .038331616 <-- 10,000 'sets' in ~0.04 seconds
    .047540910 <-- 10,000 'gets' in ~0.05 seconds

For comparison, Redis benchmarks are around 0.12 seconds for both get and set on the same machine.

I've also got examples of publish/subscribe in the testing folder. Publishing is a little different than some key-value stores: JackDB waits on a stream, rather than another publish command (as with Redis). This was a design decision that relates to using publish for streaming input over a dedicated connection (think high frequency sensors or, perhaps, tick data).

To Do
-----

* Design docs...

* Protocol specification...

* Testing suite...

* Some support for persistence (LOAD, SAVE) and administration...

* Improve storage system; I'm leaning towards dividing the map into buckets and hashing the keys...

