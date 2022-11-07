DISCONTINUATION OF PROJECT. 

This project will no longer be maintained by Intel.

This project has been identified as having known security escapes.

Intel has ceased development and contributions including, but not limited to, maintenance, bug fixes, new releases, or updates, to this project.  

Intel no longer accepts patches to this project.
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

# DISCONTINUATION OF PROJECT 

**This project will no longer be maintained by Intel.  Intel will not provide or guarantee development of or support for this project, including but not limited to, maintenance, bug fixes, new releases or updates.  Patches to this project are no longer accepted by Intel. If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the community, please create your own fork of the project.**


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
export TEST_TYPE=small
make test
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

## Documentation
There are a number of other resources you can review to learn to use this plugin:
* [Apache mod_status documentation](https://httpd.apache.org/docs/2.4/mod/mod_status.html)
* [Snap Apache examples](#examples)
* [Snap Apache JSON task example](examples/tasks/apache-file.json)

### Collected Metrics
All metrics gathered by this plugin are exposed by the [status file](https://httpd.apache.org/docs/2.4/mod/mod_status.html#machinereadable) produced by mod_status.  
This plugin has the ability to gather the following metrics:
[Available Metrics](METRICS.md)

This plugin will provide two seperate metric catalogs based on whether safe or unsafe collection is configured.

### Examples
If this is your directory structure:
```
$GOPATH/src/github.com/intelsdi-x/snap/
$GOPATH/src/github.com/intelsdi-x/snap-plugin-collector-apache/
```

In one terminal window in the /snap directory: Running snapteld with auto discovery, log level 1, and trust disabled. The config.json file has the webserver configuration parameters.
```
$ snapteld -l 1 -t 0
```
Download desired publisher plugin eg.
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
```
Load collector and publisher
```
$ snaptel plugin load snap-plugin-collector-apache
$ snaptel plugin load snap-plugin-publisher-file
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
          "config": {
            "/intel/apache": {
              "apache_mod_status_url": "https://www.apache.org/server-status?auto"
            }
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
$ snaptel task create -t ../snap-plugin-collector-apache/examples/tasks/apache-file.json
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

### Unsafe collection
Not all version of the apache status endpoint return the same metrics. To allow for consistent collections from different apache servers safe collection is enabled for the plugin by default. To override safe collection and expose all the available metrics alter the [snapteld global configuration](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPTELD_CONFIGURATION.md) to include the following:

```
control:
    plugins:
      collector:
        apache:
          all:
            safe: false
```

The provided [example config](examples/configs/config.yaml) can be used to load this plugin in unsafe mode

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
