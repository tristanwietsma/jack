[![Build Status](https://travis-ci.org/tristanwietsma/jackdb.png?branch=master)](https://travis-ci.org/tristanwietsma/jackdb)

_Update 9/28/2013: Major revision coming soon. The project is currently not working._

_Update 9/25/2013: I've spun the backend storage structure into a separate project called [MetaStore](https://github.com/tristanwietsma/metastore). MetaStore splits the map into buckets for finer lock granularity; hashing with fnv. (Hat tip to bonekeeper)_

JackDB
======

Concurrent key-value server in Go

Supports GET, SET, DEL, PUB, and SUB
