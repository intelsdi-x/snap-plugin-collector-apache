<!--
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

[![Build Status](https://travis-ci.org/intelsdi-x/snap-plugin-collector-apache.svg?branch=master)](https://travis-ci.com/intelsdi-x/snap-plugin-collector-apache)

# Snap collector plugin - Apache

This plugin collects metrics from the Apache Webserver for `mod_status`: `http://your.server.name/server-status?auto`. `?auto` is the machine-readable format for the status file.

It's used in the [Snap framework](http://github.com/intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

In order to use this plugin you need Apache2 installed. 

### Operating systems
* Linux/amd64

### Installation
#### Apache
```
sudo apt-get install apache2
```
#### Snap
You can get the pre-built binaries for your OS and architecture at Snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page.

##### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-apache
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-apache.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
(It may take a while to pull dependencies if you don't have them already.)

This builds the plugin in `/build/${GOOS}/${GOARCH}`

##### Run tests
```
TODO
```

### Configuration and Usage
#### Apache
The [Apache mod_status documentation](https://httpd.apache.org/docs/2.4/mod/mod_status.html) says to enable mod_status in httpd.conf, but it was located in mods-enabled/status.conf for me:
```
$ ls /etc/apache2/
apache2.conf  conf-available  conf-enabled  envvars  magic  mods-available  mods-enabled  ports.conf  sites-available  sites-enabled
```
Make sure mod_status is enabled with the following uncommented. Extended status needs to be on. 
If you want the mod status url to be `mod_status` instead of `server-status` just change the location line and the `apache_suffix` in the config.json file to match.

```
ExtendedStatus on
<Location /mod_status>
  SetHandler server-status
</Location>
```
If no changes needed to be made, run the following to start apache.
```
$ sudo service apache2 start 
```
If any changes needed to be made, run the following to restart apache after you save the file.
```
$ sudo service apache2 restart
```
If it doesn't seem like anything changed when you start up the server, it likely means that you had other Apache instances running.  
You can stop Apache2 and then if there are any Apache processes running, kill them, and start Apache:  
```
$ sudo service apache2 stop
$ ps aux | grep apache
$ kill -9 <PIDs> 
$ sudo service apache2 start
```
Check to see if Apache2 is running:
```
$ service apache2 status
 * apache2 is running
```
#### Snap
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported  
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build/${GOOS}/${GOARCH}`

## Documentation
There are a number of other resources you can review to learn to use this plugin:
* [Apache mod_status documentation](https://httpd.apache.org/docs/2.4/mod/mod_status.html)
* [Snap Apache examples](#examples)
* [Snap Apache JSON task example](examples/tasks/apache-file.json)

### Collected Metrics
All metrics gathered by this plugin are exposed by the [status file](https://httpd.apache.org/docs/2.4/mod/mod_status.html#machinereadable) produced by mod_status.  

This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description (optional)
----------|-----------|-----------------------
/intel/apache/BusyWorkers|float64|Busy workers
/intel/apache/BytesPerReq|float64|Bytes transferred per request
/intel/apache/BytesPerSec|float64|Bytes transferred per second
/intel/apache/CPULoad|float64|CPU load
/intel/apache/ConnsAsyncClosing|float64|Asynchronous closing connections
/intel/apache/ConnsAsyncKeepAlive|float64|Asynchronous keepalive connections
/intel/apache/ConnsAsyncWriting|float64|Asynchronous writing connections
/intel/apache/ConnsTotal|float64|Total connections
/intel/apache/IdleWorkers|float64|Idle workers
/intel/apache/ReqPerSec|float64|Requests per second
/intel/apache/Total_Accesses|float64|Total accesses
/intel/apache/Total_kBytes|float64|Total kBytes
/intel/apache/Uptime|float64|Server uptime
/intel/apache/workers/Closing|float64|Closing connection
/intel/apache/workers/DNSLookup|float64|DNS Lookup
/intel/apache/workers/Finishing|float64|Gracefully finishing
/intel/apache/workers/Idle_Cleanup|float64|Idle cleanup of worker
/intel/apache/workers/Keepalive|float64|Keepalive (read)
/intel/apache/workers/Logging|float64|Logging
/intel/apache/workers/Open|float64|Open slot with no current process
/intel/apache/workers/Reading|float64|Reading Request
/intel/apache/workers/Sending|float64|Sending Reply
/intel/apache/workers/Starting|float64|Starting up
/intel/apache/workers/Waiting|float64|Waiting for Connection


### Examples
If this is your directory structure:
```
$GOPATH/src/github.com/intelsdi-x/snap/
$GOPATH/src/github.com/intelsdi-x/snap-plugin-collector-apache/
```

In one terminal window in the /snap directory: Running snapd with auto discovery, log level 1, and trust disabled. The config.json file has the webserver configuration parameters.
```
$ $SNAP_PATH/snapd -l 1 -t 0 --config ../snap-plugin-collector-apache/config.json 
```
Download desired processor and publisher plugins eg.
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
```
Load collector, processor and publisher
```
$ $SNAP_PATH/snapctl plugin load snap-plugin-collector-apache
$ $SNAP_PATH/snapctl plugin load snap-plugin-publisher-file
```
Create task manifest for writing to a file. See [`../snap-plugin-collector-apache/examples/tasks/apache-file.json`](../snap-plugin-collector-apache/examples/tasks/apache-file.json):
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/apache/CPULoad": {},
                "/intel/apache/BytesPerSec": {},
                "/intel/apache/workers/Sending": {}
            },
            "publish": [
                {
                    "plugin_name": "file",                            
                    "config": {
                        "file": "/tmp/snap-apache-file.log"
                    }
                }
           ]
        }
    }
}
```
Another terminal window, also in /snap:
```
$ $SNAP_PATH/snapctl task create -t ../snap-plugin-collector-apache/examples/tasks/apache-file.json
```
/tmp/snap-apache-file.log
```
2016-01-27 15:08:51.09527825 -0800 PST|[intel apache BytesPerSec]|92.0861|127.0.0.1:80
2016-01-27 15:08:51.09527251 -0800 PST|[intel apache CPULoad]|.0209417|127.0.0.1:80
2016-01-27 15:08:51.095292429 -0800 PST|[intel apache workers Sending]|1|127.0.0.1:80
2016-01-27 15:08:52.096059795 -0800 PST|[intel apache BytesPerSec]|92.086|127.0.0.1:80
2016-01-27 15:08:52.096056768 -0800 PST|[intel apache CPULoad]|.0209417|127.0.0.1:80
2016-01-27 15:08:52.096089108 -0800 PST|[intel apache workers Sending]|1|127.0.0.1:80
```

### Roadmap
The next step for this plugin is to make sure it works with Lightppd and it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-apache/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-apache/pulls).

## Community Support
This repository is one of **many** plugins in the **Snap framework**: a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Tiffany Jernigan](https://github.com/tiffanyfj)
* Author: [Dan Pittman](https://github.com/danielscottt)

**Thank you!** Your contribution is incredibly important to us.
