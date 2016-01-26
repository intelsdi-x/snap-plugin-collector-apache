package apache

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	// Name of plugin
	Name = "apache"
	// Version of plugin
	Version = 1
	// Type of plugin
	Type = plugin.CollectorPluginType
)

var (
	errNoWebserver  = errors.New("apache_mod_status_url config required. Check your config JSON file")
	errBadWebserver = errors.New("Failed to parse given apache_mod_status_url")
	errReqFailed    = errors.New("Request to Apache webserver failed")

	scoreboard = map[string]string{
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

type Apache struct{}

func getMetrics(webserver string, metrics []string) ([]plugin.PluginMetricType, error) {
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(webserver)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		return nil, fmt.Errorf("%v,%v", errReqFailed, err)
	}
	defer resp.Body.Close()

	mtsmap := make(map[string]plugin.PluginMetricType)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var ns string
		var scoreboardflag bool
		line := scanner.Text()
		lineslice := strings.Split(line, ": ")
		if strings.Contains(line, ":") {
			if strings.Contains(line, "Scoreboard") {
				scoreboardflag = true
			} else {
				ns = strings.Replace(lineslice[0], " ", "_", -1)
				data, err := strconv.ParseFloat(lineslice[1], 64)
				if err != nil {
					return nil, err
				}
				mtsmap[ns] = plugin.PluginMetricType{
					Namespace_: []string{"intel", "apache", ns},
					Data_:      data,
					Source_:    webserver,
					Timestamp_: time.Now(),
				}
			}
		}
		if scoreboardflag {
			for ns, _ := range scoreboard {
				data := strings.Count(line, scoreboard[ns])
				mtsmap[ns] = plugin.PluginMetricType{
					Namespace_: []string{"intel", "apache", "scoreboard", ns},
					Data_:      data,
					Source_:    webserver,
					Timestamp_: time.Now(),
				}
			}
		}
	}
	if len(metrics) == 0 {
		mts := make([]plugin.PluginMetricType, 0, len(mtsmap))
		for _, v := range mtsmap {
			mts = append(mts, v)
		}
		return mts, nil
	}
	mts := make([]plugin.PluginMetricType, 0, len(metrics))
	for _, v := range metrics {
		mt, ok := mtsmap[v]
		if ok {
			mts = append(mts, mt)
		}
	}
	return mts, nil
}

func (a *Apache) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	config := mts[0].Config().Table()
	webservercfg, ok := config["apache_mod_status_url"]
	if !ok {
		return nil, errNoWebserver
	}
	webserver, ok := webservercfg.(ctypes.ConfigValueStr)
	if !ok {
		return nil, errBadWebserver
	}
	metrics := make([]string, len(mts))
	for i, m := range mts {
		metrics[i] = m.Namespace_[len(m.Namespace_)-1]
	}
	return getMetrics(webserver.Value, metrics)
}

func (a *Apache) GetMetricTypes(cfg plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	webservercfg, ok := cfg.Table()["apache_mod_status_url"]
	if !ok {
		return getMetrics("http://127.0.0.1:80/server-status?auto", []string{})
	}
	webserver, ok := webservercfg.(ctypes.ConfigValueStr)
	if !ok {
		return nil, errBadWebserver
	}
	return getMetrics(webserver.Value, []string{})
}

func (a *Apache) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cfg := cpolicy.New()
	rule, _ := cpolicy.NewStringRule("apache_mod_status_url", false, "http://127.0.0.1:80/server-status?auto")
	policy := cpolicy.NewPolicyNode()
	policy.Add(rule)
	cfg.Add([]string{"intel", "apache"}, policy)
	return cfg, nil
}

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		Name,
		Version,
		Type,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.Unsecure(true),
		plugin.RoutingStrategy(plugin.DefaultRouting),
	)
}
