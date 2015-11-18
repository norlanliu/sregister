# sregister
![Build Status](https://travis-ci.org/norlanliu/sregister.svg)

SRegisater is a service register part of a service-discovery framework SFinder. 

it monitor the service and report to [ETCD](https://github.com/coreos/etcd) with
the address and port of the service.

## Installation

Get the source code using the command:
> go get github.com/norlanliu/sregister

make the source code
> make

install
> sudo make install

then you can start the sregister
> sudo sregister &

# TODO
1. add http service watcher
2. add mysql service watcher
3. add docker servcie watcher
4. support [Zookeeper](http://zookeeper.apache.org/)
