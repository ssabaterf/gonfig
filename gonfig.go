package gonfig

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadBase loads the base configuration file
func LoadBase[T interface{}](basePath string, config *T) error {
	path := fmt.Sprintf("%s/%s.yaml", basePath, "base")
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = validateConfigPath(path)
	if err != nil {
		return err
	}

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return err
	}

	return nil
}

// LoadFrom loads the configuration file for the specified enviroment
func LoadFrom[T interface{}](basePath string, enviroment string, config *T) error {
	path := fmt.Sprintf("%s/%s.yaml", basePath, enviroment)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = validateConfigPath(path)
	if err != nil {
		return err
	}

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return err
	}

	return nil
}

// LoadEnviroment loads the configuration from enviroment variables
func LoadEnviroment(startWith string, splitter string, config MergeableConfig) error {
	envVars := os.Environ()
	variables := make([]string, 0)

	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, startWith+splitter) {
			variables = append(variables, strings.Join(strings.Split(envVar, splitter)[1:], "."))
		}
	}

	for _, envVar := range variables {
		value := strings.Split(envVar, "=")[1]
		key := strings.ToLower(strings.Split(envVar, "=")[0])
		err := config.SetPathValue(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func isZeroValue(value reflect.Value) bool {
	zero := reflect.Zero(value.Type())
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == zero.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() == zero.Uint()
	case reflect.Float32, reflect.Float64:
		return value.Float() == zero.Float()
	case reflect.Complex64, reflect.Complex128:
		return value.Complex() == zero.Complex()
	case reflect.Bool:
		return value.Bool() == zero.Bool()
	default:
		return reflect.DeepEqual(value.Interface(), zero.Interface())
	}
}
func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// MergeableConfig is the interface that wraps the MergeWith and SetPathValue methods.
type MergeableConfig interface {
	MergeWith(nextConfig interface{}) error
	SetPathValue(path string, value interface{}) error
}
