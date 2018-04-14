package gomanager

import (
	"time"

	"github.com/spf13/viper"
)

type SimpleConfig struct {
	file  string
	obj   interface{}
	bytes []byte
	viper *viper.Viper
}

// NewSimpleConfig...
func NewSimpleConfig(file string, obj interface{}) (IConfig, error) {
	if bytes, err := readFile(file, obj); err != nil {
		return nil, err
	} else {
		return &SimpleConfig{
			file:  file,
			obj:   obj,
			bytes: bytes,
			viper: loadViper(file),
		}, err
	}
}

// Get ...
func (simple *SimpleConfig) Get(key string) interface{} {
	return simple.viper.Get(key)
}

// GetString ...
func (simple *SimpleConfig) GetString(key string) string {
	return simple.viper.GetString(key)
}

// GetBool ...
func (simple *SimpleConfig) GetBool(key string) bool {
	return simple.viper.GetBool(key)
}

// GetInt ...
func (simple *SimpleConfig) GetInt(key string) int {
	return simple.viper.GetInt(key)
}

// GetInt64 ...
func (simple *SimpleConfig) GetInt64(key string) int64 {
	return simple.viper.GetInt64(key)
}

// GetFloat64 ...
func (simple *SimpleConfig) GetFloat64(key string) float64 {
	return simple.viper.GetFloat64(key)
}

// GetTime ...
func (simple *SimpleConfig) GetTime(key string) time.Time {
	return simple.viper.GetTime(key)
}

// GetDuration ...
func (simple *SimpleConfig) GetDuration(key string) time.Duration {
	return simple.viper.GetDuration(key)
}

// GetStringSlice ...
func (simple *SimpleConfig) GetStringSlice(key string) []string {
	return simple.viper.GetStringSlice(key)
}

// GetStringMap ...
func (simple *SimpleConfig) GetStringMap(key string) map[string]interface{} {
	return simple.viper.GetStringMap(key)
}

// GetStringMapString ...
func (simple *SimpleConfig) GetStringMapString(key string) map[string]string {
	return simple.viper.GetStringMapString(key)
}

// GetStringMapStringSlice ...
func (simple *SimpleConfig) GetStringMapStringSlice(key string) map[string][]string {
	return simple.viper.GetStringMapStringSlice(key)
}

// GetObj ...
func (simple *SimpleConfig) GetObj() interface{} {
	return simple.obj
}

// Set ...
func (simple *SimpleConfig) Set(config interface{}) {
	simple.obj = simple
}

// Reload ...
func (simple *SimpleConfig) Reload() error {
	if bytes, err := readFile(simple.file, simple.obj); err != nil {
		return err
	} else {
		simple.viper = loadViper(simple.file)
		simple.bytes = bytes
	}

	return nil
}

// Save ...
func (simple *SimpleConfig) Save() error {
	if err := writeFile(simple.file, simple.obj); err != nil {
		return err
	}

	simple.viper = loadViper(simple.file)

	return nil
}

func loadViper(file string) *viper.Viper {
	viper := viper.New()
	viper.AddConfigPath(file)
	viper.ReadInConfig()

	return viper
}
