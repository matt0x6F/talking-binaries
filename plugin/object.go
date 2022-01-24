package plugin

// Object is the representation of a plugin
type Object struct {
	// reference used throughout steps
	Ref string
	// The plugin binary name. Will be searched for on $PATH or relative to the current directory.
	Name string
	// The absolute path of the plugin binary. Path can be, and is usually, unset and determined.
	Path string
}
