package universal

import (
	"encoding/json"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
)

var defaultMetadata = Metadata{
	Name:        "app",
	Version:     "1.0.0",
	Description: "",
}

type ChartConfig struct {
	Config Config `json:"_config"`
}

type Config struct {
	Metadata    `json:"_metadata"`
	Controllers []*Controller `json:"controllers"`
}

type Metadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func NewEmptyChartConfig() *ChartConfig {
	return &ChartConfig{
		Config: Config{
			Metadata: defaultMetadata,
		},
	}
}

func PrepareChartConfig(chart string, cur int) (*ChartConfig, error) {
	chconfig := NewEmptyChartConfig()
	if chart == "" {
		chart = "{}"
	}
	err := json.Unmarshal([]byte(chart), chconfig)
	if err != nil {
		glog.Errorf("json Unmarshal error: %v", err)
		return nil, err
	}

	config := chconfig.Config
	switch {
	case len(config.Controllers) == 0:
		config.Controllers = make([]*Controller, cur+1)
		glog.V(4).Infof("controllers are empty, which will be appended %d empty elements", cur+1)
		glog.V(4).Infof("the object will be encoded in controllers[%d]", cur)
	case len(config.Controllers) < cur+1:
		glog.V(4).Infof("controllers length: %d", len(config.Controllers))
		tmp := make([]*Controller, cur+1-len(config.Controllers))
		config.Controllers = append(config.Controllers, tmp...)
		glog.V(4).Infof("controllers will be appended by %d empty elements, the new controllers length: %d", len(tmp), cur+1)
		glog.V(4).Infof("controllers[%d] is zero value struct, the object config will be encoded in it", cur)
	default:
		glog.V(4).Infof("controllers[%d] will be modified, the object config will be encoded in it", cur)
		glog.V(4).Infof("controllers[%d]'s original config: %s", cur, spew.Sdump(config.Controllers[cur]))
	}
	chconfig.Config = config
	return chconfig, nil
}
