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
	"fmt"
	"time"

	"github.com/gorilla/schema"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

// Status allows for the deserializtion of metrics and the matching of requested metrics
type Status interface {
	ReturnDesignatedMetrics(mts []plugin.Metric) ([]plugin.Metric, error)
}

// unsafeStatus allows for the deserializtion of metrics from an apache server-status endpoint
type apacheStatus struct {
	ServerVersion                string  `schema:"ServerVersion"`
	ServerMPM                    string  `schema:"ServerMPM"`
	ServerBuilt                  string  `schema:"Server_Built"`
	CurrentTime                  string  `schema:"CurrentTime"`
	RestartTime                  string  `schema:"RestartTime"`
	ServerUptime                 string  `schema:"ServerUptime"`
	Load1                        float64 `schema:"Load1"`
	Load5                        float64 `schema:"Load5"`
	Load15                       float64 `schema:"Load15"`
	CPULoad                      float64 `schema:"CPULoad"`
	ReqPerSec                    float64 `schema:"ReqPerSec"`
	BytesPerSec                  float64 `schema:"BytesPerSec"`
	BytesPerReq                  float64 `schema:"BytesPerReq"`
	CPUUser                      float64 `schema:"CPUUser"`
	CPUSystem                    float64 `schema:"CPUSystem"`
	CPUChildrenUser              float64 `schema:"CPUChildrenUser"`
	CPUChildrenSystem            float64 `schema:"CPUChildrenSystem"`
	ParentServerConfigGeneration float64 `schema:"ParentServerConfigGeneration"`
	ParentServerMPMGeneration    float64 `schema:"ParentServerMPMGeneration"`
	TotalAccesses                int     `schema:"Total_Accesses"`
	TotalkBytes                  int     `schema:"Total_kBytes"`
	ServerUptimeSeconds          int     `schema:"ServerUptimeSeconds"`
	Uptime                       int     `schema:"Uptime"`
	BusyWorkers                  int     `schema:"BusyWorkers"`
	IdleWorkers                  int     `schema:"IdleWorkers"`
	ConnsTotal                   int     `schema:"ConnsTotal"`
	ConnsAsyncWriting            int     `schema:"ConnsAsyncWriting"`
	ConnsAsyncKeepAlive          int     `schema:"ConnsAsyncKeepAlive"`
	ConnsAsyncClosing            int     `schema:"ConnsAsyncClosing"`
	Closing                      int     `schema:"Closing,required"`
	DNSLookup                    int     `schema:"DNSLookup,required"`
	Finishing                    int     `schema:"Finishing,required"`
	IdleCleanup                  int     `schema:"Idle_Cleanup,required"`
	Keepalive                    int     `schema:"Keepalive,required"`
	Logging                      int     `schema:"Logging,required"`
	Open                         int     `schema:"Open,required"`
	Reading                      int     `schema:"Reading,required"`
	Sending                      int     `schema:"Sending,required"`
	Starting                     int     `schema:"Starting,required"`
	Waiting                      int     `schema:"Waiting,required"`
}

var (
	workerMetrics = map[string]string{
		"Closing":      "workers",
		"DNSLookup":    "workers",
		"Finishing":    "workers",
		"Idle_Cleanup": "workers",
		"Keepalive":    "workers",
		"Logging":      "workers",
		"Open":         "workers",
		"Reading":      "workers",
		"Sending":      "workers",
		"Starting":     "workers",
		"Waiting":      "workers",
	}

	safeMetrics = map[string]string{
		"BusyWorkers":         "workers",
		"BytesPerSec":         "B",
		"ReqPerSec":           "req",
		"ConnsAsyncClosing":   "conn",
		"ConnsAsyncKeepAlive": "conn",
		"ConnsAsyncWriting":   "conn",
		"ConnsTotal":          "conn",
		"CPULoad":             "load",
		"IdleWorkers":         "workers",
		"Total_Accesses":      "req",
		"Total_kBytes":        "B",
		"Uptime":              "s",
	}

	unsafeMetrics = map[string]string{
		"ServerVersion":                "",
		"ServerMPM":                    "",
		"Server_Built":                 "",
		"CurrentTime":                  "",
		"RestartTime":                  "",
		"ParentServerConfigGeneration": "",
		"ParentServerMPMGeneration":    "",
		"BytesPerReq":                  "B",
		"CPUUser":                      "jiff",
		"CPUSystem":                    "jiff",
		"CPUChildrenUser":              "jiff",
		"CPUChildrenSystem":            "jiff",
		"Load1":                        "load/1M",
		"Load5":                        "load/5M",
		"Load15":                       "load/15M",
	}
)

