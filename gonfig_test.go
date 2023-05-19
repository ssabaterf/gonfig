package gonfig

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGonfig(t *testing.T) {
	os.Setenv("APP_SERVER_HOST", "localhost")
	os.Setenv("APP_SERVER_PORT", "3000")
	os.Setenv("APP_SERVER_TIMEOUT_IDLE", "3")
	config := MyConfig{}
	err := LoadBase("./configuration", &config)
	if err != nil {
		t.Errorf("Error loading config: %s", err.Error())
	}

	confProd := MyConfig{}
	err = LoadFrom("./configuration", "prod", &confProd)

	err = config.MergeWith(&confProd)
	if err != nil {
		t.Errorf("Error merging config: %s", err.Error())
	}

	err = LoadEnviroment("APP", "_", &config)
	if err != nil {
		t.Errorf("Error loading config: %s", err.Error())
	}

	assert.Equal(t, config.Server.Host, "localhost")
	assert.Equal(t, config.Server.Port, 3000)
	assert.Equal(t, config.Server.Timeout.Idle, uint32(3))
	assert.Equal(t, config.Server.Timeout.Server, uint32(30))
	assert.Equal(t, config.Server.Timeout.Write, uint32(50))
	assert.Equal(t, config.Server.Timeout.Read, uint32(15))
}

type MyConfig struct {
	Server struct {
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Timeout struct {
			Server uint32 `yaml:"server"`
			Write  uint32 `yaml:"write"`
			Read   uint32 `yaml:"read"`
			Idle   uint32 `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
}

func (c *MyConfig) MergeWith(nextConfig interface{}) error {
	next, ok := nextConfig.(*MyConfig)
	if !ok {
		return fmt.Errorf("'%s' is not a valid config", nextConfig)
	}

	if next.Server.Host != "" {
		c.Server.Host = next.Server.Host
	}
	if next.Server.Port != 0 {
		c.Server.Port = next.Server.Port
	}
	if next.Server.Timeout.Server != 0 {
		c.Server.Timeout.Server = next.Server.Timeout.Server
	}
	if next.Server.Timeout.Write != 0 {
		c.Server.Timeout.Write = next.Server.Timeout.Write
	}
	if next.Server.Timeout.Read != 0 {
		c.Server.Timeout.Read = next.Server.Timeout.Read
	}
	if next.Server.Timeout.Idle != 0 {
		c.Server.Timeout.Idle = next.Server.Timeout.Idle
	}
	return nil
}
func (c *MyConfig) SetPathValue(path string, value interface{}) error {
	switch path {
	case "server.host":
		c.Server.Host = value.(string)
	case "server.port":
		{
			val, err := strconv.Atoi(value.(string))
			if err != nil {
				return err
			}
			c.Server.Port = val
		}
	case "server.timeout.server":
		{
			val, err := strconv.Atoi(value.(string))
			if err != nil {
				return err
			}
			c.Server.Timeout.Server = uint32(val)
		}
	case "server.timeout.write":
		{
			val, err := strconv.Atoi(value.(string))
			if err != nil {
				return err
			}
			c.Server.Timeout.Write = uint32(val)
		}
	case "server.timeout.read":
		{
			val, err := strconv.Atoi(value.(string))
			if err != nil {
				return err
			}
			c.Server.Timeout.Read = uint32(val)
		}
	case "server.timeout.idle":
		{
			val, err := strconv.Atoi(value.(string))
			if err != nil {
				return err
			}
			c.Server.Timeout.Idle = uint32(val)
		}
	default:
		return fmt.Errorf("'%s' is not a valid path", path)
	}
	return nil
}
