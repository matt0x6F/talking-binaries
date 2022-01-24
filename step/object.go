package step

const (
	Waiting = iota
	Completed
)

// Object is the representation of a step
type Object struct {
	// Human-friendly name of the step; also used as a reference.
	Name string
	// The reference name of the plugin
	Plugin string
	// Configuration to be passed to the plugin.
	Config map[string]interface{}
	// Parallel indicates that the step can be performed in parallel with other steps configured in parallel. All
	// sequential parallel steps will be executed in parallel together, forming a "run group".
	Parallel bool
	// Next is the next step
	Next *Object
	// Prev is the previous step
	Prev *Object
	// Status is the status of the step
	Status int
}