// NewStatus returns a new http status object
func NewStatus(metricMap map[string][]string) (Status, error) {
	status := &apacheStatus{
		ServerVersion: "Not Found",
		ServerMPM:     "Not Found",
		ServerBuilt:   "Not Found",
		CurrentTime:   "Not Found",
		RestartTime:   "Not Found",
		ServerUptime:  "Not Found",
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(status, metricMap)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func getApacheMetrics(safe bool) []plugin.Metric {
	mts := []plugin.Metric{}
	for metric, unit := range safeMetrics {
		mts = append(mts, plugin.Metric{
			Namespace: plugin.NewNamespace("intel", "apache", metric),
			Unit:      unit,
		})
	}
	for metric, unit := range workerMetrics {
		mts = append(mts, plugin.Metric{
			Namespace: plugin.NewNamespace("intel", "apache", "workers", metric),
			Unit:      unit,
		})
	}
	if !safe {
		for metric, unit := range unsafeMetrics {
			mts = append(mts, plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "apache", metric),
				Unit:      unit,
			})
		}
	}
	return mts
}

// ReturnDesignatedMetrics returns a slice of metrics that match the metrics requested from the plugin
func (a apacheStatus) ReturnDesignatedMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	returnMetrics := []plugin.Metric{}
	now := time.Now()
	for _, metric := range mts {
		met := plugin.Metric{
			Namespace:   metric.Namespace,
			Tags:        metric.Tags,
			Unit:        metric.Unit,
			Timestamp:   now,
			Config:      metric.Config,
			Description: metric.Description,
		}
		switch metric.Namespace.Strings()[2] {
		case "ServerVersion":
			met.Data = a.ServerVersion
		case "ServerMPM":
			met.Data = a.ServerMPM
		case "Server_Built":
			met.Data = a.ServerBuilt
		case "CurrentTime":
			met.Data = a.CurrentTime
		case "RestartTime":
			met.Data = a.RestartTime
		case "ParentServerConfigGeneration":
			met.Data = a.ParentServerConfigGeneration
		case "ParentServerMPMGeneration":
			met.Data = a.ParentServerMPMGeneration
		case "ServerUptime":
			met.Data = a.ServerUptime
		case "BusyWorkers":
			met.Data = a.BusyWorkers
		case "BytesPerReq":
			met.Data = a.BytesPerReq
		case "BytesPerSec":
			met.Data = a.BytesPerSec
		case "ReqPerSec":
			met.Data = a.ReqPerSec
		case "ConnsAsyncClosing":
			met.Data = a.ConnsAsyncClosing
		case "CPUUser":
			met.Data = a.CPUUser
		case "CPUSystem":
			met.Data = a.CPUSystem
		case "CPUChildrenUser":
			met.Data = a.CPUChildrenUser
		case "CPUChildrenSystem":
			met.Data = a.CPUChildrenSystem
		case "ConnsAsyncKeepAlive":
			met.Data = a.ConnsAsyncKeepAlive
		case "ConnsAsyncWriting":
			met.Data = a.ConnsAsyncWriting
		case "ConnsTotal":
			met.Data = a.ConnsTotal
		case "CPULoad":
			met.Data = a.CPULoad
		case "Load1":
			met.Data = a.Load1
		case "Load5":
			met.Data = a.Load5
		case "Load15":
			met.Data = a.Load15
		case "IdleWorkers":
			met.Data = a.IdleWorkers
		case "Total_Accesses":
			met.Data = a.TotalAccesses
		case "Total_kBytes":
			met.Data = a.TotalkBytes
		case "Uptime":
			if a.Uptime == 0 {
				met.Data = a.ServerUptimeSeconds
			} else {
				met.Data = a.Uptime
			}
		case "workers":
			switch metric.Namespace.Strings()[3] {
			case "Closing":
				met.Data = a.Closing
			case "DNSLookup":
				met.Data = a.DNSLookup
			case "Finishing":
				met.Data = a.Finishing
			case "Idle_Cleanup":
				met.Data = a.IdleCleanup
			case "Keepalive":
				met.Data = a.Keepalive
			case "Logging":
				met.Data = a.Logging
			case "Open":
				met.Data = a.Open
			case "Reading":
				met.Data = a.Reading
			case "Sending":
				met.Data = a.Sending
			case "Starting":
				met.Data = a.Starting
			case "Waiting":
				met.Data = a.Waiting
			default:
				return nil, fmt.Errorf("metric does not exist: %v", metric.Namespace.Strings())
			}
		default:
			return nil, fmt.Errorf("metric does not exist: %v", metric.Namespace.Strings())
		}
		returnMetrics = append(returnMetrics, met)
	}
	return returnMetrics, nil
}
