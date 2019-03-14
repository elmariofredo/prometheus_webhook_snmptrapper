package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

// type Config struct {
// 	SNMPTrapAddress string
// 	SNMPCommunity   string
// 	SNMPRetries     uint
// 	WebhookAddress  string
// 	CongifFile      string
// }

// Secret is a string that must not be revealed on marshaling.
type Secret string

// MarshalYAML implements the yaml.Marshaler interface.
func (s Secret) MarshalYAML() (interface{}, error) {
	if s != "" {
		return "<secret>", nil
	}
	return nil, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for Secrets.
func (s *Secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Secret
	return unmarshal((*plain)(s))
}

// LoadConfig parses the YAML input into a Config.
func LoadConfig(s string) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	log.V(1).Infof("Loaded config:\n%+v", cfg)
	return cfg, nil
}

// LoadConfigFile parses the given YAML file into a Config.
func LoadConfigFile(filename string) (*Config, []byte, error) {
	log.V(1).Infof("Loading configuration from %q", filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	cfg, err := LoadConfig(string(content))
	if err != nil {
		return nil, nil, err
	}

	//resolveFilepaths(filepath.Dir(filename), cfg)
	return cfg, content, nil
}

// resolveFilepaths joins all relative paths in a configuration
// with a given base directory.
// func resolveFilepaths(baseDir string, cfg *Config) {
// 	join := func(fp string) string {
// 		if len(fp) == 0 || filepath.IsAbs(fp) {
// 			return fp
// 		}
// 		absFp := filepath.Join(baseDir, fp)
// 		log.V(2).Infof("Relative path %q resolved to %q", fp, absFp)
// 		return absFp
// 	}

// 	cfg.Template = join(cfg.Template)
// }

// OidConfig is the configuration for one receiver. It has a unique name and includes API access fields (URL, user
// and password) and issue fields (required -- e.g. project, issue type -- and optional -- e.g. priority).
type OidConfig struct {
	OidName   string `yaml:"Name" json:"Name"`
	OidNumber string `yaml:"Oid" json:"Oid"`
	Template  string `yaml:"Template" json:"Template"`
	Type      string `yaml:"Type" json:"Type"`
	NotEmpty  bool   `yaml:"NotEmpty" json:"NotEmpty"`
	// Fields     map[string]interface{} `yaml:"fields" json:"fields"`
	// Components []string               `yaml:"components" json:"components"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline" json:"-"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *OidConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain OidConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	// Recursively convert any maps to map[string]interface{}, filtering out all non-string keys, so the json encoder
	// doesn't blow up when marshaling JIRA requests.
	// fieldsWithStringKeys, err := tcontainer.ConvertToMarshalMap(c.Fields, func(v string) string { return v })
	// if err != nil {
	// 	return err
	// }
	// rc.Fields = fieldsWithStringKeys
	return checkOverflow(c.XXX, "oid")
}

// Config is the top-level configuration for JIRAlert's config file.
type Config struct {
	//	Defaults       *OidConfig   `yaml:"defaults,omitempty" json:"defaults,omitempty"`
	Oids            []*OidConfig `yaml:"Oids,omitempty" json:"Oids,omitempty"`
	FiringTrap      string       `yaml:"FiringTrap" json:"FiringTrap"`
	RecoveryTrap    string       `yaml:"RecoveryTrap" json:"RecoveryTrap"`
	SNMPTrapAddress string       `yaml:"TrapAddress" json:"TrapAddress"`
	SNMPCommunity   string       `yaml:"Community" json:"Community"`
	SNMPRetries     uint         `yaml:"Retries" json:"Retries"`
	WebhookAddress  string       `yaml:"WebhookAddress" json:"WebhookAddress"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline" json:"-"`
}

func (c Config) String() string {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.FiringTrap == "" {
		return fmt.Errorf("missing FiringTrap")
	}
	if c.RecoveryTrap == "" {
		return fmt.Errorf("missing RecoveryTrap")
	}
	if c.SNMPTrapAddress == "" {
		return fmt.Errorf("missing TrapAddress")
	}
	if c.SNMPCommunity == "" {
		return fmt.Errorf("missing Community")
	}
	if fmt.Sprint(c.SNMPRetries) == "" {
		return fmt.Errorf("missing Retries")
	}
	if c.WebhookAddress == "" {
		return fmt.Errorf("missing WebhookAddress")
	}

	for _, oid := range c.Oids {

		if oid.OidName == "" {
			return fmt.Errorf("missing OidName for oid %+v", oid)
		}

		if oid.OidNumber == "" {
			return fmt.Errorf("missing OidNumber for oid %+v", oid.OidName)
		}

		if oid.Template == "" {
			return fmt.Errorf("missing Template for oid %+v", oid.OidName)
		}

		if oid.Type != "" {
			if oid.Type == "int32" {

			} else if oid.Type == "string" {

			} else {
				return fmt.Errorf("Wrong Type '%s' for oid %+v", oid.Type, oid.OidName)
			}

		}

		//fmt.Printf("NotEmpty: %v\n", oid.NotEmpty)
	}

	if len(c.Oids) == 0 {
		return fmt.Errorf("no Oids defined")
	}

	return checkOverflow(c.XXX, "config")
}

//OidName loops the receiver list and returns the first instance with that oid
func (c *Config) OidName(OidName string) *OidConfig {
	for _, oid := range c.Oids {
		if oid.OidName == OidName {
			return oid
		}
	}
	return nil
}

func checkOverflow(m map[string]interface{}, ctx string) error {
	if len(m) > 0 {
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		log.Warningf("unknown fields in %s: %s", ctx, strings.Join(keys, ", "))
	}
	return nil
}
