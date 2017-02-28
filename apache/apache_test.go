/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apache

import (
	"net/http"
	"strings"
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	apacheURL    = "http://127.0.0.1:80/server-status?auto"
	legacyStatus = `Total Accesses: 10
Total kBytes: 1024
CPULoad: .0582479
Uptime: 4292
ReqPerSec: 0.123
BytesPerReq: 1.501
BytesPerSec: 1.123
BusyWorkers: 10
IdleWorkers: 50
ConnsTotal: 100
ConnsAsyncWriting: 1
ConnsAsyncKeepAlive: 2
ConnsAsyncClosing: 3
Scoreboard: CDGIKLRWS_...`

	exampleStatus = `localhost
ServerVersion: Apache/2.4.18 (Ubuntu)
ServerMPM: event
Server Built: 2016-07-14T12:32:26
CurrentTime: Thursday, 01-Dec-2016 20:09:30 UTC
RestartTime: Thursday, 01-Dec-2016 20:01:42 UTC
ParentServerConfigGeneration: 35
ParentServerMPMGeneration: 100
ServerUptimeSeconds: 4292
ServerUptime: 7 minutes 48 seconds
Load1: 0.01
Load5: 0.03
Load15: 0.05
Total Accesses: 10
Total kBytes: 1024
CPUUser: 11001
CPUSystem: 1305
CPUChildrenUser: 315
CPUChildrenSystem: 135
ReqPerSec: 0.123
BytesPerSec: 1.123
BusyWorkers: 10
IdleWorkers: 50
ConnsTotal: 100
ConnsAsyncWriting: 1
ConnsAsyncKeepAlive: 2
ConnsAsyncClosing: 3
Scoreboard: CDGIKLRWS_...`

	exampleUnknownKeysStatus = `localhost
ServerVersion: Apache/2.4.18 (Ubuntu)
ServerMPM: event
Server Built: 2016-07-14T12:32:26
CurrentTime: Thursday, 01-Dec-2016 20:09:30 UTC
RestartTime: Thursday, 01-Dec-2016 20:01:42 UTC
ParentServerConfigGeneration: 35
ParentServerMPMGeneration: 100
ServerUptimeSeconds: 4292
ServerUptime: 7 minutes 48 seconds
Load1: 0.01
Load5: 0.03
Load15: 0.05
Total Accesses: 10
Total kBytes: 1024
CPUUser: 11001
CPUSystem: 1305
CPUChildrenUser: 315
CPUChildrenSystem: 135
ReqPerSec: 0.123
BytesPerSec: 1.123
BusyWorkers: 10
IdleWorkers: 50
ConnsTotal: 100
ConnsAsyncWriting: 1
ConnsAsyncKeepAlive: 2
ConnsAsyncClosing: 3
Scoreboard: CDGIKLRWS_...
UnknownKeyhere: 123`

	allValues = map[string]interface{}{
		"intel.apache.CPULoad":                      0,
		"intel.apache.BytesPerReq":                  0,
		"intel.apache.ServerVersion":                "Apache/2.4.18 (Ubuntu)",
		"intel.apache.ServerMPM":                    "event",
		"intel.apache.Server_Built":                 "2016-07-14T12:32:26",
		"intel.apache.CurrentTime":                  "Thursday, 01-Dec-2016 20:09:30 UTC",
		"intel.apache.RestartTime":                  "Thursday, 01-Dec-2016 20:01:42 UTC",
		"intel.apache.ParentServerConfigGeneration": 35,
		"intel.apache.ParentServerMPMGeneration":    100,
		"intel.apache.ServerUptime":                 "7 minutes 48 second",
		"intel.apache.Load1":                        0.01,
		"intel.apache.Load5":                        0.03,
		"intel.apache.Load15":                       0.05,
		"intel.apache.Total_Accesses":               10,
		"intel.apache.Total_kBytes":                 1024,
		"intel.apache.CPUUser":                      11001,
		"intel.apache.CPUSystem":                    1305,
		"intel.apache.CPUChildrenUser":              315,
		"intel.apache.CPUChildrenSystem":            135,
		"intel.apache.Uptime":                       4292,
		"intel.apache.ReqPerSec":                    0.123,
		"intel.apache.BytesPerSec":                  1.123,
		"intel.apache.BusyWorkers":                  10,
		"intel.apache.IdleWorkers":                  50,
		"intel.apache.ConnsTotal":                   100,
		"intel.apache.ConnsAsyncWriting":            1,
		"intel.apache.ConnsAsyncKeepAlive":          2,
		"intel.apache.ConnsAsyncClosing":            3,
		"intel.apache.workers.Closing":              1,
		"intel.apache.workers.DNSLookup":            1,
		"intel.apache.workers.Finishing":            1,
		"intel.apache.workers.Idle_Cleanup":         1,
		"intel.apache.workers.Keepalive":            1,
		"intel.apache.workers.Logging":              1,
		"intel.apache.workers.Open":                 3,
		"intel.apache.workers.Reading":              1,
		"intel.apache.workers.Sending":              1,
		"intel.apache.workers.Starting":             1,
		"intel.apache.workers.Waiting":              1,
	}
	allValuesLegacy = map[string]interface{}{
		"intel.apache.CPULoad":                      0.0582479,
		"intel.apache.BytesPerReq":                  1.501,
		"intel.apache.ServerVersion":                "Not Found",
		"intel.apache.ServerMPM":                    "Not Found",
		"intel.apache.Server_Built":                 "Not Found",
		"intel.apache.CurrentTime":                  "Not Found",
		"intel.apache.RestartTime":                  "Not Found",
		"intel.apache.ServerUptime":                 "Not Found",
		"intel.apache.ParentServerConfigGeneration": 0,
		"intel.apache.ParentServerMPMGeneration":    0,
		"intel.apache.Load1":                        0,
		"intel.apache.Load5":                        0,
		"intel.apache.Load15":                       0,
		"intel.apache.Total_Accesses":               10,
		"intel.apache.Total_kBytes":                 1024,
		"intel.apache.CPUUser":                      0,
		"intel.apache.CPUSystem":                    0,
		"intel.apache.CPUChildrenUser":              0,
		"intel.apache.CPUChildrenSystem":            0,
		"intel.apache.Uptime":                       4292,
		"intel.apache.ReqPerSec":                    0.123,
		"intel.apache.BytesPerSec":                  1.123,
		"intel.apache.BusyWorkers":                  10,
		"intel.apache.IdleWorkers":                  50,
		"intel.apache.ConnsTotal":                   100,
		"intel.apache.ConnsAsyncWriting":            1,
		"intel.apache.ConnsAsyncKeepAlive":          2,
		"intel.apache.ConnsAsyncClosing":            3,
		"intel.apache.workers.Closing":              1,
		"intel.apache.workers.DNSLookup":            1,
		"intel.apache.workers.Finishing":            1,
		"intel.apache.workers.Idle_Cleanup":         1,
		"intel.apache.workers.Keepalive":            1,
		"intel.apache.workers.Logging":              1,
		"intel.apache.workers.Open":                 3,
		"intel.apache.workers.Reading":              1,
		"intel.apache.workers.Sending":              1,
		"intel.apache.workers.Starting":             1,
		"intel.apache.workers.Waiting":              1,
	}
)

