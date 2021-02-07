// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type CheckConfig struct {
	Name     string   `yaml:"name"`
	Path     string   `yaml:"path"`
	Category string   `yaml:"category"`
	Period   string   `yaml:"period"`
	Enabled  bool     `yaml:"enabled"`
	Program  string   `yaml:"program"`
	Params   []string `yaml:params`
}

type Config struct {
	Checks []CheckConfig
}

var DefaultConfig = Config{}
