# Gonfig

a simple config file parser for golang

**Installation**

```bash
go get github.com/ssabater/gonfig
```

**Usage**

Simple load from yaml

```go
package main

import (
    "fmt"
    "github.com/ssabater/gonfig"
)

func main() {
    config := MyConfig{}
	err := gonfig.LoadBase("./configuration", &config)
	if err != nil {
		panic(err)
	}
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
```

**Load from yaml with merge and environment variables**

For MergeWith and LoadEnviroment must implement the interface "MergeableConfig"
for the struct to be filled, in this case "MyConfig"

```go
type MergeableConfig interface {
	MergeWith(nextConfig interface{}) error
	SetPathValue(path string, value interface{}) error
}
```

```go
package main

import (
    "fmt"
    "github.com/ssabater/gonfig"
)

func main() {
    //simulating some environment variables
    os.Setenv("APP_SERVER_HOST", "localhost")
	os.Setenv("APP_SERVER_PORT", "3000")
	os.Setenv("APP_SERVER_TIMEOUT_IDLE", "3")

    config := MyConfig{}
    err := gonfig.LoadBase("./configuration", &config)
    if err != nil {
        panic(err)
    }

    //Overriding base with production values, "prod" is the name of the file and can be replaced by any string always that the file exists
    confProd := MyConfig{}
	err = gonfig.LoadFrom("./configuration", "prod", &confProd)
    if err != nil {
        panic(err)
    }

    //merging base with production
    err = config.MergeWith(&confProd)
	if err != nil {
		panic(err)
	}

    //merging with environment variables
    //"APP" is the prefix of the environment variables, 
    //"_" is the separator and "config" is the struct to be filled
    err = gonfig.LoadEnviroment("APP", "_", &config)
	if err != nil {
		panic(err)
	}
    
    fmt.Printf("%+v\n", config)
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
```
