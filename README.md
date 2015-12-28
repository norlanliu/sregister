# sregister
![Build Status](https://travis-ci.org/norlanliu/sregister.svg)

SRegisater is a service register component of a service-discovery framework SFinder which
actully is another wheel similar with [Smartstack][1], SRegister is analagous to [nerve][3].

it monitor the service and report to [ETCD][2] with the address and port of the service.

## Installation

Get the source code using the command:
> $ go get github.com/norlanliu/sregister

go to the sregister directory
> $ cd src/github.com/norlanliu/sregisater

make the source code
> $ make

install
> $ sudo make install

Then, start sregister

1.Redhat/Centos/Fedora
start the sregister with systemctl
> $ systemctl enable sregister; systemctl start sregister
> $ systemctl status sregister

2.Debain/Ubuntu
> $ sregister -h		# get the help
> $ sudo sregister &

## Configuration
The default application configuration file of SRegister is at `/etc/sfinder/sregister/sregister.conf`
* `SREG_LOG_DIR`: log directory of sregister
* `SREG_SERVICES_DIR`: path to a directory in which each json file whill be interpreted as a service, you can add/delete/modify the service file at any time without stop the sregister.

### Services Configuration
The configuration follow the [nerve][3].
The configuration contains the following options:

* `name`: the name of the service
* `host`: the default host on which to make service checks; you should make this your *public* ip to ensure your service is publically accessible
* `port`: the default port for service checks; nerve will report the `host`:`port` combo via your chosen reporter
* `reporterType`: the mechanism used to report up/down information; depending on the reporter you choose, additional parameters may be required. Defaults to `etcd`
* `checkInterval`: the frequency with which service checks will be initiated; defaults to `2s`
* `weight`: a positive integer weight value which can be used to affect the haproxy backend weighting in sdiscovery.
* `checks`: a list of checks that nerve will perform; if all of the pass, the service will be registered; otherwise, it will be un-registered

[1]: http://nerds.airbnb.com/smartstack-service-discovery-cloud/
[2]: https://github.com/coreos/etcd
[3]: https://github.com/airbnb/nerve

# TODO
1. add http service watcher
2. add mysql service watcher
3. add docker servcie watcher
4. support [Zookeeper](http://zookeeper.apache.org/)
