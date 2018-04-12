package gomanager

type SimpleConfig struct {
	file  string
	obj   interface{}
	bytes []byte
}

// NewSimpleConfig...
func NewSimpleConfig(file string, obj interface{}) (IConfig, error) {
	if bytes, err := ReadFile(file, obj); err != nil {
		return nil, err
	} else {
		return &SimpleConfig{
			file:  file,
			obj:   obj,
			bytes: bytes,
		}, err
	}
}

// Get ...
func (setup *SimpleConfig) Get() interface{} {
	return setup.obj
}

// Set ...
func (setup *SimpleConfig) Set(config interface{}) {
	setup.obj = config
}

// Reload ...
func (setup *SimpleConfig) Reload() error {
	if bytes, err := ReadFile(setup.file, setup.obj); err != nil {
		return err
	} else {
		setup.bytes = bytes
	}

	return nil
}

// Save ...
func (setup *SimpleConfig) Save() error {
	if err := WriteFile(setup.file, setup.obj); err != nil {
		return err
	}

	return nil
}
