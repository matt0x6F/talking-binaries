package protocol

// Input protocol is passed as configuration to the called application
type Input struct {
	// Config is the configuration passed to the binary
	Config map[string]interface{}
	// Indicates that a plugin is configured to be called twice in parallel execution
	Parallel bool
}

// Output protocol is used by plugins to return information
type Output struct {
	// Output to show the user
	UserOutput map[int][]byte
	Errors     map[int][]byte
	PluginLogs []byte
}

func NewOutput() Output {
	output := Output{}

	output.UserOutput = make(map[int][]byte)
	output.Errors = make(map[int][]byte)

	return output
}
