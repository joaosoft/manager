package manager

var global = make(map[string]interface{})

func init() {
	global[path_key] = defaultPath
}
