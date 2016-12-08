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
	"bufio"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

var (
	errBadWebserver = errors.New("Failed to parse given apache_mod_status_url")
	errReqFailed    = errors.New("Request to Apache webserver failed")
	errBadConfig    = errors.New("Failed to parse given safe config")
	workers         = map[string]string{
		"Closing":      "C",
		"DNSLookup":    "D",
		"Finishing":    "G",
		"Idle_Cleanup": "I",
		"Keepalive":    "K",
		"Logging":      "L",
		"Open":         ".",
		"Reading":      "R",
		"Sending":      "W",
		"Starting":     "S",
		"Waiting":      "_",
	}
)

//Apache plugin struct which gathers plugin specific data
type Apache struct{}

func parseMetrics(resp *http.Response) (map[string][]string, error) {
	mtsmap := make(map[string][]string)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var ns string
		line := scanner.Text()
		lineslice := strings.Split(line, ": ")
		if strings.Contains(line, "Scoreboard") {
			line = strings.Trim(line, "Scoreboard")
			for ns := range workers {
				data := []string{strconv.Itoa(strings.Count(line, workers[ns]))}
				mtsmap[ns] = data
			}
		} else {
			if len(lineslice) > 1 {
				ns = strings.Replace(lineslice[0], " ", "_", -1)
				mtsmap[ns] = []string{lineslice[1]}
			}
		}
	}
	return mtsmap, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (a Apache) CollectMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	webserver, err := mts[0].Config.GetString("apache_mod_status_url")
	if err != nil {
		return nil, errBadWebserver
	}
	client := &http.Client{}
	resp, err := client.Get(webserver)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errReqFailed
	}
	metricMap, err := parseMetrics(resp)
	if err != nil {
		return nil, err
	}
	status, err := NewStatus(metricMap)
	if err != nil {
		return nil, err
	}
	return status.ReturnDesignatedMetrics(mts)
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (a Apache) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	safe, ok := cfg.GetBool("safe")
	if ok != nil {
		return nil, errBadConfig
	}
	mts := getApacheMetrics(safe)
	return mts, nil
}

// GetConfigPolicy returns config policy
// It returns error in case retrieval was not successful
func (a Apache) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	cfg := plugin.NewConfigPolicy()
	cfg.AddNewStringRule(
		[]string{"intel", "apache"},
		"apache_mod_status_url",
		false,
		plugin.SetDefaultString("http://127.0.0.1:80/server-status?auto"),
	)
	cfg.AddNewBoolRule(
		[]string{"intel", "apache"},
		"safe",
		false,
		plugin.SetDefaultBool(true),
	)
	return *cfg, nil
}