func TestApachePlugin(t *testing.T) {
	Convey("Create Apache Collector", t, func() {
		apacheCollector := Apache{}
		Convey("apacheCollector.GetConfigPolicy() should return a config policy", func() {
			cfgPolicy, _ := apacheCollector.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(cfgPolicy, ShouldNotBeNil)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(cfgPolicy, ShouldHaveSameTypeAs, plugin.ConfigPolicy{})
			})
		})
	})
}

func TestSafeApachePlugin(t *testing.T) {
	config := plugin.Config{
		"apache_mod_status_url": apacheURL,
		"safe":                  true,
	}
	Convey("Run collector with safe configuration", t, func() {
		apacheCollector := Apache{}

		Convey("Get Metrics Types for safe configuration and return 23 available safe metrics", func() {
			metrics, err := apacheCollector.GetMetricTypes(config)
			So(len(metrics), ShouldResemble, 23)
			So(err, ShouldBeNil)
		})

		metrics := getApacheMetrics(true)
		for i := range metrics {
			metrics[i].Config = config
		}

		Convey("Collect safe metrics from legacy endpoint", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", apacheURL,
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, legacyStatus)
					return resp, nil
				},
			)
			collectedMetrics, err := apacheCollector.CollectMetrics(metrics)
			So(err, ShouldBeNil)
			So(len(collectedMetrics), ShouldResemble, 23)
			for _, m := range collectedMetrics {
				val, metricExists := allValuesLegacy[strings.Join(m.Namespace.Strings(), ".")]
				So(metricExists, ShouldBeTrue)
				So(m.Data, ShouldEqual, val)
			}
		})

		Convey("Collect safe metrics from an updated status endpoint", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", apacheURL,
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, exampleStatus)
					return resp, nil
				},
			)
			collectedMetrics, err := apacheCollector.CollectMetrics(metrics)
			So(err, ShouldBeNil)
			So(len(collectedMetrics), ShouldResemble, 23)
			for _, m := range collectedMetrics {
				val, metricExists := allValues[strings.Join(m.Namespace.Strings(), ".")]
				So(metricExists, ShouldBeTrue)
				So(m.Data, ShouldEqual, val)
			}
		})

		Convey("Collect metrics from a bad status endpoint and ensure a failed collection", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", apacheURL,
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(401, "")
					return resp, nil
				},
			)
			collectedMetrics, err := apacheCollector.CollectMetrics(metrics)
			So(err, ShouldNotBeNil)
			So(collectedMetrics, ShouldBeNil)
		})

		Convey("Collect metrics from a bad endpoint with safe collector and ensure a failed collection", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", apacheURL,
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, "")
					return resp, nil
				},
			)
			collectedMetrics, err := apacheCollector.CollectMetrics(metrics)
			So(err, ShouldNotBeNil)
			So(collectedMetrics, ShouldBeNil)
		})
	})
}

