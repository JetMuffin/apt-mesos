package impl

import (
	"encoding/json"
	"net/http"

	"github.com/icsnju/apt-mesos/registry"
)

// FetchMetricData connect to mesos-master and get raw metric data,
func (core *Core) FetchMetricData() (*registry.MetricsData, error) {
	resp, err := http.Get("http://" + core.master + "/master/state.json")
	if err != nil {
		return nil, err
	}

	data := new(registry.MetricsData)
	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}
	resp.Body.Close()

	return data, nil
}
