package config

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"github.com/toolkits/pkg/file"
)

type ConfYaml struct {
	Logger LoggerSection       `yaml:"logger"`
	Probe  probeSection        `yaml:"probe"`
	Ping   map[string][]string `yaml:"ping"`
	Url    map[string][]string `yaml:"url"`
	Server serverSection       `yaml:"server"`
}

type serverSection struct {
	RpcMethod string `yaml:"rpcMethod"`
}

type probeSection struct {
	Region   string            `yaml:"region"`
	Timeout  int64             `yaml:"timeout"`
	Limit    int64             `yaml:"limit"`
	Interval int               `yaml:"interval"`
	Headers  map[string]string `yaml:"headers"`
}

var (
	Config   *ConfYaml
	lock     = new(sync.RWMutex)
	Endpoint string
	Cwd      string
)

// Get configuration file
func Get() *ConfYaml {
	lock.RLock()
	defer lock.RUnlock()
	return Config
}

func Parse(conf string) error {
	bs, err := file.ReadBytes(conf)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	lock.Lock()
	defer lock.Unlock()

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(bs))
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	viper.SetDefault("server.rpcMethod", "Transfer.Push")

	err = viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("Unmarshal %v", err)
	}

	return nil
}