func TestUnsafeApachePlugin(t *testing.T) {
	config := plugin.Config{
		"apache_mod_status_url": apacheURL,
		"safe":                  false,
	}
	Convey("Run collector with unsafe configuration", t, func() {
		apacheCollector := Apache{}
		Convey("Get Metrics Types for unsafe configuration and return 38 available unsafe metrics", func() {
			metrics, err := apacheCollector.GetMetricTypes(config)
			So(len(metrics), ShouldResemble, 38)
			So(err, ShouldBeNil)
		})

		metrics := getApacheMetrics(false)
		for i := range metrics {
			metrics[i].Config = config
		}

		Convey("Collect unsafe metrics from legacy endpoint", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", apacheURL,
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, legacyStatus)
					return resp, nil
				},
			)
			collectedMetrics, err := apacheCollector.CollectMetrics(metrics)
			So(err, ShouldBeNil)
			So(len(collectedMetrics), ShouldResemble, 38)
			for _, m := range collectedMetrics {
				val, metricExists := allValuesLegacy[strings.Join(m.Namespace.Strings(), ".")]
				So(metricExists, ShouldBeTrue)
				So(m.Data, ShouldEqual, val)
			}
		})

		Convey("Collect unsafe metrics from updated status endpoint", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", apacheURL,
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, exampleStatus)
					return resp, nil
				},
			)
			collectedMetrics, err := apacheCollector.CollectMetrics(metrics)
			So(err, ShouldBeNil)
			So(len(collectedMetrics), ShouldResemble, 38)
			for _, m := range collectedMetrics {
				val, metricExists := allValues[strings.Join(m.Namespace.Strings(), ".")]
				So(metricExists, ShouldBeTrue)
				So(m.Data, ShouldEqual, val)
			}
		})

		Convey("Collect metrics from a bad status endpoint with an unsafe collector", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", apacheURL,
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(401, "")
					return resp, nil
				},
			)
			collectedMetrics, err := apacheCollector.CollectMetrics(metrics)
			So(err, ShouldNotBeNil)
			So(collectedMetrics, ShouldBeNil)
		})

		// This is enforced by a required worker values in the metrics schema
		Convey("Collect metrics from a bad endpoint with unsafe collector", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder("GET", apacheURL,
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(200, "")
					return resp, nil
				},
			)
			collectedMetrics, err := apacheCollector.CollectMetrics(metrics)
			So(err, ShouldNotBeNil)
			So(collectedMetrics, ShouldBeNil)
		})
	})
}

func TestNewStatus(t *testing.T) {
	config := plugin.Config{
		"apache_mod_status_url": apacheURL,
		"safe":                  true,
	}

	Convey("Collect safe legacy endpoint", t, func() {
		apacheCollector := Apache{}
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", apacheURL,
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, exampleUnknownKeysStatus)
				return resp, nil
			},
		)
		metrics := getApacheMetrics(false)
		for i := range metrics {
			metrics[i].Config = config
		}
		_, err := apacheCollector.CollectMetrics(metrics)
		So(err, ShouldBeNil)
	})
}
