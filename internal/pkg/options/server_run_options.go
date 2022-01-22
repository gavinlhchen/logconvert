package options

import (
	"github.com/spf13/pflag"
)

type ServerRunOptions struct {
	IsaGlobalConfigPath string `json:"isa-global-config-path"   mapstructure:"isa-global-config-path"`
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters.
func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		IsaGlobalConfigPath: "/usr/local/app/pcmgr_bigdata/isa_global.ini",
	}
}

// Validate checks validation of ServerRunOptions.
func (s *ServerRunOptions) Validate() []error {
	errors := []error{}

	return errors
}

func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.IsaGlobalConfigPath, "server.isa-global-config-path", s.IsaGlobalConfigPath, ""+
		"Isa global config path.")
}
